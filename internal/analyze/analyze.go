package analyze

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/g10z3r/archx/internal/analyze/types"
)

func ParseGoFile(filePath string) (types.NodeType, error) {
	_, node, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	nodeData := make(types.NodeType)
	var currentStructName string

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {
				currentStructName = t.Name.String()
				nodeData[currentStructName] = types.MakeStructType(structType)
			}

		case *ast.FuncDecl:
			if t.Recv == nil || len(t.Recv.List) == 0 {
				return true
			}

			se, ok := t.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				return true
			}

			ident, ok := se.X.(*ast.Ident)
			if !ok {
				return true
			}

			currentStructName = ident.Name
			methodName := t.Name.Name

			if nodeData[currentStructName] == nil {
				return true
			}

			// Collect information about the fields used within this method
			usedFields := make(map[string]bool)
			ast.Inspect(t, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if _, ok := nodeData[currentStructName].Field[ident.Name]; ok {
						usedFields[ident.Name] = true
					}
				}
				return true
			})

			for fieldName := range usedFields {
				nodeData[currentStructName].Method[methodName] = append(
					nodeData[currentStructName].Method[methodName], fieldName,
				)
			}

		}
		return true
	})

	return nodeData, nil
}

func parseFile(filePath string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	return fset, node, nil
}
