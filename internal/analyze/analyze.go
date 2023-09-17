package analyze

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/g10z3r/archx/internal/analyze/snapshot"
	"github.com/g10z3r/archx/internal/analyze/types"
)

func ParseGoFile(filePath string) (*snapshot.FileManifest, error) {
	pkgName, err := parsePackage(filePath)
	if err != nil {
		return nil, err
	}

	pkgDir, _ := filepath.Split(filePath)
	pkgPath := fmt.Sprintf("%s%s", pkgDir, pkgName)
	fileManifest := snapshot.NewFileManifest(pkgPath)

	fset, node, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	var currentStructName string

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {
				currentStructName = t.Name.String()
				sType, err := types.NewStructType(fset, structType, types.NotEmbedded)
				if err != nil {
					return false
				}

				fileManifest.AddStructType(currentStructName, sType)
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

			// Collect information about the fields used within this method
			usedFields := make(map[string]bool)
			ast.Inspect(t, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if exists, _ := fileManifest.IsFieldPresent(currentStructName, ident.Name); exists {
						usedFields[ident.Name] = true
					}
				}

				return true
			})

			for fieldName := range usedFields {
				if err := fileManifest.AddMethodToStruct(currentStructName, methodName, fieldName); err != nil {
					return false
				}
			}

		}
		return true
	})

	return fileManifest, nil
}

// func MustParseGoFile(filePath string) *SnapshotManifest {
// 	pkg, err := parsePackage(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if Snapshot == nil {
// 		Snapshot = NewSnapshot()
// 	}

// 	pkgPath := Snapshot.UpsertPackageManifest(filePath, pkg)

// 	fset, node, err := parseFile(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	var currentStructName string

// 	ast.Inspect(node, func(n ast.Node) bool {
// 		switch t := n.(type) {
// 		case *ast.TypeSpec:
// 			if structType, ok := t.Type.(*ast.StructType); ok {
// 				currentStructName = t.Name.String()
// 				sType, err := types.NewStructType(fset, structType, types.NotEmbedded)
// 				if err != nil {
// 					panic(err)
// 				}

// 				Snapshot.AddStructType(pkgPath, currentStructName, sType)
// 			}

// 		case *ast.FuncDecl:
// 			if t.Recv == nil || len(t.Recv.List) == 0 {
// 				return true
// 			}

// 			se, ok := t.Recv.List[0].Type.(*ast.StarExpr)
// 			if !ok {
// 				return true
// 			}

// 			ident, ok := se.X.(*ast.Ident)
// 			if !ok {
// 				return true
// 			}

// 			currentStructName = ident.Name
// 			methodName := t.Name.Name

// 			if !Snapshot.HasStructType(pkgPath, currentStructName) {
// 				return true
// 			}

// 			// Collect information about the fields used within this method
// 			usedFields := make(map[string]bool)
// 			ast.Inspect(t, func(n ast.Node) bool {
// 				if ident, ok := n.(*ast.Ident); ok {
// 					if exists, _ := Snapshot.IsFieldPresent(pkgPath, currentStructName, ident.Name); exists {
// 						usedFields[ident.Name] = true
// 					}
// 				}

// 				return true
// 			})

// 			for fieldName := range usedFields {
// 				if err := Snapshot.AddMethodToStruct(pkgPath, currentStructName, methodName, fieldName); err != nil {
// 					panic(err)
// 				}

// 			}

// 		}
// 		return true
// 	})

// 	return Snapshot
// }

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
