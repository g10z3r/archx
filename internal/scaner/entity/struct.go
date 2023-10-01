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
)

type Field struct {
	_        [0]int    `bson:"-"`
	pos      token.Pos `bson:"pos"`
	end      token.Pos `bson:"end"`
	Type     string    `bson:"type"`
	Embedded *Struct   `bson:"embedded"`
	IsPublic bool      `bson:"isPublic"`
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

type Dependency struct {
	ImportIndex int `bson:"importIndex"`
	Usage       int `bson:"usage"`
}

type Struct struct {
	_     [0]int       `bson:"-"`
	Mutex sync.RWMutex `bson:"-"`

	Pos token.Pos `bson:"pos"`
	End token.Pos `bson:"end"`

	Fields      []*Field       `bson:"fields"`
	FieldsIndex map[string]int `bson:"fieldsIndex"`

	Methods      []*Method      `bson:"methods"`
	MethodsIndex map[string]int `bson:"methodsIndex"`

	Dependencies      []*Dependency  `bson:"dependencies"`
	DependenciesIndex map[string]int `bson:"dependenciesIndex"`

	Incomplete bool `bson:"incomplete"`
	IsEmbedded bool `bson:"isEmbedded"`
}

func (s *Struct) AddDependency(importIndex int, element string) {
	if index, exists := s.DependenciesIndex[element]; exists {
		s.Dependencies[index].ImportIndex = importIndex
		s.Dependencies[index].Usage++
	} else {
		dep := &Dependency{
			ImportIndex: importIndex,
			Usage:       1,
		}
		s.Dependencies = append(s.Dependencies, dep)
		s.DependenciesIndex[element] = len(s.Dependencies) - 1
	}
}

func (s *Struct) AddMethod(metdod *Method, name string) {

	s.Methods = append(s.Methods, metdod)
	s.MethodsIndex[name] = len(s.Methods) - 1
}

func (s *Struct) SyncMethods(from *Struct) {

	for methodName, i := range from.MethodsIndex {
		if _, exists := s.MethodsIndex[methodName]; !exists {
			s.AddMethod(from.Methods[i], methodName)
		}
	}
}

func (s *Struct) SyncDependencies(from *Struct) {

	for element, i := range from.DependenciesIndex {
		if _, exists := s.DependenciesIndex[element]; !exists {
			s.AddDependency(from.Dependencies[i].ImportIndex, element)
		} else {
			s.Dependencies[s.DependenciesIndex[element]].Usage += from.Dependencies[i].Usage
		}
	}
}

func NewStructPreInit(name string) *Struct {
	return &Struct{
		Methods:           make([]*Method, 0),
		MethodsIndex:      make(map[string]int),
		Dependencies:      make([]*Dependency, 0),
		DependenciesIndex: make(map[string]int),
		IsEmbedded:        NotEmbedded,
		Incomplete:        onlyPreinitialized,
	}
}

func NewStructType(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*Struct, []UsedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	var methods []*Method
	var methodsIndex map[string]int
	var dependencies []*Dependency
	var dependenciesIndex map[string]int

	if !isEmbedded {
		methods = []*Method{}
		methodsIndex = make(map[string]int)
		dependencies = make([]*Dependency, 0, len(mapMetaData.usedPackages))
		dependenciesIndex = make(map[string]int, len(mapMetaData.usedPackages))
	}

	return &Struct{
			Pos:               res.Pos(),
			End:               res.End(),
			Fields:            mapMetaData.fieldsSet,
			FieldsIndex:       mapMetaData.fieldsIndex,
			Methods:           methods,
			MethodsIndex:      methodsIndex,
			Dependencies:      dependencies,
			DependenciesIndex: dependenciesIndex,
			IsEmbedded:        isEmbedded,
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
	fieldsSet    []*Field
	fieldsIndex  map[string]int
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*Field, 0, len(fieldList))
	fieldsIndex := make(map[string]int, len(fieldList))
	usedPackages := []UsedPackage{}

	for i, field := range fieldList {
		fieldMetaData, err := extractFieldType(fset, field.Type)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fieldMetaData.usedPackages); i++ {
			usedPackages = append(usedPackages, fieldMetaData.usedPackages[i])
		}

		for _, name := range field.Names {
			fieldsIndex[name.Name] = i
			fields = append(fields, &Field{
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
	embeddedStruct *Struct
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (*fieldTypeMetaData, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, usedPackages, err := NewStructType(fset, ft, true)
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
