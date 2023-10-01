package dto

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

	onlyPreinitialized = false
)

type FieldDTO struct {
	_        [0]int
	pos      token.Pos
	end      token.Pos
	Type     string
	Embedded *StructDTO
	IsPublic bool
}

type MethodDTO struct {
	Start      token.Pos
	End        token.Pos
	UsedFields map[string]int
	IsPublic   bool
}

func NewMethodDTO(res *ast.FuncDecl) *MethodDTO {
	return &MethodDTO{
		Start:      res.Pos(),
		End:        res.End(),
		UsedFields: make(map[string]int),
		IsPublic:   unicode.IsUpper(rune(res.Name.Name[0])),
	}
}

type DependencyDTO struct {
	ImportIndex int
	Usage       int
}

type StructDTO struct {
	_ [0]int

	Pos token.Pos
	End token.Pos

	Fields      []*FieldDTO
	FieldsIndex map[string]int

	Methods      []*MethodDTO
	MethodsIndex map[string]int

	Dependencies      []*DependencyDTO
	DependenciesIndex map[string]int

	Incomplete bool
	IsEmbedded bool
}

func (s *StructDTO) AddDependency(importIndex int, element string) {
	if index, exists := s.DependenciesIndex[element]; exists {
		s.Dependencies[index].ImportIndex = importIndex
		s.Dependencies[index].Usage++
	} else {
		dep := &DependencyDTO{
			ImportIndex: importIndex,
			Usage:       1,
		}
		s.Dependencies = append(s.Dependencies, dep)
		s.DependenciesIndex[element] = len(s.Dependencies) - 1
	}
}

func (s *StructDTO) AddMethod(metdod *MethodDTO, name string) {
	s.Methods = append(s.Methods, metdod)
	s.MethodsIndex[name] = len(s.Methods) - 1
}

func NewStructDTO(fset *token.FileSet, res *ast.StructType, isEmbedded bool) (*StructDTO, []UsedPackage, error) {
	mapMetaData, err := extractFieldMap(fset, res.Fields.List)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract field map: %w", err)
	}

	var methods []*MethodDTO
	var methodsIndex map[string]int
	var dependencies []*DependencyDTO
	var dependenciesIndex map[string]int

	if !isEmbedded {
		methods = []*MethodDTO{}
		methodsIndex = make(map[string]int)
		dependencies = make([]*DependencyDTO, 0, len(mapMetaData.usedPackages))
		dependenciesIndex = make(map[string]int, len(mapMetaData.usedPackages))
	}

	return &StructDTO{
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
	fieldsSet    []*FieldDTO
	fieldsIndex  map[string]int
}

func extractFieldMap(fset *token.FileSet, fieldList []*ast.Field) (*fieldMapMetaData, error) {
	fields := make([]*FieldDTO, 0, len(fieldList))
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
			fields = append(fields, &FieldDTO{
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
	embeddedStruct *StructDTO
}

func extractFieldType(fset *token.FileSet, fieldType ast.Expr) (*fieldTypeMetaData, error) {
	switch ft := fieldType.(type) {
	case *ast.StructType:
		embedded, usedPackages, err := NewStructDTO(fset, ft, true)
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
