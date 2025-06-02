package main

import (
	"context"
	"github.com/ethereal3x/mint-file/service/upload"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"os"
	"testing"
)

const (
	TOSAccessKeyID     = "A*********************************************U"
	TOSSecretAccessKey = "T**************************************************="
	TOSRegion          = "cn-guangzhou"
	TOSBucketName      = "mint-test"
	TOSEndpoint        = "https://tos-cn-guangzhou.volces.com"
	TOSLocation        = "/test"

	MinioAccessKey  = "o******************M"
	MinioSecretKey  = "v*********************************************t"
	MinioEndPoint   = "*.*.*.*:9000"
	MinioBucketName = "test"
	MinioLocation   = ""
)

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	client, err := tos.NewClientV2(TOSEndpoint,
		tos.WithRegion(TOSRegion),
		tos.WithCredentials(tos.NewStaticCredentials(TOSAccessKeyID, TOSSecretAccessKey)))
	if err != nil {
		panic("初始化 TOS 客户端失败: " + err.Error())
	}

	// 初始化 tosService
	tosService = upload.NewTosObjectService(client)

	minioClient, err = minio.New(MinioEndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
		Secure: false,
	})

	// 初始化 minioService
	minioService = upload.NewMinioObjectUploadService(minioClient)

	// 启动测试
	code := m.Run()
	os.Exit(code)
}
