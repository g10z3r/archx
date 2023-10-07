package entity

type PackageEntity struct {
	Name              string
	Path              string
	Imports           []string
	SideEffectImports []string
	Structs           []*StructEntity
	Functions         []*FunctionEntity
}

func NewPackageEntity(path, name string) *PackageEntity {
	return &PackageEntity{
		Path: path,
		Name: name,
	}
}
