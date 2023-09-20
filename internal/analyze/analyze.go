package analyze

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"path/filepath"

	"github.com/g10z3r/archx/internal/analyze/entity"
	"github.com/g10z3r/archx/internal/analyze/snapshot"
)

func ParsePackage(dirPath string, mod string) (*PackageBuffer, error) {
	var buf *PackageBuffer

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	for pkgName, pkg := range pkgs {
		buf = NewPackageBuffer(pkgName)

		for fileName, file := range pkg.Files {
			log.Printf("Processing file: %s", fileName)

			for _, imp := range file.Imports {
				if imp.Path != nil && imp.Path.Value != "" {
					buf.AddImport(imp, mod)
				}
			}

			for _, decl := range file.Decls {
				switch d := decl.(type) {

				case *ast.FuncDecl:
					if err := processFuncDecl(buf, fset, d); err != nil {
						return nil, err
					}
				case *ast.GenDecl:
					if err := processGenDecl(buf, fset, d); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return buf, nil
}

func processFuncDecl(pkgBuf *PackageBuffer, fs *token.FileSet, funcDecl *ast.FuncDecl) error {
	if funcDecl.Recv != nil {
		if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
			if ident, ok := starExpr.X.(*ast.Ident); ok {
				fmt.Printf("Method belongs to struct: %s\n", ident.Name)
				ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
					// Add logic to process nodes within the function body
					return true
				})
			}
		}
	}
	return nil
}

func processGenDecl(buf *PackageBuffer, fs *token.FileSet, genDecl *ast.GenDecl) error {
	if genDecl.Tok != token.TYPE {
		return nil
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		sType, usedPackages, err := entity.NewStructType(fs, structType, entity.NotEmbedded)
		if err != nil {
			return fmt.Errorf("failed to create new struct type: %w", err)
		}

		for _, p := range usedPackages {
			if importPath, exist := buf.Imports[p.Alias]; exist {
				sType.AddDependency(importPath, p.Element)
			}
		}

		log.Printf("Found a struct: %s\n", typeSpec.Name.Name)
		buf.AddStruct(sType, typeSpec.Name.Name)
	}
	return nil
}

func ParseGoFile(filePath string, mod string) (*snapshot.FileManifest, error) {
	// Parse the package name
	pkgName, err := parsePackage(filePath)
	if err != nil {
		return nil, err
	}

	pkgDir, _ := filepath.Split(filePath)
	pkgPath := path.Join(pkgDir, pkgName)
	fileManifest := snapshot.NewFileManifest(pkgPath)

	// Parse the file to get the AST
	fset, node, err := parseFile(filePath)
	if err != nil {
		return nil, err
	}

	var currentStructName string
	methodsInfo := make(map[string]map[string]entity.Usage)

	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.ImportSpec:
			if t.Path != nil && t.Path.Value != "" {
				fileManifest.AddImport(t, mod)
			}

		case *ast.TypeSpec:
			if structType, ok := t.Type.(*ast.StructType); ok {
				currentStructName = t.Name.Name
				sType, _, err := entity.NewStructType(fset, structType, entity.NotEmbedded)
				if err != nil {
					return false
				}

				fileManifest.AddStruct(currentStructName, sType)
			}

			if interfaceType, ok := t.Type.(*ast.InterfaceType); ok {
				iType := entity.NewInterfaceType(interfaceType)
				fileManifest.AddInterface(t.Name.String(), iType)
			}

		case *ast.SelectorExpr:
			if xIdent, ok := t.X.(*ast.Ident); ok {
				if _, exists := fileManifest.Imports[xIdent.Name]; exists {
					// Get the index of the current struct from the StructsIndex map
					structIndex := fileManifest.StructsIndex[currentStructName]
					// Use the index to get the reference to the current StructInfo from the Structs slice
					currentStruct := fileManifest.Structs[structIndex]
					if xIdent.Pos() > currentStruct.Pos && xIdent.Pos() < currentStruct.End {
						// fmt.Println(xIdent.Name, t.Sel.Name, xIdent.Pos())
						currentStruct.AddDependency(fileManifest.Imports[xIdent.Name], t.Sel.Name)
					} else {
						for _, s := range fileManifest.Structs {
							fmt.Println(filePath, xIdent.Name, t.Sel.Name, xIdent.Pos(), xIdent.NamePos)
							fmt.Println(s.Pos, s.End)
							fmt.Println("+++++++++++++++++++")
							if xIdent.Pos() > s.Pos && xIdent.Pos() < s.End {
								s.AddDependency(fileManifest.Imports[xIdent.Name], t.Sel.Name)
							}
						}
					}
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
				methodFields = make(map[string]entity.Usage)
				methodsInfo[methodName] = methodFields
			}

			ast.Inspect(t, func(n ast.Node) bool {
				if ident, ok := n.(*ast.Ident); ok {
					if exists, _ := fileManifest.IsFieldPresent(currentStructName, ident.Name); exists {
						fieldUsage := methodFields[ident.Name]

						fieldUsage.Total++

						// If the field is encountered for the first time in this method, increment Uniq
						if _, seen := encounteredFields[ident.Name]; !seen {
							fieldUsage.Uniq++
							encounteredFields[ident.Name] = struct{}{}
						}

						methodFields[ident.Name] = fieldUsage
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

	for _, v := range fileManifest.Structs {
		v.DepsTree.Compress()
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

// func indexEncode(index, startPos, endPos uint64) uint64 {
// 	return (index << 48) | (startPos << 24) | endPos
// }

// func indexDecode(combinedIndex uint64) (uint64, uint64, uint64) {
// 	extractedIndex := (combinedIndex >> 48) & 0xFFFF
// 	extractedStartPos := (combinedIndex >> 24) & 0xFFFFFF
// 	extractedEndPos := combinedIndex & 0xFFFFFF

// 	return extractedIndex, extractedStartPos, extractedEndPos
// }
