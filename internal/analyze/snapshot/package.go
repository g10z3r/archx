package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/entity"
)

type PackageManifest struct {
	StructTypeMap    map[string]*entity.StructInfo    `json:"NodeMap"`
	InterfaceTypeMap map[string]*entity.InterfaceType `json:"AbstractMap"`
}

func (pm *PackageManifest) CountInterfaces() int {
	return len(pm.InterfaceTypeMap)
}

func (pm *PackageManifest) CountStructs() int {
	return len(pm.StructTypeMap)
}

func NewPackageManifest() *PackageManifest {
	return &PackageManifest{
		StructTypeMap:    make(map[string]*entity.StructInfo),
		InterfaceTypeMap: make(map[string]*entity.InterfaceType),
	}
}
