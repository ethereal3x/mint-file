package download

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ethereal3x/mint-file/service"
	"github.com/minio/minio-go/v7"
	"io"
	"path"
)

type MinioObjectDownloadService struct {
	client *minio.Client
}

func NewMinioObjectDownloadService(client *minio.Client) *MinioObjectDownloadService {
	return &MinioObjectDownloadService{
		client: client,
	}
}

func (t *MinioObjectDownloadService) DownloadFileToBrowser(w *app.RequestContext, a *service.DownloadFileAggregation) error {
	if a == nil || a.Aggregation == nil {
		return fmt.Errorf("invalid download aggregation")
	}

	ctx := a.Aggregation.Ctx
	bucket := a.Aggregation.BucketName
	objectName := path.Join(a.Aggregation.Location, a.Aggregation.FileName)

	// 从 MinIO 获取对象流
	object, err := t.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	service.CheckMinioError(err)
	defer func(object *minio.Object) {
		err = object.Close()
	}(object)

	// 获取文件信息（元数据）
	stat, err := object.Stat()
	service.CheckMinioError(err)

	// 设置 HTTP 响应头，告诉浏览器下载文件
	w.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, a.Aggregation.FileName))
	w.Response.Header.Set("Content-Type", stat.ContentType)
	w.Response.Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size))

	// 通过 BodyWriter 写数据流给客户端
	_, err = io.Copy(w.Response.BodyWriter(), object)
	service.CheckMinioError(err)

	return nil
}
