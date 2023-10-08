package model

import (
	"github.com/g10z3r/archx/internal/domain/obj"
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

	Methods []*MethodDAO `bson:"methods"`

	Dependencies      []*DependencyDAO `bson:"dependencies"`
	DependenciesIndex map[string]int   `bson:"dependenciesIndex"`
}

// func MapMethodEntity(e *entity.MethodEntity) *MethodDAO {
// 	return &MethodDAO{
// 		UsedFields: e.UsedFields,
// 		IsPublic:   e.Visibility,
// 	}
// }

func MapStructEntity(e *obj.StructObj) *StructDAO {
	fields := make([]*FieldDAO, 0, len(e.Fields))
	for i := 0; i < len(e.Fields); i++ {
		fields = append(fields, mapFieldEntity(e.Fields[i]))
	}

	// methods := make([]*MethodDAO, 0, len(e.Methods))
	// for i := 0; i < len(e.Methods); i++ {
	// 	methods = append(methods, mapMethodEntity(e.Methods[i]))
	// }

	// deps := make([]*DependencyDAO, 0, len(e.Dependencies))
	// for i := 0; i < len(e.Dependencies); i++ {
	// 	deps = append(deps, mapDependencyEntity(e.Dependencies[i]))
	// }

	return &StructDAO{
		Fields: fields,
		// FieldsIndex: e.FieldsIndex,
		// Methods:           methods,
		// Dependencies: deps,
		// DependenciesIndex: e.DependenciesIndex,
	}
}

func mapFieldEntity(e *obj.FieldObj) *FieldDAO {
	var embedded *StructDAO
	if e.Embedded != nil {
		embedded = MapStructEntity(e.Embedded)
	}

	return &FieldDAO{
		Type:     e.Type,
		Embedded: embedded,
		IsPublic: e.Visibility,
	}
}

// func mapMethodEntity(e *entity.MethodEntity) *MethodDAO {
// 	return &MethodDAO{
// 		UsedFields: e.UsedFields,
// 		IsPublic:   e.Visibility,
// 	}
// }

func mapDependencyEntity(e *obj.DepObj) *DependencyDAO {
	return &DependencyDAO{
		ImportIndex: e.ImportIndex,
		Usage:       e.Usage,
	}
}
