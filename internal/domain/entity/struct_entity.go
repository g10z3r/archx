package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"unicode"
)

const (
	Embedded    = true
	NotEmbedded = false
)

type FieldEntity struct {
	_ [0]int

	Name string

	start token.Pos
	end   token.Pos

	Type     string
	Embedded *StructEntity
	IsPublic bool
}

type MethodEntity struct {
	start token.Pos
	end   token.Pos

	Name         string
	ParentStruct string

	Dependencies map[string]*DependencyEntity

	UsedFields map[string]int
	IsPublic   bool
}

func (s *MethodEntity) AddDependency(importIndex int, element string) {
	if _, exists := s.Dependencies[element]; !exists {
		s.Dependencies[element] = &DependencyEntity{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	s.Dependencies[element].Usage++
}

func NewMethodEntity(res *ast.FuncDecl, parentStructName string) *MethodEntity {
	return &MethodEntity{
		start:        res.Pos(),
		end:          res.End(),
		Name:         res.Name.Name,
		ParentStruct: parentStructName,
		UsedFields:   make(map[string]int),
		Dependencies: make(map[string]*DependencyEntity),
		IsPublic:     unicode.IsUpper(rune(res.Name.Name[0])),
	}
}

type DependencyEntity struct {
	ImportIndex int
	Usage       int
}

type StructEntity struct {
	_            [0]int
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
		fieldMetaData, err := extractFieldType(fset, field.Type)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fieldMetaData.usedPackages); i++ {
			usedPackages = append(usedPackages, fieldMetaData.usedPackages[i])
		}

		if len(field.Names) == 0 {
			fields = append(fields, &FieldEntity{
				Type:     fieldMetaData._type,
				Embedded: fieldMetaData.embeddedStruct,
				IsPublic: false,
			})
			continue
		}

		for _, name := range field.Names {
			fields = append(fields, &FieldEntity{
				start:    name.Pos(),
				end:      name.End(),
				Name:     name.Name,
				Type:     fieldMetaData._type,
				Embedded: fieldMetaData.embeddedStruct,
				IsPublic: name.IsExported(),
			})
		}
	}

	return &fieldMapMetaData{
		fieldsSet:    fields,
		usedPackages: usedPackages,
	}, nil
}

type fieldTypeMetaData struct {
	_type          string
	usedPackages   []UsedPackage
	embeddedStruct *StructEntity
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (*fieldTypeMetaData, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, usedPackages, err := NewStructEntity(fset, ft, true, nil)
		if err != nil {
			return nil, err
		}

		return &fieldTypeMetaData{
			_type:          "struct",
			usedPackages:   usedPackages,
			embeddedStruct: embedded,
		}, nil

	case *ast.SelectorExpr:
		if ident, ok := ft.X.(*ast.Ident); ok {
			return &fieldTypeMetaData{
				_type:        ft.Sel.Name,
				usedPackages: []UsedPackage{{Alias: ident.Name, Element: ft.Sel.Name}},
			}, nil
		}

		return &fieldTypeMetaData{
			_type: ft.Sel.Name,
		}, nil

	default:
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, fieldType); err != nil {
			return nil, fmt.Errorf("failed to format node: %w", err)
		}

		return &fieldTypeMetaData{
			_type: buf.String(),
		}, nil
	}
}
