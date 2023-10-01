package dto

type PackageDTO struct {
	Name              string
	Path              string
	Structs           []StructDTO
	Imports           []string
	SideEffectImports []int
}
