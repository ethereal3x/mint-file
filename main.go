package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/ethereal3x/mint-file/service"
	"github.com/ethereal3x/mint-file/service/download"
	"github.com/ethereal3x/mint-file/service/upload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"net/http"
	"strconv"
)

var client *tos.ClientV2
var minioClient *minio.Client
var err error

// TOS配置
const (
	TOSAccessKeyID     = "YOUR_TOSAccessKeyID"     // 替换为实际密钥
	TOSSecretAccessKey = "YOUR_TOSSecretAccessKey" // 替换为实际密钥
	TOSRegion          = "YOUR_TOSRegion"          // 替换为实际区域
	TOSBucketName      = "YOUR_TOSBucketName"      // 替换为实际存储桶
	TOSEndpoint        = "YOUR_TOSEndpoint"
	TOSLocation        = "YOUR_TOSLocation"
)

const (
	MinioAccessKey  = "ogLR75LInO4769GebWXK"
	MinioSecretKey  = "2kHSuK5ISaNCabKmgUU42cGoWXpzB1viqMq9fIj9"
	MinioEndPoint   = "127.0.0.1:9000"
	MinioBucketName = "mint"
	MinioLocation   = "test"
)

func main() {
	// 支持最大 200M 文件上传
	h := server.Default(server.WithMaxRequestBodySize(200 << 20))

	// 初始化tos客户端
	client, err = tos.NewClientV2(TOSEndpoint, tos.WithRegion(TOSRegion), tos.WithCredentials(tos.NewStaticCredentials(TOSAccessKeyID, TOSSecretAccessKey)))

	// 初始话minio客户端
	// 使用 accessKey 和 secretKey 初始化 minio.Client
	minioClient, err = minio.New(MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false,
	})

	service.CheckVolceTosErr(err)

	// 上传文件
	h.POST("/upload", HandlerUploadBinaryFunc)

	// 上传 URL 文件
	h.POST("/upload/url", HandlerUploadUrlFile)

	// 分片上传
	h.POST("/upload/fragment", HandleShardFileFunc)

	// 文件下载
	h.GET("/download/browser", HandlerDownloadFileFunc)

	// 查询文件下载进度
	h.GET("/download/progress", HandlerQueryDownloadStatus)

	// 解析文件
	h.POST("/file/parse", HandlerParseFileFunc)

	h.Spin()
}

func HandlerParseFileFunc(ctx context.Context, c *app.RequestContext) {
	fileName := c.Query("file_name")
	data := c.Request.Body()
	ext := service.GetFileExtension(fileName)
	parser := service.GetParserByExtension(ext)
	if parser == nil {
		c.JSON(http.StatusOK, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "不支持当前文件格式",
			"error":   err.Error(),
		})
		return
	}
	byData, err := parser.ParseByData(data)
	if err != nil {
		c.JSON(http.StatusOK, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "文件解析失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.H{
		"code":    http.StatusOK,
		"message": "文件解析成功",
		"data":    byData,
	})
}

func HandlerUploadBinaryFunc(ctx context.Context, c *app.RequestContext) {
	fileName := c.Query("file_name")
	data := c.Request.Body()

	//uploadBaseAggregation := &service.UploadBaseAggregation{
	//	Ctx:        ctx,
	//	Location:   TOSLocation,
	//	FileName:   fileName,
	//	BucketName: TOSBucketName,
	//}

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   MinioLocation,
		FileName:   fileName,
		BucketName: MinioBucketName,
	}

	aggregation := &service.UploadBinaryFileAggregation{Aggregation: uploadBaseAggregation, Data: data}
	// file, err := upload.NewTosObjectService(client).UploadBinaryFile(aggregation)
	file, err := upload.NewMinioObjectUploadService(minioClient).UploadBinaryFile(aggregation)
	if err != nil {
		c.JSON(http.StatusOK, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "上传到TOS失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":       http.StatusOK,
		"message":    "上传成功",
		"request_id": file,
	})
}

func HandlerUploadUrlFile(ctx context.Context, c *app.RequestContext) {
	url := c.Query("url")
	fileName := c.Query("file_name")

	//uploadBaseAggregation := &service.UploadBaseAggregation{
	//	Ctx:        ctx,
	//	Location:   TOSLocation,
	//	FileName:   fileName,
	//	BucketName: TOSBucketName,
	//}

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   MinioLocation,
		FileName:   fileName,
		BucketName: MinioBucketName,
	}

	aggregation := &service.UploadUrlFileFileAggregation{
		Aggregation: uploadBaseAggregation,
		Url:         url,
	}

	// requestID, err := upload.NewTosObjectService(client).UploadUrlFile(aggregation)
	file, err := upload.NewMinioObjectUploadService(minioClient).UploadUrlFile(aggregation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "上传URL文件到TOS失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":    http.StatusOK,
		"message": "上传成功",
		// "request_id": requestID,
		"path": file,
	})
}

func HandleShardFileFunc(ctx context.Context, c *app.RequestContext) {
	fileName := c.Query("file_name")
	shardSizeStr := c.Query("shard_size")

	if fileName == "" || shardSizeStr == "" {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "file_name 和 shard_size 是必需参数",
		})
		return
	}

	shardSize, err := strconv.ParseInt(shardSizeStr, 10, 64)
	if err != nil || shardSize <= 0 {
		c.JSON(http.StatusBadRequest, utils.H{
			"code":    http.StatusBadRequest,
			"message": "shard_size 参数无效",
		})
		return
	}

	data := c.Request.Body()

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   TOSLocation,
		FileName:   fileName,
		BucketName: TOSBucketName,
	}

	aggregation := &service.UploadFragmentFileAggregation{
		Aggregation: uploadBaseAggregation,
		ShardSize:   shardSize,
		Data:        data,
	}

	requestID, err := upload.NewTosObjectService(client).UploadFragmentFile(aggregation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "分片上传失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":       http.StatusOK,
		"message":    "分片上传成功",
		"request_id": requestID,
	})
}

func HandlerDownloadFileFunc(ctx context.Context, c *app.RequestContext) {
	fileName := c.Query("file_name")
	agg := &service.DownloadFileAggregation{
		Aggregation: &service.DownloadBaseAggregation{
			Location:   MinioLocation,
			FileName:   fileName,
			BucketName: MinioBucketName,
			Ctx:        ctx,
		},
	}
	err = download.NewMinioObjectDownloadService(minioClient).DownloadFileToBrowser(c, agg)
	if err != nil {
		c.JSON(200, utils.H{
			"code":    http.StatusBadRequest,
			"message": "下载失败",
		})
		return
	}
}

// 查询下载进度

func HandlerQueryDownloadStatus(ctx context.Context, c *app.RequestContext) {
	downloadID := c.Query("file_name")
	progress := download.NewTosObjectDownloadService(client).GetDownloadProgress(downloadID)
	if progress == nil {
		c.JSON(200, utils.H{
			"code":    http.StatusBadRequest,
			"message": "获取失败",
		})
		return
	}
	c.JSON(200, progress)
}
