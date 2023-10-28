package analyzer

import (
	"go/ast"
	"go/token"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
	"github.com/g10z3r/archx/internal/domain/service/anthill/common"
)

type FunctionAnalyzer struct{}

func (a *FunctionAnalyzer) Name() string {
	return "func"
}

func (a *FunctionAnalyzer) Check(node ast.Node) bool {
	_, ok := node.(*ast.FuncDecl)
	return ok
}

func (a *FunctionAnalyzer) Save(f *obj.FileObj, object common.Object) {
	funcObj, ok := object.(*obj.FuncObj)
	if !ok {
		log.Fatal("not a func objects")
	}

	f.AppendFunc(funcObj)
}

func (a *FunctionAnalyzer) Analyze(f *obj.FileObj, node ast.Node) common.Object {
	funcDecl, _ := node.(*ast.FuncDecl)

	var parentStruct *ast.Ident
	if funcDecl.Recv != nil {
		starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}

		parentStruct, ok = starExpr.X.(*ast.Ident)
		if !ok {
			return nil
		}
	}

	params, deps, err := processFuncParams(f.FileSet, funcDecl, f.Entities.Imports.InternalImportsMeta)
	if err != nil {
		log.Fatal(err)
	}

	funcObj := obj.NewFuncObj(f.FileSet, funcDecl, params, deps, parentStruct)

	if err := inspectFuncBody(funcDecl, funcObj, f.Entities.Imports.InternalImportsMeta); err != nil {
		log.Fatal(err)
	}

	return funcObj
}

func inspectFuncBody(funcDecl *ast.FuncDecl, funcEntity *obj.FuncObj, impMeta map[string]int) error {
	// Get the recipient's name if it is a structure method
	var receiverName string
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		receiverName = funcDecl.Recv.List[0].Names[0].Name
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				// Structure field call found
				if ident.Name == receiverName {
					if usage, exists := funcEntity.Fields[expr.Sel.Name]; !exists {
						funcEntity.Fields[expr.Sel.Name] = usage
					}

					funcEntity.Fields[expr.Sel.Name]++
				}

				// Found using another internal package
				if index, exists := impMeta[ident.Name]; exists {
					funcEntity.AddDependency(index, expr.Sel.Name)
				}
			}

		case *ast.CallExpr:
			// Check for recursion in regular functions
			if ident, ok := expr.Fun.(*ast.Ident); ok {
				if ident.Name == funcEntity.Name {
					funcEntity.Metadata.IsRecursive = true
				}
			}

			// Check for recursion in methods
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

type FuncObjParamMap map[string]*obj.FuncObjParam
type DepObjMap map[string]*obj.EntityDepObj

func processFuncParams(fset *token.FileSet, funcDecl *ast.FuncDecl, impMeta map[string]int) (FuncObjParamMap, DepObjMap, error) {
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
