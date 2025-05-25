package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/ethereal3x/mint-file/service"
	"github.com/ethereal3x/mint-file/service/download"
	"github.com/ethereal3x/mint-file/service/upload"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"net/http"
	"strconv"
)

// TOS配置
const (
	TOSAccessKeyID     = "YOUR_TOSAccessKeyID"     // 替换为实际密钥
	TOSSecretAccessKey = "YOUR_TOSSecretAccessKey" // 替换为实际密钥
	TOSRegion          = "YOUR_TOSRegion"          // 替换为实际区域
	TOSBucketName      = "YOUR_TOSBucketName"      // 替换为实际存储桶
	TOSEndpoint        = "YOUR_TOSEndpoint"
	TOSLocation        = "YOUR_TOSLocation"
)

var client *tos.ClientV2
var err error

func main() {
	// 支持最大 200M 文件上传
	h := server.Default(server.WithMaxRequestBodySize(200 << 20))

	// 初始化客户端
	client, err = tos.NewClientV2(TOSEndpoint, tos.WithRegion(TOSRegion), tos.WithCredentials(tos.NewStaticCredentials(TOSAccessKeyID, TOSSecretAccessKey)))

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
	h.GET("/download/progress", HandlerQueryDonwloadStatus)

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

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   TOSLocation,
		FileName:   fileName,
		BucketName: TOSBucketName,
	}

	aggregation := &service.UploadBinaryFileAggregation{Aggregation: uploadBaseAggregation, Data: data}
	file, err := upload.NewTosObjectService(client).UploadBinaryFile(aggregation)
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

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   TOSLocation,
		FileName:   fileName,
		BucketName: TOSBucketName,
	}

	aggregation := &service.UploadUrlFileFileAggregation{
		Aggregation: uploadBaseAggregation,
		Url:         url,
	}

	requestID, err := upload.NewTosObjectService(client).UploadUrlFile(aggregation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.H{
			"code":    http.StatusInternalServerError,
			"message": "上传URL文件到TOS失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.H{
		"code":       http.StatusOK,
		"message":    "上传成功",
		"request_id": requestID,
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
			Location:   TOSLocation,
			FileName:   fileName,
			BucketName: TOSBucketName,
			Ctx:        ctx,
		},
	}
	err = download.NewTosObjectDownloadService(client).DownloadFileToBrowser(c, agg)
	if err != nil {
		c.JSON(200, utils.H{
			"code":    http.StatusBadRequest,
			"message": "下载失败",
		})
		return
	}
}

// 查询下载进度

func HandlerQueryDonwloadStatus(ctx context.Context, c *app.RequestContext) {
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
