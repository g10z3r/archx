package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
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
	Element string
	Usage   Usage
}

type DependencyNode struct {
	Prefix   string
	Children map[string]*DependencyNode
	Index    int
}

type DepsRadixTree struct {
	Root     *DependencyNode
	BasePath string
}

func NewDepsRadixTree(basePath string) *DepsRadixTree {
	tree := &DepsRadixTree{
		Root: &DependencyNode{
			Prefix:   basePath,
			Children: make(map[string]*DependencyNode),
			Index:    -1,
		},
		BasePath: basePath,
	}
	return tree
}

func (t *DepsRadixTree) Insert(path string, index int) {
	if strings.HasPrefix(path, t.BasePath) {
		path = strings.TrimPrefix(path, t.BasePath)
	}

	currentNode := t.Root
	pathSegments := strings.Split(path, "/")

	for _, segment := range pathSegments {
		if segment == "" {
			continue // Skip empty segments
		}

		if nextNode, ok := currentNode.Children[segment]; ok {
			currentNode = nextNode
		} else {
			newNode := &DependencyNode{
				Prefix:   segment,
				Children: make(map[string]*DependencyNode),
				Index:    index,
			}
			currentNode.Children[segment] = newNode
			currentNode = newNode
		}
	}
}
func (t *DepsRadixTree) Find(path string) (bool, int) {
	currentNode := t.Root
	pathSegments := strings.Split(path, "/")

	for _, segment := range pathSegments {
		if segment == "" {
			// Skip empty segments that may occur due to double slashes or slashes at the beginning/end of the line
			continue
		}

		if nextNode, ok := currentNode.Children[segment]; ok {
			currentNode = nextNode
		} else {
			return false, -1 // Index -1 means no dependency
		}
	}
	return true, currentNode.Index
}

type StructInfo struct {
	_   [0]int
	pos token.Pos
	end token.Pos

	Fields      []*FieldInfo
	FieldsIndex map[string]int

	Methods      []*MethodInfo
	MethodsIndex map[string]int

	Deps     []*DependencyInfo
	DepsTree *DepsRadixTree

	IsEmbedded bool
}

func (st *StructInfo) AddDependency(importPath, element string) {
	exists, index := st.DepsTree.Find(importPath)
	if !exists {
		dependencyInfo := &DependencyInfo{
			Element: element,
			Usage:   Usage{Total: 1, Uniq: 1},
		}
		st.Deps = append(st.Deps, dependencyInfo)
		st.DepsTree.Insert(importPath, len(st.Deps)-1)
	} else {
		st.Deps[index].Usage.Total++
		if st.Deps[index].Element != element {
			st.Deps[index].Usage.Uniq++
			st.Deps[index].Element = element
		}
	}
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool, modBasePath string, createDepTree bool) (*StructInfo, error) {
	fields, fieldsIndex, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	methods := []*MethodInfo{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
		pos:          res.Pos(),
		end:          res.End(),
		Fields:       fields,
		FieldsIndex:  fieldsIndex,
		Methods:      methods,
		MethodsIndex: methodsIndex,
		Deps:         []*DependencyInfo{},
		DepsTree:     NewDepsRadixTree(modBasePath),
		IsEmbedded:   isEmbedded,
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
		embedded, err := NewStructType(fset, ft, true, "", false)
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
