package dto

type ScannerResultDTO struct {
	Timestamp     int64
	Packages      []PackageDTO
	PackagesIndex map[string]int
}
