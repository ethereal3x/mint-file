package upload

import (
	"bytes"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/ethereal3x/mint-file/service"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"
	"io"
	"net/http"
)

type TosObjectUploadService struct {
	client *tos.ClientV2
}

func NewTosObjectService(client *tos.ClientV2) *TosObjectUploadService {
	return &TosObjectUploadService{client: client}
}

func (t *TosObjectUploadService) UploadBinaryFile(a *service.UploadBinaryFileAggregation) (string, error) {
	fileLocal := service.GenerateUniqueFilePath(a.Aggregation.Location, a.Aggregation.FileName)
	output, err := t.client.PutObjectV2(a.Aggregation.Ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: a.Aggregation.BucketName,
			Key:    fileLocal,
		},
		Content: bytes.NewReader(a.Data),
	})
	service.CheckVolceTosErr(err)
	return output.RequestID, nil
}

func (t *TosObjectUploadService) UploadUrlFile(a *service.UploadUrlFileFileAggregation) (string, error) {
	fileLocal := service.GenerateUniqueFilePath(a.Aggregation.Location, a.Aggregation.FileName)

	resp, _ := http.Get(a.Url)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	output, err := t.client.PutObjectV2(a.Aggregation.Ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket: a.Aggregation.BucketName,
			Key:    fileLocal,
		},
		Content: resp.Body,
	})
	service.CheckVolceTosErr(err)
	return output.RequestID, nil
}

func (t *TosObjectUploadService) UploadFragmentFile(a *service.UploadFragmentFileAggregation) (string, error) {
	fileLocal := service.GenerateUniqueFilePath(a.Aggregation.Location, a.Aggregation.FileName)

	// 1. 初始化分片上传任务
	createMultipartOutput, err := t.client.CreateMultipartUploadV2(a.Aggregation.Ctx, &tos.CreateMultipartUploadV2Input{
		Bucket:       a.Aggregation.BucketName,
		Key:          fileLocal,
		ACL:          enum.ACLPublicRead,
		StorageClass: enum.StorageClassIa,
		Meta:         map[string]string{"key": "value"},
	})
	service.CheckVolceTosErr(err)

	uploadID := createMultipartOutput.UploadID
	totalSize := int64(len(a.Data))
	var parts []tos.UploadedPartV2

	// 2. 分片上传
	for i := int64(0); i < totalSize; i += a.ShardSize {
		end := i + a.ShardSize
		if end > totalSize {
			end = totalSize
		}

		partNumber := int(i/a.ShardSize) + 1
		partReader := bytes.NewReader(a.Data[i:end])

		uploadPartOutput, err := t.client.UploadPartV2(a.Aggregation.Ctx, &tos.UploadPartV2Input{
			UploadPartBasicInput: tos.UploadPartBasicInput{
				Bucket:     a.Aggregation.BucketName,
				Key:        fileLocal,
				UploadID:   uploadID,
				PartNumber: partNumber,
			},
			Content:       partReader,
			ContentLength: end - i,
		})
		service.CheckVolceTosErr(err)

		parts = append(parts, tos.UploadedPartV2{
			ETag:       uploadPartOutput.ETag,
			PartNumber: partNumber,
		})
	}

	// 3. 完成上传（所有分片都上传后）
	completeOutput, err := t.client.CompleteMultipartUploadV2(a.Aggregation.Ctx, &tos.CompleteMultipartUploadV2Input{
		Bucket:   a.Aggregation.BucketName,
		Key:      fileLocal,
		UploadID: uploadID,
		Parts:    parts,
	})
	service.CheckVolceTosErr(err)
	return completeOutput.RequestID, nil
}

func (t *TosObjectUploadService) ListUploadedFragments(a *service.ListUploadedFragmentsAggregation, c chan<- *service.ShardPart) error {
	// 生成对象在 TOS 中的完整路径（Key），用于唯一标识对象
	fileLocal := service.GenerateUniqueFilePath(a.Aggregation.Location, a.Aggregation.FileName)

	// 用于控制分页读取结果
	truncated := true
	marker := 0

	// 循环分页获取所有已上传的分片信息
	for truncated {
		// 请求列出当前 UploadID 下的分片（从 marker 开始）
		output, err := t.client.ListParts(a.Aggregation.Ctx, &tos.ListPartsInput{
			Bucket:           a.Aggregation.BucketName, // Bucket 名称
			Key:              fileLocal,                // 对象在 TOS 中的路径
			UploadID:         a.UploadID,               // 当前分片上传的 UploadID
			PartNumberMarker: marker,                   // 上一页最后一个 PartNumber，第一次为 0
		})
		service.CheckVolceTosErr(err)

		// 是否还有下一页结果（如果为 true，则继续循环）
		truncated = output.IsTruncated

		// 设置下一页起始位置
		marker = output.NextPartNumberMarker

		// 遍历当前页的所有分片信息
		for _, part := range output.Parts {
			// 打印分片编号（从 1 开始递增）
			hlog.Infof("Part Number: %v", part.PartNumber)

			// 打印分片的 ETag（上传后由 TOS 返回的校验值）
			hlog.Infof("ETag: %v", part.ETag)

			// 打印分片的大小（单位：字节）
			hlog.Infof("Size: %v", part.Size)

			for _, part := range output.Parts {
				select {
				case <-a.Aggregation.Ctx.Done():
					return a.Aggregation.Ctx.Err() // 支持中断
				case c <- &service.ShardPart{
					PartNumber:   part.PartNumber,
					ETag:         part.ETag,
					LastModified: part.LastModified,
					Size:         part.Size,
				}:
				}
			}
		}
	}
	return nil
}

func (t *TosObjectUploadService) CancelFragmentUpload(a *service.CancelFragmentUploadAggregation) {
	fileLocal := service.GenerateUniqueFilePath(a.Aggregation.Location, a.Aggregation.FileName)
	// 取消分片上传
	_, err := t.client.AbortMultipartUpload(a.Aggregation.Ctx, &tos.AbortMultipartUploadInput{
		Bucket:   a.Aggregation.BucketName,
		Key:      fileLocal,
		UploadID: a.UploadID,
	})
	service.CheckVolceTosErr(err)
}
