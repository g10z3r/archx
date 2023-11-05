package obj

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"unicode"
)

type (
	FuncObjMeta struct {
		LineCount      int
		Arity          int
		ReturnCount    int
		IsRecursive    bool
		HasSideEffects bool
	}

	FuncObjParam struct {
		Type  string
		Usage int
	}

	FuncDeclObj struct {
		Start        token.Pos
		End          token.Pos
		Name         string
		Receiver     *string
		FieldAccess  map[string]int
		Parameters   map[string]*FuncObjParam
		Dependencies map[string]*DependencyObj
		Visibility   bool
		Metadata     *FuncObjMeta
	}
)

func (o *FuncDeclObj) Type() string {
	return "func"
}

func (o *FuncDeclObj) AddDependency(importIndex int, element string) {
	if _, exists := o.Dependencies[element]; !exists {
		o.Dependencies[element] = &DependencyObj{
			ImportIndex: importIndex,
			Usage:       1,
		}

		return
	}

	o.Dependencies[element].Usage++
}

func NewFuncDeclObj(fset *token.FileSet, res *ast.FuncDecl, params map[string]*FuncObjParam, initDeps map[string]*DependencyObj, receiver *ast.Ident) *FuncDeclObj {
	var receiverName *string
	var fields map[string]int

	if receiver != nil {
		receiverName = &receiver.Name
		fields = make(map[string]int)
	}

	return &FuncDeclObj{
		Start:        res.Pos(),
		End:          res.End(),
		Name:         res.Name.Name,
		Receiver:     receiverName,
		FieldAccess:  fields,
		Dependencies: initDeps,
		Parameters:   params,
		Visibility:   unicode.IsUpper(rune(res.Name.Name[0])),
		Metadata: &FuncObjMeta{
			LineCount: CalcEntityLOC(fset, res),
			Arity:     len(params),
		},
	}
}

type (
	// TODO: get rid of this structure
	FuncFieldObj struct {
		Name string
		Type string
	}

	FuncTypeObj struct {
		Params  []*StructFieldObj         // TODO: convert it into a map
		Results map[string]*FuncFieldObj  // TODO: Not implemented
		Deps    map[string]*DependencyObj // TODO: Not implemented
	}
)

func NewFuncTypeObj(fset *token.FileSet, node ast.Node) (*FuncTypeObj, error) {
	ts, ok := node.(*ast.TypeSpec)
	if !ok {
		return nil, fmt.Errorf("node is not a TypeSpec: %s", reflect.TypeOf(node))
	}

	funcType, ok := ts.Type.(*ast.FuncType)
	if !ok {
		return nil, fmt.Errorf("node is not a FuncType: %s", reflect.TypeOf(node))
	}

	extrParamsData, err := extractFieldMap(fset, funcType.Params.List)
	if err != nil {
		return nil, fmt.Errorf("failed to extract func field map: %w", err)
	}

	return &FuncTypeObj{
		Params: extrParamsData.fieldsSet,
	}, nil
}
