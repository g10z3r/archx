package entity

import "time"

type SnapshotEntity struct {
	Timestamp int64
	BasePath  string
}

func NewSnapshotEntity(basePath string) *SnapshotEntity {
	return &SnapshotEntity{
		Timestamp: time.Now().Unix(),
		BasePath:  basePath,
	}
}
