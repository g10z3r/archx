package scaner

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sync"

	"github.com/g10z3r/archx/internal/scaner/buffer"
	"github.com/g10z3r/archx/internal/scaner/entity"
)

var errChan = make(chan error, 1)

func ScanPackage(dirPath string, mod string) (*buffer.BufferEventBus, error) {
	var buf *buffer.BufferEventBus

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		imports, total := processImports(pkg.Files)
		buf = buffer.NewBufferEventBus(mod, total, errChan)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf.Open()
		}()

		for i := 0; i < len(imports); i++ {
			buf.SendEvent(&buffer.AddImportEvent{Import: imports[i]})
		}

		for fileName, file := range pkg.Files {
			log.Printf("Processing file: %s", fileName)

			buf.WaitGroup.Add(1)
			go func(file *ast.File) {
				defer buf.WaitGroup.Done()

				for _, decl := range file.Decls {
					switch d := decl.(type) {
					case *ast.FuncDecl:
						processFuncDecl(buf, fset, d)
					case *ast.GenDecl:
						processGenDecl(buf, fset, d)
					}
				}
			}(file)
		}

		buf.WaitGroup.Wait()
		buf.Close()
		wg.Wait()
	}

	return buf, nil
}

func processImports(files map[string]*ast.File) ([]*entity.Import, int) {
	var impTotal int
	var imports []*entity.Import

	for _, file := range files {
		impTotal = impTotal + len(file.Imports)

		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				imports = append(imports, entity.NewImport(imp))
			}
		}
	}

	return imports, impTotal
}
