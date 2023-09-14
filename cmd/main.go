package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type TypeData struct {
	Fields  map[string]struct{}
	Methods map[string]map[string]struct{}
}

func main() {
	data, err := analyzeGoFile("./cmd/main.go")
	if err != nil {
		fmt.Println("Error analyzing Go file:", err)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

}

func analyzeGoFile(filePath string) (map[string]*TypeData, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	data := make(map[string]*TypeData)
	currentTypeName := ""

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if _, ok := t.Type.(*ast.StructType); ok {
				currentTypeName = t.Name.Name
				data[currentTypeName] = &TypeData{
					Fields:  make(map[string]struct{}),
					Methods: make(map[string]map[string]struct{}),
				}
			}
		case *ast.Field:
			if currentTypeName != "" {
				for _, name := range t.Names {
					data[currentTypeName].Fields[name.Name] = struct{}{}
				}
			}
		case *ast.FuncDecl:
			if t.Recv != nil && len(t.Recv.List) > 0 {
				if se, ok := t.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := se.X.(*ast.Ident); ok {
						currentTypeName = ident.Name
						methodName := t.Name.Name
						data[currentTypeName].Methods[methodName] = make(map[string]struct{})
						ast.Inspect(t, func(n ast.Node) bool {
							if ident, ok := n.(*ast.Ident); ok {
								if _, ok := data[currentTypeName].Fields[ident.Name]; ok {
									data[currentTypeName].Methods[methodName][ident.Name] = struct{}{}
								}
							}
							return true
						})
					}
				}
			}
		}
		return true
	})

	return data, nil
}
