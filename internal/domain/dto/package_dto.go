package dto

type PackageDTO struct {
	Name string
	Path string
}

func NewPackageDTO(path, name string) *PackageDTO {
	return &PackageDTO{
		Path: path,
		Name: name,
	}
}
