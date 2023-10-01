package dto

import "time"

type ScanResultDTO struct {
	Timestamp int64
	BasePath  string
}

func NewScanResultDTO(basePath string) *ScanResultDTO {
	return &ScanResultDTO{
		Timestamp: time.Now().Unix(),
		BasePath:  basePath,
	}
}
