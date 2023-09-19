package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

const (
	Embedded    = true
	NotEmbedded = false

	CustomTypeStruct = "struct"
)

type FieldInfo struct {
	_        [0]int
	pos      token.Pos
	end      token.Pos
	Type     string
	Embedded *StructInfo
	IsPublic bool
}

type Usage struct {
	Total int
	Uniq  int
}

type MethodInfo struct {
	Pos      token.Pos
	End      token.Pos
	Usages   map[string]Usage
	IsPublic bool
}

type DependencyInfo struct {
	ElementName string
	Usage       Usage
}

type StructInfo struct {
	_   [0]int
	pos token.Pos
	end token.Pos

	Fields      []*FieldInfo
	FieldsIndex map[string]int

	Methods      []*MethodInfo
	MethodsIndex map[string]int

	Dependencies      []*DependencyInfo
	DependenciesIndex map[string]int

	IsEmbedded bool
}

func (st *StructInfo) AddDependency(importPath, elementName string) {
	index, exists := st.DependenciesIndex[importPath]
	if !exists {
		dependencyInfo := &DependencyInfo{
			ElementName: elementName,
			Usage:       Usage{Total: 1, Uniq: 1},
		}
		st.Dependencies = append(st.Dependencies, dependencyInfo)
		st.DependenciesIndex[importPath] = len(st.Dependencies) - 1
	} else {
		st.Dependencies[index].Usage.Total++
		if st.Dependencies[index].ElementName != elementName {
			st.Dependencies[index].Usage.Uniq++
			st.Dependencies[index].ElementName = elementName
		}
	}
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructInfo, error) {
	fields, fieldsIndex, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	methods := []*MethodInfo{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
		pos:               res.Pos(),
		end:               res.End(),
		Fields:            fields,
		FieldsIndex:       fieldsIndex,
		Methods:           methods,
		MethodsIndex:      methodsIndex,
		Dependencies:      []*DependencyInfo{},
		DependenciesIndex: make(map[string]int),
		IsEmbedded:        isEmbedded,
	}, nil
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) ([]*FieldInfo, map[string]int, error) {
	fields := make([]*FieldInfo, 0, len(fieldList))
	fieldsIndex := make(map[string]int, len(fieldList))

	for i, field := range fieldList {
		fieldTypeString, embedded, err := extractFieldType(fset, field.Type)
		if err != nil {
			return nil, nil, err
		}
		for _, name := range field.Names {
			fieldsIndex[name.Name] = i
			fields = append(fields, &FieldInfo{
				pos:      name.Pos(),
				end:      name.End(),
				Type:     fieldTypeString,
				Embedded: embedded,
				IsPublic: name.IsExported(),
			})
		}
	}
	return fields, fieldsIndex, nil
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (string, *StructInfo, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, err := NewStructType(fset, ft, true)
		if err != nil {
			return "", nil, err
		}
		return CustomTypeStruct, embedded, nil
	case *ast.SelectorExpr:
		return ft.Sel.Name, nil, nil
	default:
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, fieldType); err != nil {
			return "", nil, fmt.Errorf("failed to format node: %w", err)
		}
		return buf.String(), nil, nil
	}
}
