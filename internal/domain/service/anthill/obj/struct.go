package obj

import (
	"fmt"
	"go/ast"
	"go/token"
)

const (
	Embedded    = true
	NotEmbedded = false
)

type FieldObjMetadata struct {
	LineCount int
}

type FieldObj struct {
	Name       string
	Type       string
	Embedded   *StructObj
	Visibility bool
	Metadata   *FieldObjMetadata
}

type DepObj struct {
	ImportIndex int
	Usage       int
}

type StructObj struct {
	start        token.Pos
	end          token.Pos
	Name         *string
	Fields       []*FieldObj
	Dependencies map[string]*DepObj
	isEmbedded   bool
}

func (s *StructObj) Type() string {
	return "struct"
}

func (s *StructObj) AddDependency(importIndex int, element string) {
	if _, exists := s.Dependencies[element]; !exists {
		s.Dependencies[element] = &DepObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	s.Dependencies[element].Usage++
}

func NewStructObj(fset *token.FileSet, res *ast.StructType, isEmbedded bool, name *string) (*StructObj, []UsedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	var dependencies map[string]*DepObj
	if !isEmbedded {
		dependencies = make(map[string]*DepObj, len(mapMetaData.usedPackages))
	}

	return &StructObj{
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
	fieldsSet    []*FieldObj
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*FieldObj, 0, len(fieldList))
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
			fields = append(fields, &FieldObj{
				Type:       fieldMetaData.Type,
				Embedded:   fieldMetaData.EmbeddedStruct,
				Visibility: false,
			})
			continue
		}

		for _, name := range field.Names {
			fields = append(fields, &FieldObj{
				Name:       name.Name,
				Type:       fieldMetaData.Type,
				Embedded:   fieldMetaData.EmbeddedStruct,
				Visibility: name.IsExported(),
				Metadata: &FieldObjMetadata{
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
