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

type FieldObjMeta struct {
	LineCount int
}

type FieldObj struct {
	Name       string
	Type       string
	Embedded   *StructObj
	Visibility bool
	Metadata   *FieldObjMeta
}

type StructObjMeta struct {
	LineCount int
}

type StructObj struct {
	Start        token.Pos
	End          token.Pos
	Name         *string
	Fields       []*FieldObj
	Dependencies map[string]*EntityDepObj
	Incomplete   bool
	Valid        bool
	Metadata     *StructObjMeta
	isEmbedded   bool
}

func (o *StructObj) Type() string {
	return "struct"
}

func (o *StructObj) AddDependency(importIndex int, element string) {
	if _, exists := o.Dependencies[element]; !exists {
		o.Dependencies[element] = &EntityDepObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	o.Dependencies[element].Usage++
}

func NewStructObj(fset *token.FileSet, res *ast.StructType, isEmbedded bool, name *string) (*StructObj, []UsedPackage, error) {
	mapMeta, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	var dependencies map[string]*EntityDepObj
	if !isEmbedded {
		dependencies = make(map[string]*EntityDepObj, len(mapMeta.usedPackages))
	}

	return &StructObj{
			Name:         name,
			Start:        res.Pos(),
			End:          res.End(),
			Fields:       mapMeta.fieldsSet,
			Dependencies: dependencies,
			Incomplete:   res.Incomplete,
			Valid:        res.Struct.IsValid(),
			Metadata: &StructObjMeta{
				LineCount: CalcEntityLOC(fset, res),
			},
			isEmbedded: isEmbedded,
		},
		mapMeta.usedPackages,
		nil
}

type UsedPackage struct {
	_              [0]int
	Alias, Element string
}

type fieldMapMeta struct {
	usedPackages []UsedPackage
	fieldsSet    []*FieldObj
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMeta, error) {
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
				Metadata: &FieldObjMeta{
					LineCount: CalcEntityLOC(fset, name),
				},
			})
		}
	}

	return &fieldMapMeta{
		fieldsSet:    fields,
		usedPackages: usedPackages,
	}, nil
}
