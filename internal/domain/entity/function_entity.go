package entity

import (
	"go/ast"
	"go/token"
	"unicode"
)

type FuncMetadata struct {
	LineCount      int
	Arity          int
	ReturnCount    int
	IsRecursive    bool
	HasSideEffects bool
}

type FuncParam struct {
	Type  string
	Usage int
}

type FunctionEntity struct {
	Name         string
	Receiver     *string
	Fields       map[string]int
	Parameters   map[string]*FuncParam
	Dependencies map[string]*DependencyEntity
	Visibility   bool
	Metadata     *FuncMetadata
}

func (s *FunctionEntity) AddDependency(importIndex int, element string) {
	if _, exists := s.Dependencies[element]; !exists {
		s.Dependencies[element] = &DependencyEntity{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	s.Dependencies[element].Usage++
}

func NewFunctionEntity(fset *token.FileSet, res *ast.FuncDecl, params map[string]*FuncParam, initDeps map[string]*DependencyEntity, receiver *string) *FunctionEntity {
	var fields map[string]int

	if receiver != nil {
		fields = make(map[string]int)
	}

	return &FunctionEntity{
		Name:         res.Name.Name,
		Receiver:     receiver,
		Fields:       fields,
		Dependencies: initDeps,
		Parameters:   params,
		Visibility:   unicode.IsUpper(rune(res.Name.Name[0])),
		Metadata: &FuncMetadata{
			LineCount: calcLineCount(fset, res),
			Arity:     len(params),
		},
	}
}
