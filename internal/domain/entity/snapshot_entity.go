package entity

import (
	"time"
)

type SnapshotEntity struct {
	Timestamp int64
	BasePath  string
	Packages  []*PackageEntity
}

func NewSnapshotEntity(mod string, pkgCount int) *SnapshotEntity {
	return &SnapshotEntity{
		Timestamp: time.Now().Unix(),
		BasePath:  mod,
	}
}
