package types

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
)

const (
	CustomTypeStruct = "struct"
)

type StructType struct {
	_   [0]int
	pos token.Pos
	end token.Pos
	// Fields holds a mapping between field names and their respective metadata
	Fields map[string]*FieldType
	// Methods maps method names to the fields that are utilized within them
	Methods map[string]map[string]FieldUsage
	// Dependencies maps package paths to the names of the types they contain
	Dependencies map[string]map[string]int
	// Flag indicating whether the struct is embedded
	IsEmbedded bool
}

func (st *StructType) AddDependency(importPath string, elementName string) {
	if _, exists := st.Dependencies[importPath]; !exists {
		st.Dependencies[importPath] = make(map[string]int)
	}

	if _, exists := st.Dependencies[importPath][elementName]; !exists {
		st.Dependencies[importPath][elementName] = 1
		return
	}

	st.Dependencies[importPath][elementName]++
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructType, error) {
	fMap, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	methods := make(map[string]map[string]FieldUsage)
	if isEmbedded {
		methods = nil
	}

	return &StructType{
		pos:          res.Pos(),
		end:          res.End(),
		Fields:       fMap,
		Methods:      methods,
		Dependencies: make(map[string]map[string]int),
		IsEmbedded:   isEmbedded,
	}, nil
}
func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (map[string]*FieldType, error) {
	if len(fieldList) < 1 {
		return nil, nil
	}

	fieldMap := make(map[string]*FieldType, len(fieldList))

	for _, field := range fieldList {
		fieldTypeString, embedded, err := extractFieldType(fset, field.Type)
		if err != nil {
			return nil, err
		}

		for _, name := range field.Names {

			fieldMap[name.Name] = &FieldType{
				pos:      name.Pos(),
				end:      name.End(),
				Type:     fieldTypeString,
				Embedded: embedded,
				IsPublic: name.IsExported(),
			}
		}
	}

	return fieldMap, nil
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (string, *StructType, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, err := NewStructType(fset, ft, Embedded)
		if err != nil {
			return "", nil, err
		}
		return CustomTypeStruct, embedded, nil

	case *ast.SelectorExpr:
		return ft.Sel.Name, nil, nil

	default:
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, fieldType); err != nil {
			return "", nil, fmt.Errorf("failed to format node: %v", err)
		}
		return buf.String(), nil, nil
	}
}

type FieldType struct {
	_        [0]int
	pos      token.Pos
	end      token.Pos
	Type     string
	Embedded *StructType
	IsPublic bool
}

type FieldUsage struct {
	Total int
	Uniq  int
}
