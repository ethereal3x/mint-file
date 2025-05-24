package main

import (
	"github.com/ethereal3x/mint-file/service"
)

type UploadFileService interface {
	UploadBinaryFile(*service.UploadBinaryFileAggregation) (string, error)
	UploadUrlFile(*service.UploadUrlFileFileAggregation) (string, error)
	UploadFragmentFile(*service.UploadFragmentFileAggregation) (string, error)
	ListUploadedFragments(*service.ListUploadedFragmentsAggregation, chan<- *service.ShardPart) error
	CancelFragmentUpload(*service.CancelFragmentUploadAggregation) error
}
