package listener

import (
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos/enum"
	"sync"
)

type ProgressInfo struct {
	ConsumedBytes int64  `json:"consumed_bytes"`
	TotalBytes    int64  `json:"total_bytes"`
	Percentage    int    `json:"percentage"`
	Status        string `json:"status"` // started, running, success, failed
}

// TosDownLoadListenerService 自定义进度回调，需要实现 tos.DataTransferStatusChange 接口
type TosDownLoadListenerService struct {
	ProgressMap *sync.Map
	DownloadID  string
}

func (l *TosDownLoadListenerService) DataTransferStatusChange(event *tos.DataTransferStatus) {
	val, _ := l.ProgressMap.LoadOrStore(l.DownloadID, &ProgressInfo{})
	progress := val.(*ProgressInfo)

	switch event.Type {
	case enum.DataTransferStarted:
		progress.Status = "started"
	case enum.DataTransferRW:
		progress.ConsumedBytes = event.ConsumedBytes
		progress.TotalBytes = event.TotalBytes
		progress.Status = "running"
		if event.TotalBytes > 0 {
			progress.Percentage = int(event.ConsumedBytes * 100 / event.TotalBytes)
		}
	case enum.DataTransferSucceed:
		progress.Status = "success"
		progress.Percentage = 100
	case enum.DataTransferFailed:
		progress.Status = "failed"
	}
}
