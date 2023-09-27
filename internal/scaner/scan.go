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
		impTotal := countTotalImports(pkg.Files)
		buf = buffer.NewBufferEventBus(mod, impTotal, errChan)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf.Open()
		}()

		processImports(buf, pkg.Files)

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

func processImports(buf *buffer.BufferEventBus, files map[string]*ast.File) {
	for _, file := range files {
		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				buf.SendEvent(&buffer.AddImportEvent{
					Import: entity.NewImport(imp),
				})
			}
		}
	}
}

func countTotalImports(files map[string]*ast.File) int {
	var impTotal int
	for _, file := range files {
		impTotal = impTotal + len(file.Imports)
	}

	return impTotal
}
