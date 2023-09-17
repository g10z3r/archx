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

type NodeType map[string]*StructType

type StructType struct {
	_          [0]int
	pos        token.Pos
	end        token.Pos
	Field      map[string]*FieldType
	Method     map[string]map[string]struct{}
	IsEmbedded bool
}

func MakeStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructType, error) {
	fMap, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, err
	}

	var method map[string]map[string]struct{}
	if !isEmbedded {
		method = make(map[string]map[string]struct{})
	}

	return &StructType{
		pos:        res.Pos(),
		end:        res.End(),
		Field:      fMap,
		Method:     method,
		IsEmbedded: isEmbedded,
	}, nil
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (map[string]*FieldType, error) {
	if len(fieldList) < 1 {
		return nil, nil
	}

	var (
		embedded        *StructType
		fieldTypeString string
		err             error
	)

	fieldMap := make(map[string]*FieldType, len(fieldList))
	for _, field := range fieldList {
		// Check that the current field is not a structure, since built-in structures are possible
		switch fieldType := field.Type.(type) {
		case *ast.StructType:
			// Rename here coz otherwise will be just a JSON string of embedded structure
			fieldTypeString = CustomTypeStruct
			embedded, err = MakeStructType(fset, fieldType, Embedded)
			if err != nil {
				return nil, err
			}

		default:
			// Getting the field type
			var buf bytes.Buffer
			if err := format.Node(&buf, fset, field.Type); err != nil {
				return nil, fmt.Errorf("Failed to format node: %v", err)
			}

			fieldTypeString = buf.String()
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

type FieldType struct {
	_        [0]int
	pos      token.Pos
	end      token.Pos
	Type     string
	Embedded *StructType
	IsPublic bool
}

func MakeFieldType(res *ast.Ident) *FieldType {
	return &FieldType{
		pos:      res.Pos(),
		end:      res.End(),
		IsPublic: res.IsExported(),
	}
}
