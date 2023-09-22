package analyze

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"path/filepath"

	"github.com/g10z3r/archx/internal/analyze/buffer"
	"github.com/g10z3r/archx/internal/analyze/entity"
	"github.com/g10z3r/archx/internal/analyze/snapshot"
)

func ParsePackage(dirPath string, mod string) (*buffer.ManagerBuffer, error) {
	var buf *buffer.ManagerBuffer
	errChan := make(chan error, 1)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		buf = buffer.NewManagerBuffer(errChan)
		go buf.Start()

		for fileName, file := range pkg.Files {
			log.Printf("Processing file: %s", fileName)

			for _, imp := range file.Imports {
				if imp.Path != nil && imp.Path.Value != "" {
					buf.SendEvent(&buffer.AddImportEvent{
						Import: entity.NewImport(imp),
						Mod:    mod,
					})
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

		buf.WaitGroup.Wait()
		buf.Stop()
	}

	return buf, nil
}

func processFuncDecl(buf *buffer.ManagerBuffer, fs *token.FileSet, funcDecl *ast.FuncDecl) error {
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

	newMethod := entity.NewMethod(funcDecl)

	var sType *entity.StructInfo
	var structIndex int
	var isNew bool

	if !buf.StructBuffer.IsPresent(parentStruct.Name) {
		isNew = true
		sType = entity.NewStructPreInit(parentStruct.Name)
	} else {
		structIndex = buf.StructBuffer.GetIndex(parentStruct.Name)
		sType = buf.StructBuffer.GetByIndex(structIndex)
	}

	log.Printf("Method %s belongs to struct: %s\n", funcDecl.Name.Name, parentStruct.Name)

	receiverName := funcDecl.Recv.List[0].Names[0].Name
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				stName, _ := expr.X.(*ast.Ident)
				fmt.Println(stName.Name, expr.Sel.Name, funcDecl.Name.Name, ident.Name, parentStruct.Name)
				if ident.Name == receiverName {
					usage, exists := newMethod.UsedFields[expr.Sel.Name]
					if exists {
						usage.Total++
					} else {
						usage = entity.Usage{
							Total: 1,
							Uniq:  1,
						}
					}
					newMethod.UsedFields[expr.Sel.Name] = usage
					fmt.Printf("Accessing field %s of struct %s\n", expr.Sel.Name, ident.Name)
				}
			}
		}

		return true
	})

	if isNew {
		buf.SendEvent(
			&buffer.UpsertStructEvent{
				StructInfo: sType,
				StructName: parentStruct.Name,
			},
		)
	}

	buf.WaitGroup.Add(1)
	go func() {
		defer buf.WaitGroup.Done()

		buf.SendEvent(
			&buffer.AddMethodEvent{
				StructIndex: structIndex,
				Method:      newMethod,
				MethodName:  funcDecl.Name.Name,
			},
		)
	}()

	return nil
}

func processGenDecl(buf *buffer.ManagerBuffer, fs *token.FileSet, genDecl *ast.GenDecl) error {
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
			if importPath, exists := buf.ImportBuffer.IsPresent(p.Alias); exists {
				sType.AddDependency(importPath, p.Element)
			}
		}

		log.Printf("Found a struct: %s\n", typeSpec.Name.Name)

		buf.WaitGroup.Add(1)
		go func() {
			defer buf.WaitGroup.Done()

			buf.SendEvent(
				&buffer.UpsertStructEvent{
					StructInfo: sType,
					StructName: typeSpec.Name.Name,
				},
			)
		}()

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
