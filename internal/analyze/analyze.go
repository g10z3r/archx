package analyze

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"

	"github.com/g10z3r/archx/internal/analyze/snapshot"
	"github.com/g10z3r/archx/internal/analyze/types"
)

func ParseGoFile(filePath string, mod string) (*snapshot.FileManifest, error) {
	pkgName, err := parsePackage(filePath)
	if err != nil {
		return nil, err
	}

	pkgDir, _ := filepath.Split(filePath)
	pkgPath := path.Join(pkgDir, pkgName)
	fileManifest := snapshot.NewFileManifest(pkgPath)

	fset, node, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	var currentStructName string
	methodsInfo := make(map[string]map[string]types.FieldUsage)

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.ImportSpec:
			if t.Path != nil && t.Path.Value != "" {
				fileManifest.AddImport(t, mod)
			}

		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {
				currentStructName = t.Name.String()
				sType, err := types.NewStructType(fset, structType, types.NotEmbedded)
				if err != nil {
					return false
				}

				fileManifest.AddStructType(currentStructName, sType)
			}

			if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
				iType := types.NewInterfaceType(interfaceType)
				fileManifest.AddInterfaceType(t.Name.String(), iType)
			}

		case *ast.SelectorExpr:
			if xIdent, ok := t.X.(*ast.Ident); ok {
				if _, exists := fileManifest.Imports[xIdent.Name]; exists {
					fileManifest.StructTypeMap[currentStructName].AddDependency(
						fileManifest.Imports[xIdent.Name],
						t.Sel.Name,
					)
				}
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

			if !fileManifest.HasStructType(currentStructName) {
				return true
			}

			// Tracks unique field accesses within a method.
			encounteredFields := make(map[string]struct{})

			// Collect information about the fields used within this method
			methodFields, exists := methodsInfo[methodName]
			if !exists {
				methodFields = make(map[string]types.FieldUsage)
				methodsInfo[methodName] = methodFields
			}

			ast.Inspect(t, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if exists, _ := fileManifest.IsFieldPresent(currentStructName, ident.Name); exists {
						// Increment the Total counter each time a field is encountered
						methodFields[ident.Name] = types.FieldUsage{
							Total: methodFields[ident.Name].Total + 1,
							Uniq:  methodFields[ident.Name].Uniq,
						}

						// If the field is seen for the first time, increment the Uniq counter
						if _, seen := encounteredFields[ident.Name]; !seen {
							methodFields[ident.Name] = types.FieldUsage{
								Total: methodFields[ident.Name].Total,
								Uniq:  methodFields[ident.Name].Uniq + 1,
							}
							encounteredFields[ident.Name] = struct{}{}
						}
					}
				}
				return true
			})

		}
		return true
	})

	// Add methods information to FileManifest outside the AST inspect to avoid repeated work
	for methodName, fields := range methodsInfo {
		for fieldName, fieldUsage := range fields {
			if err := fileManifest.AddMethodToStruct(currentStructName, methodName, fieldName, fieldUsage); err != nil {
				return nil, err
			}
		}
	}

	return fileManifest, nil
}

func parseFile(filePath string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse AST in %s: %v", filePath, err)
	}

	return fset, node, nil
}

func parsePackage(filePath string) (string, error) {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, filePath, nil, parser.PackageClauseOnly)
	if err != nil {
		return "", fmt.Errorf("failed to parse package in %s: %v", filePath, err)
	}

	return file.Name.Name, nil
}

func handleCallExpr(expr *ast.CallExpr) {
	fun := expr.Fun
	switch fn := fun.(type) {
	case *ast.Ident:
		fmt.Printf("Function called: %s\n", fn.Name)
	case *ast.SelectorExpr:
		xIdent, ok := fn.X.(*ast.Ident)
		if ok {
			fmt.Printf("Function called: %s.%s\n", xIdent.Name, fn.Sel.Name)
		}
	}
}
