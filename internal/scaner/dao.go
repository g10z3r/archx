package scaner

import (
	"github.com/g10z3r/archx/internal/scaner/entity"
)

type Package struct {
	Name              string          `bson:"name"`
	Path              string          `bson:"path"`
	Structs           []entity.Struct `bson:"structs"`
	Imports           []string        `bson:"imports"`
	SideEffectImports []int           `bson:"sideEffectImports"`
}

type ScanResult struct {
	Timestamp     int64          `bson:"timestamp"`
	Packages      []Package      `bson:"packages"`
	PackagesIndex map[string]int `bson:"packagesIndex"`
}
