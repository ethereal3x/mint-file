package upload

import (
	"bytes"
	"fmt"
	"github.com/ethereal3x/mint-file/service"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"path"
)

type MinioObjectUploadService struct {
	client *minio.Client
}

func NewMinioObjectUploadService(client *minio.Client) *MinioObjectUploadService {
	return &MinioObjectUploadService{client: client}
}

func (m MinioObjectUploadService) UploadBinaryFile(aggregation *service.UploadBinaryFileAggregation) (string, error) {

	ctx := aggregation.Aggregation.Ctx
	bucket := aggregation.Aggregation.BucketName
	objectName := path.Join(aggregation.Aggregation.Location, aggregation.Aggregation.FileName)
	data := aggregation.Data

	contentType := "application/octet-stream"
	reader := bytes.NewReader(data)
	size := int64(len(data))

	// 上传文件
	_, err := m.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	service.CheckMinioError(err)

	// 构造返回路径，格式为 bucket/object
	return fmt.Sprintf("%s/%s", bucket, objectName), nil
}

func (m MinioObjectUploadService) UploadUrlFile(aggregation *service.UploadUrlFileFileAggregation) (string, error) {
	if aggregation == nil || aggregation.Aggregation == nil {
		return "", fmt.Errorf("invalid upload aggregation")
	}

	ctx := aggregation.Aggregation.Ctx
	bucket := aggregation.Aggregation.BucketName
	objectName := path.Join(aggregation.Aggregation.Location, aggregation.Aggregation.FileName)
	fileUrl := aggregation.Url

	// 下载 URL 文件内容
	resp, _ := http.Get(fileUrl)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 读取数据到 buffer
	data, err := io.ReadAll(resp.Body)
	service.CheckMinioError(err)
	size := int64(len(data))
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传到 MinIO
	reader := bytes.NewReader(data)
	_, err = m.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	service.CheckMinioError(err)

	return fmt.Sprintf("%s/%s", bucket, objectName), nil
}

func (m MinioObjectUploadService) UploadFragmentFile(aggregation *service.UploadFragmentFileAggregation) (string, error) {
	ctx := aggregation.Aggregation.Ctx
	bucket := aggregation.Aggregation.BucketName
	objectName := path.Join(aggregation.Aggregation.Location, aggregation.Aggregation.FileName)
	data := aggregation.Data

	// 默认 ContentType 和读入数据
	reader := bytes.NewReader(data)
	size := int64(len(data))

	// 只要数据大小超过 5MB，它就会自动使用 multipart
	_, err := m.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return fmt.Sprintf("%s/%s", bucket, objectName), nil

}

// ListUploadedFragments is Deprecated: 该方法暂无实现
func (m MinioObjectUploadService) ListUploadedFragments(aggregation *service.ListUploadedFragmentsAggregation, parts chan<- *service.ShardPart) error {
	// TODO
	panic("implement me")
}

// CancelFragmentUpload is Deprecated: 该方法暂无实现
func (m MinioObjectUploadService) CancelFragmentUpload(aggregation *service.CancelFragmentUploadAggregation) error {
	//TODO implement me
	panic("implement me")
}
