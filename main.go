package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// 支持最大 200M 文件上传
	h := server.Default(server.WithMaxRequestBodySize(200 << 20))

	//service.CheckVolceTosErr(err)
	//
	//// 文件下载
	//h.GET("/download/browser", HandlerDownloadFileFunc)
	//
	//// 查询文件下载进度
	//h.GET("/download/progress", HandlerQueryDownloadStatus)
	//
	//// 解析文件
	//h.POST("/file/parse", HandlerParseFileFunc)

	h.Spin()
}

//func HandlerParseFileFunc(ctx context.Context, c *app.RequestContext) {
//	fileName := c.Query("file_name")
//	data := c.Request.Body()
//	ext := service.GetFileExtension(fileName)
//	parser := service.GetParserByExtension(ext)
//	if parser == nil {
//		c.JSON(http.StatusOK, utils.H{
//			"code":    http.StatusInternalServerError,
//			"message": "不支持当前文件格式",
//			"error":   err.Error(),
//		})
//		return
//	}
//	byData, err := parser.ParseByData(data)
//	if err != nil {
//		c.JSON(http.StatusOK, utils.H{
//			"code":    http.StatusInternalServerError,
//			"message": "文件解析失败",
//			"error":   err.Error(),
//		})
//		return
//	}
//	c.JSON(http.StatusOK, utils.H{
//		"code":    http.StatusOK,
//		"message": "文件解析成功",
//		"data":    byData,
//	})
//}
//
//func HandlerDownloadFileFunc(ctx context.Context, c *app.RequestContext) {
//	fileName := c.Query("file_name")
//	agg := &service.DownloadFileAggregation{
//		Aggregation: &service.DownloadBaseAggregation{
//			Location:   MinioLocation,
//			FileName:   fileName,
//			BucketName: MinioBucketName,
//			Ctx:        ctx,
//		},
//	}
//	err = download.NewMinioObjectDownloadService(minioClient).DownloadFileToBrowser(c, agg)
//	if err != nil {
//		c.JSON(200, utils.H{
//			"code":    http.StatusBadRequest,
//			"message": "下载失败",
//		})
//		return
//	}
//}
//
//// 查询下载进度
//
//func HandlerQueryDownloadStatus(ctx context.Context, c *app.RequestContext) {
//	downloadID := c.Query("file_name")
//	progress := download.NewTosObjectDownloadService(client).GetDownloadProgress(downloadID)
//	if progress == nil {
//		c.JSON(200, utils.H{
//			"code":    http.StatusBadRequest,
//			"message": "获取失败",
//		})
//		return
//	}
//	c.JSON(200, progress)
//}
