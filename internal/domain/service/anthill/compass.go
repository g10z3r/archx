package anthill

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"sync"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
	"github.com/g10z3r/archx/internal/domain/service/anthill/common"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
)

type Manager struct {
	analyzers map[string]common.Analyzer
}

func (m *Manager) Register(a common.Analyzer) {
	m.analyzers[a.Name()] = a
}

type compassEvent interface {
	Name() string
}

type Compass struct {
	mutex sync.Mutex

	manager *Manager

	eventCh       chan compassEvent
	unsubscribeCh chan struct{}
}

func NewCompass() *Compass {
	return &Compass{
		manager: &Manager{
			analyzers: make(map[string]common.Analyzer),
		},

		eventCh:       make(chan compassEvent, 1),
		unsubscribeCh: make(chan struct{}),
	}
}

func (r *Compass) Subscribe() (<-chan compassEvent, chan struct{}) {
	return r.eventCh, r.unsubscribeCh
}

func (c *Compass) Parse(info *collector.Info, targetDir string) {
	importAlz := &analyzer.ImportAnalyzer{}
	c.manager.Register(importAlz)

	structAlz := &analyzer.StructAnalyzer{}
	c.manager.Register(structAlz)

	funcAlz := &analyzer.FunctionAnalyzer{}
	c.manager.Register(funcAlz)

	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, targetDir, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, pkgAst := range pkg {
		wg.Add(1)
		go func(pkgAst *ast.Package) {
			c.eventCh <- &event.PackageFormedEvent{
				Package: c.ParsePkg(fset, pkgAst, targetDir, info.ModuleName),
			}

			wg.Done()
		}(pkgAst)
	}

	wg.Wait()
}

func (c *Compass) ParsePkg(fset *token.FileSet, pkgAst *ast.Package, targetDir, moduleName string) *obj.PackageObj {
	pkgObj := obj.NewPackageObj(pkgAst, targetDir)

	var wg sync.WaitGroup
	for fileName, fileAst := range pkgAst.Files {
		wg.Add(1)
		go func(fileAst *ast.File, fileName string) {
			fileObj := c.ParseFile(fset, fileAst, moduleName, fileName)
			pkgObj.AppendFile(fileObj)

			wg.Done()
		}(fileAst, fileName)
	}

	wg.Wait()
	return pkgObj
}

func (c *Compass) ParseFile(fset *token.FileSet, fileAst *ast.File, moduleName, fileName string) *obj.FileObj {
	fileObj := obj.NewFileObj(fset, moduleName, filepath.Base(fileName))
	visitor := NewVisitor(fileObj, c.manager.analyzers)
	ast.Walk(visitor, fileAst)

	return fileObj
}
