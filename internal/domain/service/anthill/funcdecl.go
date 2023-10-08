package anthill

import (
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

func (f *forager) processFuncDecl(fset *token.FileSet, funcDecl *ast.FuncDecl, impMeta map[string]int, fileName string) error {
	if funcDecl.Recv == nil {
		return nil
	}

	starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil
	}

	parentStruct, ok := starExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}

	params, deps, err := processFuncParams(fset, funcDecl, impMeta)
	if err != nil {
		return err
	}

	funcEntity := obj.NewFuncObj(fset, funcDecl, params, deps, &parentStruct.Name)
	receiver := funcDecl.Recv.List[0].Names[0]

	if err := inspectFuncBody(funcDecl, receiver, funcEntity, impMeta); err != nil {
		return err
	}

	bucket, _ := f.storage.funcBucket.Load(fileName)
	f.storage.funcBucket.Store(fileName, append(bucket, funcEntity))

	return nil
}

func inspectFuncBody(funcDecl *ast.FuncDecl, receiver *ast.Ident, funcEntity *obj.FuncObj, impMeta map[string]int) error {
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				// Structure field call found
				if ident.Name == receiver.Name {
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
		}
		return true
	})

	return nil
}

type FuncObjParamMap map[string]*obj.FuncObjParam
type DepObjMap map[string]*obj.DepObj

func processFuncParams(fset *token.FileSet, funcDecl *ast.FuncDecl, impMeta map[string]int) (FuncObjParamMap, DepObjMap, error) {
	var params map[string]*obj.FuncObjParam
	if len(funcDecl.Type.Params.List) > 0 {
		params = make(map[string]*obj.FuncObjParam)
	}

	deps := map[string]*obj.DepObj{}
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
				deps[typ.UsedPackages[0].Element] = &obj.DepObj{
					ImportIndex: index,
				}
			}

			deps[typ.UsedPackages[0].Element].Usage++
		}
	}

	return params, deps, nil
}
