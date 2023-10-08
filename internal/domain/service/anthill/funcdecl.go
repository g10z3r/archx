package anthill

import (
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/domain/entity"
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

	var params map[string]*entity.FuncParam
	if len(funcDecl.Type.Params.List) > 0 {
		params = make(map[string]*entity.FuncParam)
	}

	paramsDeps := map[string]*entity.DependencyEntity{}
	for _, param := range funcDecl.Type.Params.List {
		for _, name := range param.Names {
			typ, err := entity.ExtractExprAsType(fset, param.Type)
			if err != nil {
				return err
			}

			params[name.Name] = &entity.FuncParam{
				Type: typ.Type,
			}

			if len(typ.UsedPackages) < 1 {
				continue
			}

			if index, exists := impMeta[typ.UsedPackages[0].Alias]; exists {
				paramsDeps[typ.UsedPackages[0].Element] = &entity.DependencyEntity{
					ImportIndex: index,
				}
			}

			paramsDeps[typ.UsedPackages[0].Element].Usage++
		}
	}

	funcEntity := entity.NewFunctionEntity(fset, funcDecl, params, paramsDeps, &parentStruct.Name)
	receiver := funcDecl.Recv.List[0].Names[0]

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				if ident.Name == receiver.Name {
					if usage, exists := funcEntity.Fields[expr.Sel.Name]; !exists {
						funcEntity.Fields[expr.Sel.Name] = usage
					}
					funcEntity.Fields[expr.Sel.Name]++
				}

				if ident.Name != receiver.Name {
					if index, exists := impMeta[ident.Name]; exists {
						funcEntity.AddDependency(index, expr.Sel.Name)
					}
				}
			}
		}
		return true
	})

	bucket, _ := f.storage.funcBucket.Load(fileName)
	f.storage.funcBucket.Store(fileName, append(bucket, funcEntity))

	return nil
}
