package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/entity"
)

type PackageManifest struct {
	Structs      []*entity.StructInfo
	StructsIndex map[string]int

	Interfaces      []*entity.InterfaceType
	InterfacesIndex map[string]int
}

func (pm *PackageManifest) CountInterfaces() int {
	return len(pm.Interfaces)
}

func (pm *PackageManifest) CountStructs() int {
	return len(pm.Structs)
}

func NewPackageManifest() *PackageManifest {
	return &PackageManifest{
		Structs:         []*entity.StructInfo{},
		StructsIndex:    make(map[string]int),
		Interfaces:      []*entity.InterfaceType{},
		InterfacesIndex: make(map[string]int),
	}
}
