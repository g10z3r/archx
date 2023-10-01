package dao

import (
	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type PackageDAO struct {
	Name              string      `bson:"name"`
	Path              string      `bson:"path"`
	Structs           []StructDAO `bson:"structs"`
	Imports           []string    `bson:"imports"`
	SideEffectImports []int       `bson:"sideEffectImports"`
}

func MapPackageDTO(dto *domainDTO.PackageDTO) PackageDAO {
	return PackageDAO{
		Name:              dto.Name,
		Path:              dto.Path,
		Structs:           make([]StructDAO, 0),
		Imports:           make([]string, 0),
		SideEffectImports: make([]int, 0),
	}
}
