package entity

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	Embedded    = true
	NotEmbedded = false
)

type FieldMetadata struct {
	LineCount int
}

type FieldEntity struct {
	Name       string
	Type       string
	Embedded   *StructEntity
	Visibility bool
	Metadata   *FieldMetadata
}

type DependencyEntity struct {
	ImportIndex int
	Usage       int
}

type StructEntity struct {
	start        token.Pos
	end          token.Pos
	Name         *string
	Fields       []*FieldEntity
	Dependencies map[string]*DependencyEntity
	isEmbedded   bool
}

func (s *StructEntity) AddDependency(importIndex int, element string) {
	if _, exists := s.Dependencies[element]; !exists {
		s.Dependencies[element] = &DependencyEntity{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	s.Dependencies[element].Usage++
}

func NewStructEntity(fset *token.FileSet, res *ast.StructType, isEmbedded bool, name *string) (*StructEntity, []UsedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	var dependencies map[string]*DependencyEntity
	if !isEmbedded {
		dependencies = make(map[string]*DependencyEntity, len(mapMetaData.usedPackages))
	}

	return &StructEntity{
			Name:         name,
			start:        res.Pos(),
			end:          res.End(),
			Fields:       mapMetaData.fieldsSet,
			Dependencies: dependencies,
			isEmbedded:   isEmbedded,
		},
		mapMetaData.usedPackages,
		nil
}

type UsedPackage struct {
	_              [0]int
	Alias, Element string
}

type fieldMapMetaData struct {
	usedPackages []UsedPackage
	fieldsSet    []*FieldEntity
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*FieldEntity, 0, len(fieldList))
	usedPackages := make([]UsedPackage, 0, len(fieldList))

	for _, field := range fieldList {
		fieldMetaData, err := ExtractExprAsType(fset, field.Type)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fieldMetaData.UsedPackages); i++ {
			usedPackages = append(usedPackages, fieldMetaData.UsedPackages[i])
		}

		if len(field.Names) == 0 {
			fields = append(fields, &FieldEntity{
				Type:       fieldMetaData.Type,
				Embedded:   fieldMetaData.EmbeddedStruct,
				Visibility: false,
			})
			continue
		}

		for _, name := range field.Names {
			fields = append(fields, &FieldEntity{
				Name:       name.Name,
				Type:       fieldMetaData.Type,
				Embedded:   fieldMetaData.EmbeddedStruct,
				Visibility: name.IsExported(),
				Metadata: &FieldMetadata{
					LineCount: calcLineCount(fset, name),
				},
			})
		}
	}

	return &fieldMapMetaData{
		fieldsSet:    fields,
		usedPackages: usedPackages,
	}, nil
}
