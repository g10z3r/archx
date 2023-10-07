package entity

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
)

type astEntity interface {
	Pos() token.Pos
	End() token.Pos
}

func calcLineCount(fset *token.FileSet, res astEntity) int {
	return fset.Position(res.End()).Line - fset.Position(res.Pos()).Line + 1
}

type exprTypeMetaData struct {
	Type           string
	UsedPackages   []UsedPackage
	EmbeddedStruct *StructEntity
}

func ExtractExprAsType(fset *token.FileSet, expr ast.Expr) (*exprTypeMetaData, error) {
	switch ft := expr.(type) {
	case *ast.StructType:
		embedded, usedPackages, err := NewStructEntity(fset, ft, true, nil)
		if err != nil {
			return nil, err
		}

		return &exprTypeMetaData{
			Type:           "struct",
			UsedPackages:   usedPackages,
			EmbeddedStruct: embedded,
		}, nil

	case *ast.SelectorExpr:
		if ident, ok := ft.X.(*ast.Ident); ok {
			return &exprTypeMetaData{
				Type:         ft.Sel.Name,
				UsedPackages: []UsedPackage{{Alias: ident.Name, Element: ft.Sel.Name}},
			}, nil
		}

		return &exprTypeMetaData{
			Type: ft.Sel.Name,
		}, nil

	default:
		var buf bytes.Buffer
		if err := format.Node(&buf, fset, expr); err != nil {
			return nil, fmt.Errorf("failed to format node: %w", err)
		}

		return &exprTypeMetaData{
			Type: buf.String(),
		}, nil
	}
}
