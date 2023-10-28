package obj

import (
	"go/ast"
	"go/token"
	"unicode"
)

type FuncObjMeta struct {
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
	Start        token.Pos
	End          token.Pos
	Name         string
	Receiver     *string
	Fields       map[string]int
	Parameters   map[string]*FuncObjParam
	Dependencies map[string]*EntityDepObj
	Visibility   bool
	Metadata     *FuncObjMeta
}

func (o *FuncObj) Type() string {
	return "func"
}

func (o *FuncObj) AddDependency(importIndex int, element string) {
	if _, exists := o.Dependencies[element]; !exists {
		o.Dependencies[element] = &EntityDepObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	o.Dependencies[element].Usage++
}

func NewFuncObj(fset *token.FileSet, res *ast.FuncDecl, params map[string]*FuncObjParam, initDeps map[string]*EntityDepObj, receiver *ast.Ident) *FuncObj {
	var receiverName *string
	var fields map[string]int

	if receiver != nil {
		receiverName = &receiver.Name
		fields = make(map[string]int)
	}

	return &FuncObj{
		Start:        res.Pos(),
		End:          res.End(),
		Name:         res.Name.Name,
		Receiver:     receiverName,
		Fields:       fields,
		Dependencies: initDeps,
		Parameters:   params,
		Visibility:   unicode.IsUpper(rune(res.Name.Name[0])),
		Metadata: &FuncObjMeta{
			LineCount: CalcEntityLOC(fset, res),
			Arity:     len(params),
		},
	}
}
