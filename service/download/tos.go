package download

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ethereal3x/mint-file/service"
	"github.com/ethereal3x/mint-file/service/listener"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"io"
	"sync"
)

type TosObjectDownloadService struct {
	client      *tos.ClientV2
	progressMap *sync.Map // 用于跟踪下载进度
}

func NewTosObjectDownloadService(client *tos.ClientV2) *TosObjectDownloadService {
	return &TosObjectDownloadService{
		client:      client,
		progressMap: &sync.Map{},
	}
}

func (t *TosObjectDownloadService) DownloadFileToBrowser(w *app.RequestContext, a *service.DownloadFileAggregation) error {
	fileLocal := service.GenerateDownloadFilePath(a.Aggregation.Location, a.Aggregation.FileName)

	// 安装下载监听器
	getOutput, err := t.client.GetObjectV2(a.Aggregation.Ctx, &tos.GetObjectV2Input{
		Bucket: a.Aggregation.BucketName,
		Key:    fileLocal,
		//DataTransferListener: &listener.TosDownLoadListenerService{
		//	ProgressMap: t.progressMap,
		//	DownloadID:  a,
		//},
	})
	service.CheckVolceTosErr(err)
	defer func(Content io.ReadCloser) {
		_ = Content.Close()
	}(getOutput.Content)

	w.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, a.Aggregation.FileName))
	w.Response.Header.Set("Content-Type", "application/octet-stream")
	w.Response.Header.Set("Content-Length", fmt.Sprintf("%d", getOutput.ContentLength))

	// 设定状态码
	w.Response.SetStatusCode(200)

	// 把流写入响应体
	_, err = io.Copy(w.Response.BodyWriter(), getOutput.Content)
	return err
}

func (t *TosObjectDownloadService) GetDownloadProgress(downloadID string) *listener.ProgressInfo {
	if val, ok := t.progressMap.Load(downloadID); ok {
		return val.(*listener.ProgressInfo)
	}
	return nil
}
