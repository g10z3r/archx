package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/types"
)

type PackageManifest struct {
	StructTypeMap    map[string]*types.StructType    `json:"NodeMap"`
	InterfaceTypeMap map[string]*types.InterfaceType `json:"AbstractMap"`
}

func NewPackageManifest() *PackageManifest {
	return &PackageManifest{
		StructTypeMap:    make(map[string]*types.StructType),
		InterfaceTypeMap: make(map[string]*types.InterfaceType),
	}
}
