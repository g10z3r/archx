package scaner

import (
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/scaner/buffer"
	"github.com/g10z3r/archx/internal/scaner/entity"
)

func processFuncDecl(buf *buffer.BufferEventBus, fs *token.FileSet, funcDecl *ast.FuncDecl) {
	if funcDecl.Recv == nil {
		return
	}

	starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return
	}

	parentStruct, ok := starExpr.X.(*ast.Ident)
	if !ok {
		return
	}

	newMethod, structIndex := processMethod(buf, funcDecl, parentStruct.Name)
	notifyMethodAddition(buf, structIndex, newMethod, funcDecl.Name.Name)
}

func processMethod(buf *buffer.BufferEventBus, funcDecl *ast.FuncDecl, parentStructName string) (*entity.Method, int) {
	newMethod := entity.NewMethod(funcDecl)
	sType, structIndex := getOrCreateStruct(buf, parentStructName)

	receiverName := funcDecl.Recv.List[0].Names[0].Name
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				if ident.Name == receiverName {
					if usage, exists := newMethod.UsedFields[expr.Sel.Name]; !exists {
						newMethod.UsedFields[expr.Sel.Name] = usage
					}

					newMethod.UsedFields[expr.Sel.Name]++
				}

				if ident.Name != receiverName {
					if importIndex, exists := buf.ImportBuffer.GetIndexByAlias(ident.Name); exists {
						sType.AddDependency(importIndex, expr.Sel.Name)
					}
				}
			}
		}
		return true
	})

	return newMethod, structIndex
}

func getOrCreateStruct(buf *buffer.BufferEventBus, structName string) (*entity.Struct, int) {
	var structIndex int
	var sType *entity.Struct

	if !buf.StructBuffer.IsPresent(structName) {
		sType = entity.NewStructPreInit(structName)
		notifyStructCreation(buf, sType, structName)

		structIndex = buf.StructBuffer.GetIndex(structName)
	} else {
		structIndex = buf.StructBuffer.GetIndex(structName)
		sType = buf.StructBuffer.GetByIndex(structIndex)
	}

	return sType, structIndex
}

func notifyStructCreation(buf *buffer.BufferEventBus, sType *entity.Struct, structName string) {
	buf.SendEvent(
		&buffer.UpsertStructEvent{
			StructInfo: sType,
			StructName: structName,
		},
	)
}

func notifyMethodAddition(buf *buffer.BufferEventBus, structIndex int, newMethod *entity.Method, methodName string) {
	buf.SendEvent(
		&buffer.AddMethodEvent{
			StructIndex: structIndex,
			Method:      newMethod,
			MethodName:  methodName,
		},
	)
}
