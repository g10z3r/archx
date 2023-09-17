package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/types"
)

type PackageManifest struct {
	StructTypeMap    map[string]*types.StructType    `json:"NodeMap"`
	InterfaceTypeMap map[string]*types.InterfaceType `json:"AbstractMap"`
}

func (pm *PackageManifest) CountInterfaces() int {
	return len(pm.InterfaceTypeMap)
}

func (pm *PackageManifest) CountStructs() int {
	return len(pm.StructTypeMap)
}

func NewPackageManifest() *PackageManifest {
	return &PackageManifest{
		StructTypeMap:    make(map[string]*types.StructType),
		InterfaceTypeMap: make(map[string]*types.InterfaceType),
	}
}
