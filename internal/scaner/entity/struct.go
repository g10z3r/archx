package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"unicode"
)

const (
	Embedded    = true
	NotEmbedded = false

	onlyPreinitialized = false

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

type Method struct {
	Start      token.Pos
	End        token.Pos
	UsedFields map[string]Usage
	IsPublic   bool
}

func NewMethod(res *ast.FuncDecl) *Method {
	return &Method{
		Start:      res.Pos(),
		End:        res.End(),
		UsedFields: make(map[string]Usage),
		IsPublic:   unicode.IsUpper(rune(res.Name.Name[0])),
	}
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
	Root *DependencyNode
}

func NewDepsRadixTree() *DepsRadixTree {
	tree := &DepsRadixTree{
		Root: &DependencyNode{
			Children: make(map[string]*DependencyNode),
			Index:    -1,
		},
	}
	return tree
}

func (t *DepsRadixTree) Insert(path string, index int) {
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
			continue
		}

		if nextNode, ok := currentNode.Children[segment]; ok {
			currentNode = nextNode
		} else {
			return false, -1
		}
	}
	return true, currentNode.Index
}

func (t *DepsRadixTree) Compress() {
	compress(t.Root)
}

func compress(node *DependencyNode) {
	for prefix, child := range node.Children {
		compress(child)

		if len(child.Children) == 1 {
			for childPrefix, grandChild := range child.Children {
				node.Children[prefix+"/"+childPrefix] = grandChild
				delete(node.Children, prefix)
			}
		}
	}
}

type StructInfo struct {
	_   [0]int
	Pos token.Pos
	End token.Pos

	Fields      []*FieldInfo
	FieldsIndex map[string]int

	Methods      []*Method
	MethodsIndex map[string]int

	// dependencies
	Deps     []*DependencyInfo
	DepsTree *DepsRadixTree

	Incomplete bool
	isEmbedded bool
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

func (si *StructInfo) AddMethod(metdod *Method, name string) {
	si.Methods = append(si.Methods, metdod)
	si.MethodsIndex[name] = len(si.Methods) - 1
}

func (si *StructInfo) SyncMethods(from *StructInfo) {
	for methodName, i := range from.MethodsIndex {
		if _, exists := si.MethodsIndex[methodName]; !exists {
			fmt.Println(from.MethodsIndex)
			fmt.Println(si.MethodsIndex)
			si.AddMethod(from.Methods[i], methodName)
		}
	}
}

func NewStructPreInit(name string) *StructInfo {
	methods := []*Method{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
		Methods:      methods,
		MethodsIndex: methodsIndex,
		Deps:         []*DependencyInfo{},
		DepsTree:     NewDepsRadixTree(),
		isEmbedded:   NotEmbedded,
		Incomplete:   onlyPreinitialized,
	}
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructInfo, []usedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	methods := []*Method{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
			Pos:          res.Pos(),
			End:          res.End(),
			Fields:       mapMetaData.fieldsSet,
			FieldsIndex:  mapMetaData.fieldsIndex,
			Methods:      methods,
			MethodsIndex: methodsIndex,
			Deps:         []*DependencyInfo{},
			DepsTree:     NewDepsRadixTree(),
			isEmbedded:   isEmbedded,
			Incomplete:   true,
		},
		mapMetaData.usedPackages,
		nil
}

type usedPackage struct {
	_              [0]int
	Alias, Element string
}

type fieldMapMetaData struct {
	usedPackages []usedPackage
	fieldsSet    []*FieldInfo
	fieldsIndex  map[string]int
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*FieldInfo, 0, len(fieldList))
	fieldsIndex := make(map[string]int, len(fieldList))
	usedPackages := []usedPackage{}

	for i, field := range fieldList {
		fieldMetaData, err := extractFieldType(fset, field.Type)
		if err != nil {
			return nil, err
		}

		if fieldMetaData.isImported {
			for i := 0; i < len(fieldMetaData.usedPackages); i++ {
				usedPackages = append(usedPackages, fieldMetaData.usedPackages[i])
			}
		}

		for _, name := range field.Names {
			fieldsIndex[name.Name] = i
			fields = append(fields, &FieldInfo{
				pos:      name.Pos(),
				end:      name.End(),
				Type:     fieldMetaData._type,
				Embedded: fieldMetaData.embeddedStruct,
				IsPublic: name.IsExported(),
			})
		}
	}

	return &fieldMapMetaData{
		fieldsSet:    fields,
		usedPackages: usedPackages,
		fieldsIndex:  fieldsIndex,
	}, nil
}

type fieldTypeMetaData struct {
	_type          string
	usedPackages   []usedPackage
	isImported     bool
	embeddedStruct *StructInfo
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (*fieldTypeMetaData, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, usedPackages, err := NewStructType(fset, ft, true)
		if err != nil {
			return nil, err
		}

		return &fieldTypeMetaData{
			_type:          CustomTypeStruct,
			usedPackages:   usedPackages,
			embeddedStruct: embedded,
		}, nil

	case *ast.SelectorExpr:
		if ident, ok := ft.X.(*ast.Ident); ok {
			return &fieldTypeMetaData{
				_type:        ft.Sel.Name,
				usedPackages: []usedPackage{{Alias: ident.Name, Element: ft.Sel.Name}},
				isImported:   true,
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
