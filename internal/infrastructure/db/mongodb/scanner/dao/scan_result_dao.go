package dao

import (
	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type ScanResultDAO struct {
	Timestamp     int64          `bson:"timestamp"`
	BasePath      string         `bson:"basePath"`
	Packages      []PackageDAO   `bson:"packages"`
	PackagesIndex map[string]int `bson:"packagesIndex"`
}

func MapScanResultDTO(dto *domainDTO.ScanResultDTO) *ScanResultDAO {
	return &ScanResultDAO{
		Timestamp:     dto.Timestamp,
		BasePath:      dto.BasePath,
		Packages:      make([]PackageDAO, 0),
		PackagesIndex: make(map[string]int),
	}
}
