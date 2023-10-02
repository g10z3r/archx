package entity

type PackageEntity struct {
	Name string
	Path string
}

func NewPackageEntity(path, name string) *PackageEntity {
	return &PackageEntity{
		Path: path,
		Name: name,
	}
}
