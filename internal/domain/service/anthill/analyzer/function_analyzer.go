package analyzer

import (
	"context"
	"errors"
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

func NewFuncAnalyzer(file *obj.FileObj) Analyzer[ast.Node, obj.Object] {
	return NewAnalyzer[ast.Node, obj.Object](
		file,
		analyzeFuncNode,
	)
}

// func checkFuncNode(node ast.Node) bool {
// 	_, ok := node.(*ast.FuncDecl)
// 	return ok
// }

func analyzeFuncNode(ctx context.Context, f *obj.FileObj, node ast.Node) (obj.Object, error) {
	funcDecl, _ := node.(*ast.FuncDecl)

	ps, err := getParentStruct(funcDecl)
	if err != nil {
		return nil, err
	}

	params, deps, err := processFuncParams(f.FileSet, funcDecl, f.Entities.Imports.InternalImportsMeta)
	if err != nil {
		return nil, errors.New("some error from analyzeFuncNode 3") // TODO: add normal error return message
	}

	funcObj := obj.NewFuncObj(f.FileSet, funcDecl, params, deps, ps)

	if err := inspectFuncBody(funcDecl, funcObj, f.Entities.Imports.InternalImportsMeta); err != nil {
		return nil, errors.New("some error from analyzeFuncNode 4") // TODO: add normal error return message
	}

	return funcObj, nil
}

func getParentStruct(funcDecl *ast.FuncDecl) (*ast.Ident, error) {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return nil, nil
	}

	receiverType := funcDecl.Recv.List[0].Type

	switch t := receiverType.(type) {
	case *ast.StarExpr:
		// if the receiver's type is a pointer, attempt to get the identifier of the struct
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident, nil
		}

	case *ast.Ident:
		// if the receiver's type is not a pointer, it's a regular struct, so return the identifier of the struct
		return t, nil
	}

	return nil, errors.New("invalid receiver type in method declaration")
}

func inspectFuncBody(funcDecl *ast.FuncDecl, funcEntity *obj.FuncObj, impMeta map[string]int) error {
	// get the recipient's name if it is a structure method
	var receiverName string
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		receiverName = funcDecl.Recv.List[0].Names[0].Name
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				// structure field call found
				if ident.Name == receiverName {
					if usage, exists := funcEntity.Fields[expr.Sel.Name]; !exists {
						funcEntity.Fields[expr.Sel.Name] = usage
					}

					funcEntity.Fields[expr.Sel.Name]++
				}

				// found using another internal package
				if index, exists := impMeta[ident.Name]; exists {
					funcEntity.AddDependency(index, expr.Sel.Name)
				}
			}

		case *ast.CallExpr:
			// check for recursion in regular functions
			if ident, ok := expr.Fun.(*ast.Ident); ok {
				if ident.Name == funcEntity.Name {
					funcEntity.Metadata.IsRecursive = true
				}
			}

			// check for recursion in methods
			if sel, ok := expr.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == receiverName && sel.Sel.Name == funcEntity.Name {
					funcEntity.Metadata.IsRecursive = true
				}
			}
		}
		return true
	})

	return nil
}

type funcObjParamMap map[string]*obj.FuncObjParam
type depObjMap map[string]*obj.EntityDepObj

func processFuncParams(fset *token.FileSet, funcDecl *ast.FuncDecl, impMeta map[string]int) (funcObjParamMap, depObjMap, error) {
	var params map[string]*obj.FuncObjParam
	if len(funcDecl.Type.Params.List) > 0 {
		params = make(map[string]*obj.FuncObjParam)
	}

	deps := map[string]*obj.EntityDepObj{}
	for _, param := range funcDecl.Type.Params.List {
		for _, name := range param.Names {
			typ, err := obj.ExtractExprAsType(fset, param.Type)
			if err != nil {
				return nil, nil, err
			}

			params[name.Name] = &obj.FuncObjParam{
				Type: typ.Type,
			}

			if len(typ.UsedPackages) < 1 {
				continue
			}

			if index, exists := impMeta[typ.UsedPackages[0].Alias]; exists {
				deps[typ.UsedPackages[0].Element] = &obj.EntityDepObj{
					ImportIndex: index,
				}
			}

			deps[typ.UsedPackages[0].Element].Usage++
		}
	}

	return params, deps, nil
}
