package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/entity"
)

type PackageManifest struct {
	StructTypeMap    map[string]*entity.StructType    `json:"NodeMap"`
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
		StructTypeMap:    make(map[string]*entity.StructType),
		InterfaceTypeMap: make(map[string]*entity.InterfaceType),
	}
}
