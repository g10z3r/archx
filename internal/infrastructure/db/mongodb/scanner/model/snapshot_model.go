package model

import (
	"github.com/g10z3r/archx/internal/domain/entity"
)

type SnapshotDAO struct {
	Timestamp     int64          `bson:"timestamp"`
	BasePath      string         `bson:"basePath"`
	Packages      []PackageDAO   `bson:"packages"`
	PackagesIndex map[string]int `bson:"packagesIndex"`
}

func MapSnapshotEntity(e *entity.SnapshotEntity) *SnapshotDAO {
	return &SnapshotDAO{
		Timestamp:     e.Timestamp,
		BasePath:      e.BasePath,
		Packages:      make([]PackageDAO, 0),
		PackagesIndex: make(map[string]int),
	}
}
