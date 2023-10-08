package model

import (
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type PackageDAO struct {
	Name string `bson:"name"`
	Path string `bson:"path"`

	Structs      []StructDAO    `bson:"structs"`
	StructsIndex map[string]int `bson:"structsIndex"`

	Imports           []string `bson:"imports"`
	SideEffectImports []int    `bson:"sideEffectImports"`
}

func MapPackageEntity(e *obj.PackageObj) PackageDAO {
	return PackageDAO{
		Name: e.Name,
		Path: e.Path,

		Structs:      make([]StructDAO, 0),
		StructsIndex: make(map[string]int),

		Imports:           make([]string, 0),
		SideEffectImports: make([]int, 0),
	}
}
