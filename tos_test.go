package main

import (
	"context"
	"github.com/ethereal3x/mint-file/service"
	"github.com/ethereal3x/mint-file/service/upload"
	"os"
	"strconv"
	"testing"
	"time"
)

/*
   tos_endpoint: your_tos_endpoint
   tos_access_key: your_tos_access_key
   tos_access_secret: your_tos_access_secret
   tos_region: your_tos_region
   tos_bucket_name: your_tos_bucket_name
*/

var (
	tosService *upload.TosObjectUploadService
)

// Test 获取临时图片Url - Tos
func TestGenerateTempURLBaseTos(t *testing.T) {
	ctx := context.Background()

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	var fileName = "1748098370012326100_123123.png"

	g := &service.GenerateBaseAggregation{
		Ctx:        ctx,
		Location:   TOSLocation,
		FileName:   fileName,
		BucketName: TOSBucketName,
		// 1 到 604800（即 1 秒 到 7 天）
		// Expires: int64((10 * time.Minute).Seconds()),  // 正确
		EffectiveDate: 10 * time.Minute,
	}

	url, err := tosService.GenerateTemporaryLink(g)
	if err != nil {
		t.Fatalf("GeneratePresignedURL failed: %v", err)
	}

	t.Logf("Generated URL: %s", url)
	if url == "" {
		t.Errorf("Expected a valid URL, got empty string")
	}
}

// Test 测试文件上传 - Tos
func TestUploadFileBinaryBaseTos(t *testing.T) {
	fileName := "1693836967115.png"
	readFile, err := os.ReadFile("C:\\Users\\ASUS\\Pictures\\笔记\\1693836967115.png")

	uploadBaseAggregation := &service.UploadBaseAggregation{
		Ctx:        ctx,
		Location:   TOSLocation,
		FileName:   fileName,
		BucketName: TOSBucketName,
	}

	aggregation := &service.UploadBinaryFileAggregation{Aggregation: uploadBaseAggregation, Data: readFile}
	file, err := tosService.UploadBinaryFile(aggregation)

	checkError(file, err, t)
}

// Test 测试url链接上传 - Minio
func TestUploadUrlFileBaseTos(t *testing.T) {

	fileName := "123.png"
	url := "https://sfile.chatglm.cn/testpath/8aeab57f-b9c1-5905-864e-164f449f1440_0.png"

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

	file, err := tosService.UploadUrlFile(aggregation)
	checkError(file, err, t)

}

func TestUploadShardFileBaseTos(t *testing.T) {
	fileName := "VALORANT-2025-05-01-22-23-24.mp4"
	shardSizeStr := "5242880"

	readFile, err := os.ReadFile("C:\\Users\\ASUS\\Videos\\Captures\\VALORANT.mp4")

	shardSize, err := strconv.ParseInt(shardSizeStr, 10, 64)

	checkError("", err, t)

	data := readFile

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

	requestID, err := tosService.UploadFragmentFile(aggregation)
	checkError(requestID, err, t)
}
