package main

import (
	"github.com/ethereal3x/mint-file/service"
	"net/http"
)

type DownLoadFileService interface {
	DownloadFileToBrowser(w http.ResponseWriter, s *service.DownloadFileAggregation) error
}
