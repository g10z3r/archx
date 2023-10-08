package obj

import (
	"go/ast"
	"go/token"
	"unicode"
)

type FuncObjMetadata struct {
	LineCount      int
	Arity          int
	ReturnCount    int
	IsRecursive    bool
	HasSideEffects bool
}

type FuncObjParam struct {
	Type  string
	Usage int
}

type FuncObj struct {
	Name         string
	Receiver     *string
	Fields       map[string]int
	Parameters   map[string]*FuncObjParam
	Dependencies map[string]*DepObj
	Visibility   bool
	Metadata     *FuncObjMetadata
}

func (f *FuncObj) AddDependency(importIndex int, element string) {
	if _, exists := f.Dependencies[element]; !exists {
		f.Dependencies[element] = &DepObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	f.Dependencies[element].Usage++
}

func NewFuncObj(fset *token.FileSet, res *ast.FuncDecl, params map[string]*FuncObjParam, initDeps map[string]*DepObj, receiver *ast.Ident) *FuncObj {
	var receiverName *string
	var fields map[string]int

	if receiver != nil {
		receiverName = &receiver.Name
		fields = make(map[string]int)
	}

	return &FuncObj{
		Name:         res.Name.Name,
		Receiver:     receiverName,
		Fields:       fields,
		Dependencies: initDeps,
		Parameters:   params,
		Visibility:   unicode.IsUpper(rune(res.Name.Name[0])),
		Metadata: &FuncObjMetadata{
			LineCount: calcLineCount(fset, res),
			Arity:     len(params),
		},
	}
}
