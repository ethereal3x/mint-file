package main

import (
	"github.com/ethereal3x/mint-file/service"
	"github.com/ethereal3x/mint-file/service/upload"
	"os"
	"testing"
	"time"
)

var (
	minioService *upload.MinioObjectUploadService
)

// Test 获取临时图片Url - Minio
func TestGenerateTempURLBaseMinio(t *testing.T) {

	fileName := "微信图片_20240910155206.jpg"

	g := &service.GenerateBaseAggregation{
		Ctx:           ctx,
		Location:      MinioLocation,
		FileName:      fileName,
		BucketName:    MinioBucketName,
		EffectiveDate: 10 * time.Minute,
	}

	url, err := minioService.GenerateTemporaryLink(g)

	if err != nil {
		t.Fatalf("GeneratePresignedURL failed: %v", err)
	}

	t.Logf("Generated URL: %s", url)
	if url == "" {
		t.Errorf("Expected a valid URL, got empty string")
	}
}

// Test 测试文件上传 - Minio
func TestUploadFileBinaryBaseMinio(t *testing.T) {
	fileName := "1693836967115.png"
	readFile, err := os.ReadFile("C:\\Users\\ASUS\\Pictures\\笔记\\1693836967115.png")

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   MinioLocation,
		FileName:   fileName,
		BucketName: MinioBucketName,
	}

	aggregation := &service.UploadBinaryFileAggregation{Aggregation: uploadBaseAggregation, Data: readFile}
	file, err := minioService.UploadBinaryFile(aggregation)

	checkError(file, err, t)
}

// Test 测试url链接上传 - Minio
func TestUploadUrlFileBaseMinio(t *testing.T) {

	fileName := "123.png"
	url := "https://sfile.chatglm.cn/testpath/8aeab57f-b9c1-5905-864e-164f449f1440_0.png"

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

	file, err := minioService.UploadUrlFile(aggregation)
	checkError(file, err, t)

}

func checkError(obj interface{}, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Test Error: %v", err)
	}

	t.Logf("res: %s", obj)
}
