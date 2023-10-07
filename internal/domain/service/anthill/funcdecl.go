package anthill

import (
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/entity"
)

func (f *forager) processFuncDecl(funcDecl *ast.FuncDecl, impMeta map[string]int, fileName string) error {
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

	methodEntity := entity.NewMethodEntity(funcDecl, parentStruct.Name)
	receiver := funcDecl.Recv.List[0].Names[0]

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				if ident.Name == receiver.Name {
					if usage, exists := methodEntity.UsedFields[expr.Sel.Name]; !exists {
						methodEntity.UsedFields[expr.Sel.Name] = usage
					}

					methodEntity.UsedFields[expr.Sel.Name]++
				}

				if ident.Name != receiver.Name {
					if index, exists := impMeta[ident.Name]; exists {
						methodEntity.AddDependency(index, expr.Sel.Name)
					}
				}
			}
		}
		return true
	})

	bucket, _ := f.storage.methodBucket.Load(fileName)
	f.storage.methodBucket.Store(fileName, append(bucket, methodEntity))

	return nil
}
