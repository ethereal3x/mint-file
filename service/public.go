package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"path"
	"time"
)

type UploadBaseAggregation struct {
	Ctx        context.Context
	Location   string
	FileName   string
	BucketName string
}

type DownloadBaseAggregation struct {
	Ctx        context.Context
	Location   string
	FileName   string
	BucketName string
}

type GenerateBaseAggregation struct {
	Ctx           context.Context
	Location      string
	FileName      string
	BucketName    string
	EffectiveDate time.Duration
}

type DownloadFileAggregation struct {
	Aggregation *DownloadBaseAggregation
	Data        []byte
}

type UploadBinaryFileAggregation struct {
	Aggregation *UploadBaseAggregation
	Data        []byte
}

type UploadUrlFileFileAggregation struct {
	Aggregation *UploadBaseAggregation
	Url         string
}

type UploadFragmentFileAggregation struct {
	Aggregation *UploadBaseAggregation
	Data        []byte
	ShardSize   int64
}

type ListUploadedFragmentsAggregation struct {
	Aggregation *UploadBaseAggregation
	UploadID    string
}

type CancelFragmentUploadAggregation struct {
	Aggregation *UploadBaseAggregation
	UploadID    string
}

func GenerateUniqueFilePath(dir, fileName string) string {
	return path.Join(dir, fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName))
}

func GenerateDownloadFilePath(dir, fileName string) string {
	return path.Join(dir, fileName)
}

func CheckVolceTosErr(err error) {
	if err != nil {
		var serverErr *tos.TosServerError
		if errors.As(err, &serverErr) {
			hlog.Infof("Error:%v", serverErr.Error())
			hlog.Infof("Request ID:%v", serverErr.RequestID)
			hlog.Infof("Response Status Code:%v", serverErr.StatusCode)
			hlog.Infof("Response Header:%v", serverErr.Header)
			hlog.Infof("Response Err Code:%v", serverErr.Code)
			hlog.Infof("Response Err Msg:%v", serverErr.Message)
		}
		hlog.Fatalf(err.Error())
	}
}

func CheckMinioError(err error) {
	if err != nil {
		hlog.Infof("Object Resource Error:%v", err.Error())
	}
}

type ShardPart struct {
	PartNumber   int       `json:"PartNumber,omitempty"`   // Part编号
	ETag         string    `json:"ETag,omitempty"`         // ETag
	LastModified time.Time `json:"LastModified,omitempty"` // 最后一次修改时间
	Size         int64     `json:"Size,omitempty"`         // Part大小
}
