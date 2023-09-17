package snapshot

import "github.com/g10z3r/archx/internal/analyze/types"

type PackageManifest struct {
	StructTypeMap map[string]*types.StructType `json:"NodeMap"`
}

func NewPackageManifest() *PackageManifest {
	return &PackageManifest{
		StructTypeMap: make(map[string]*types.StructType),
	}
}
