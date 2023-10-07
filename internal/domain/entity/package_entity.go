package entity

type PackageEntity struct {
	Name              string
	Path              string
	Imports           []string
	SideEffectImports []string
	Structs           []*StructEntity
	Methods           []*MethodEntity
}

func NewPackageEntity(path, name string) *PackageEntity {
	return &PackageEntity{
		Path: path,
		Name: name,
	}
}
