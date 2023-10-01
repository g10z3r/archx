package dao

import (
	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type FieldDAO struct {
	Type     string     `bson:"type"`
	Embedded *StructDAO `bson:"embedded"`
	IsPublic bool       `bson:"isPublic"`
}

type MethodDAO struct {
	UsedFields map[string]int
	IsPublic   bool
}

type DependencyDAO struct {
	ImportIndex int `bson:"importIndex"`
	Usage       int `bson:"usage"`
}

type StructDAO struct {
	Fields      []*FieldDAO    `bson:"fields"`
	FieldsIndex map[string]int `bson:"fieldsIndex"`

	Methods      []*MethodDAO   `bson:"methods"`
	MethodsIndex map[string]int `bson:"methodsIndex"`

	Dependencies      []*DependencyDAO `bson:"dependencies"`
	DependenciesIndex map[string]int   `bson:"dependenciesIndex"`
}

func MapStructDTO(dto *domainDTO.StructDTO) *StructDAO {
	fields := make([]*FieldDAO, len(dto.Fields))
	for i := 0; i < len(dto.Fields); i++ {
		fields = append(fields, mapFieldDTO(dto.Fields[i]))
	}

	methods := make([]*MethodDAO, len(dto.Methods))
	for i := 0; i < len(dto.Methods); i++ {
		methods = append(methods, mapMethodDTO(dto.Methods[i]))
	}

	deps := make([]*DependencyDAO, len(dto.Dependencies))
	for i := 0; i < len(dto.Dependencies); i++ {
		deps = append(deps, mapDependencyDTO(dto.Dependencies[i]))
	}

	return &StructDAO{
		Fields:            fields,
		FieldsIndex:       dto.FieldsIndex,
		Methods:           methods,
		MethodsIndex:      dto.MethodsIndex,
		Dependencies:      deps,
		DependenciesIndex: dto.DependenciesIndex,
	}
}

func mapFieldDTO(dto *domainDTO.FieldDTO) *FieldDAO {
	var embedded *StructDAO
	if dto.Embedded != nil {
		embedded = MapStructDTO(dto.Embedded)
	}

	return &FieldDAO{
		Type:     dto.Type,
		Embedded: embedded,
		IsPublic: dto.IsPublic,
	}
}

func mapMethodDTO(dto *domainDTO.MethodDTO) *MethodDAO {
	return &MethodDAO{
		UsedFields: dto.UsedFields,
		IsPublic:   dto.IsPublic,
	}
}

func mapDependencyDTO(dto *domainDTO.DependencyDTO) *DependencyDAO {
	return &DependencyDAO{
		ImportIndex: dto.ImportIndex,
		Usage:       dto.Usage,
	}
}
