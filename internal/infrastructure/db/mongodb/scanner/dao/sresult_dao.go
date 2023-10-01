package dao

type ScannerResultDAO struct {
	Timestamp     int64          `bson:"timestamp"`
	Packages      []PackageDAO   `bson:"packages"`
	PackagesIndex map[string]int `bson:"packagesIndex"`
}
