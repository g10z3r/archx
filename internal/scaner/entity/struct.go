package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"sync"
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
	UsedFields map[string]int
	IsPublic   bool
}

func NewMethod(res *ast.FuncDecl) *Method {
	return &Method{
		Start:      res.Pos(),
		End:        res.End(),
		UsedFields: make(map[string]int),
		IsPublic:   unicode.IsUpper(rune(res.Name.Name[0])),
	}
}

type DependencyInfo struct {
	ImportIndex int
	Usage       int
}

type StructInfo struct {
	_     [0]int
	Mutex sync.RWMutex

	Pos token.Pos
	End token.Pos

	Fields      []*FieldInfo
	FieldsIndex map[string]int

	Methods      []*Method
	MethodsIndex map[string]int

	Dependencies      []*DependencyInfo
	DependenciesIndex map[string]int

	Incomplete bool
	isEmbedded bool
}

func (s *StructInfo) AddDependency(importIndex int, element string) {
	if index, exists := s.DependenciesIndex[element]; exists {
		s.Dependencies[index].ImportIndex = importIndex
		s.Dependencies[index].Usage++
	} else {
		dep := &DependencyInfo{
			ImportIndex: importIndex,
			Usage:       1,
		}
		s.Dependencies = append(s.Dependencies, dep)
		s.DependenciesIndex[element] = len(s.Dependencies) - 1
	}
}

func (s *StructInfo) AddMethod(metdod *Method, name string) {

	s.Methods = append(s.Methods, metdod)
	s.MethodsIndex[name] = len(s.Methods) - 1
}

func (s *StructInfo) SyncMethods(from *StructInfo) {

	for methodName, i := range from.MethodsIndex {
		if _, exists := s.MethodsIndex[methodName]; !exists {
			s.AddMethod(from.Methods[i], methodName)
		}
	}
}

func (s *StructInfo) SyncDependencies(from *StructInfo) {

	for element, i := range from.DependenciesIndex {
		if _, exists := s.DependenciesIndex[element]; !exists {
			s.AddDependency(from.Dependencies[i].ImportIndex, element)
		} else {
			s.Dependencies[s.DependenciesIndex[element]].Usage += from.Dependencies[i].Usage
		}
	}
}

func NewStructPreInit(name string) *StructInfo {
	methods := []*Method{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
		Methods:           methods,
		MethodsIndex:      methodsIndex,
		Dependencies:      make([]*DependencyInfo, 0),
		DependenciesIndex: make(map[string]int),
		isEmbedded:        NotEmbedded,
		Incomplete:        onlyPreinitialized,
	}
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructInfo, []UsedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	methods := []*Method{}
	methodsIndex := make(map[string]int)

	return &StructInfo{
			Pos:               res.Pos(),
			End:               res.End(),
			Fields:            mapMetaData.fieldsSet,
			FieldsIndex:       mapMetaData.fieldsIndex,
			Methods:           methods,
			MethodsIndex:      methodsIndex,
			Dependencies:      make([]*DependencyInfo, 0),
			DependenciesIndex: make(map[string]int),
			isEmbedded:        isEmbedded,
			Incomplete:        true,
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
	fieldsSet    []*FieldInfo
	fieldsIndex  map[string]int
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*FieldInfo, 0, len(fieldList))
	fieldsIndex := make(map[string]int, len(fieldList))
	usedPackages := []UsedPackage{}

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
	usedPackages   []UsedPackage
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
				usedPackages: []UsedPackage{{Alias: ident.Name, Element: ft.Sel.Name}},
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
