package scanner

import (
	"context"
	"go/ast"
	"log"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type packageBuffer interface {
	AddMethod(structName string, method *entity.MethodEntity)
	GetAndClearMethods(structName string) []*entity.MethodEntity
}

func (pa *packageActor) processFuncDecl(ctx context.Context, funcDecl *ast.FuncDecl) error {
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

	newMethod := entity.NewMethodEntity(funcDecl)
	receiver := funcDecl.Recv.List[0].Names[0]

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				if ident.Name == receiver.Name {
					if usage, exists := newMethod.UsedFields[expr.Sel.Name]; !exists {
						newMethod.UsedFields[expr.Sel.Name] = usage
					}
				}

				if ident.Name != receiver.Name {
					if index := pa.cache.GetImportIndex(ident.Name); index >= 0 {
						newMethod.AddDependency(index, expr.Sel.Name)
					}
				}
			}
		}
		return true
	})

	structIndex := pa.cache.GetStructIndex(parentStruct.Name)
	if structIndex < 0 {
		log.Printf("Sending method %s to buffer", newMethod.Name)
		pa.buf.AddMethod(parentStruct.Name, newMethod)
		return nil
	}

	// s.db.PackageAcc().StructAcc().AddMethod(ctx, newMethod, structIndex, pkgPath)

	return nil
}
