package obj

type PackageObj struct {
	Name              string
	Path              string
	Imports           []string
	SideEffectImports []string
	Structs           []*StructObj
	Functions         []*FuncObj
}

func NewPackageObj(path, name string) *PackageObj {
	return &PackageObj{
		Path: path,
		Name: name,
	}
}
