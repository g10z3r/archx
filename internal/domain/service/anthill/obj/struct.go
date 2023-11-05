package obj

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

const (
	Embedded    = true
	NotEmbedded = false
)

type StructObjMeta struct {
	LineCount int
}

type StructTypeObj struct {
	Name         *string
	Fields       []*FieldObj
	Dependencies map[string]*DependencyObj
	Incomplete   bool
	Valid        bool
	Metadata     *StructObjMeta
}

func (o *StructTypeObj) Type() string {
	return "struct"
}

func (o *StructTypeObj) AddDependency(importIndex int, element string) {
	if _, exists := o.Dependencies[element]; !exists {
		o.Dependencies[element] = &DependencyObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	o.Dependencies[element].Usage++
}

func NewStructObj(fset *token.FileSet, node ast.Node, name *string) (*StructTypeObj, []UsedPackage, error) {
	var structType *ast.StructType

	typeSpec, ok := node.(*ast.TypeSpec)
	if ok {
		structType, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return nil, nil, errors.New("some error from NewStructObj 1") // TODO: add normal error return message
		}
	} else {
		structType, ok = node.(*ast.StructType)
		if !ok {
			return nil, nil, errors.New("some error from NewStructObj 2") // TODO: add normal error return message
		}
	}

	extractedFieldsData, err := extractFieldMap(fset, structType.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	return &StructTypeObj{
			Name:         name,
			Fields:       extractedFieldsData.fieldsSet,
			Dependencies: make(map[string]*DependencyObj, len(extractedFieldsData.usedPackages)),
			Incomplete:   structType.Incomplete,
			Valid:        structType.Struct.IsValid(),
			Metadata: &StructObjMeta{
				LineCount: CalcEntityLOC(fset, structType),
			},
		},
		extractedFieldsData.usedPackages,
		nil
}

type UsedPackage struct {
	_              [0]int
	Alias, Element string
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*extractedFieldsData, error) {
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
			})
		}
	}

	return &extractedFieldsData{
		fieldsSet:    fields,
		usedPackages: usedPackages,
	}, nil
}
