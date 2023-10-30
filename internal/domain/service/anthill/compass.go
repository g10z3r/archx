package anthill

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"sync"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
	"github.com/g10z3r/archx/internal/domain/service/anthill/config"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
	"github.com/g10z3r/archx/internal/domain/service/anthill/pipe"
	"github.com/g10z3r/archx/internal/domain/service/anthill/pipe/plugin"
)

type Compass struct {
	mutex    sync.Mutex
	pipeline *pipe.Pipeline
	config   *config.Config

	eventCh       chan event.Event
	unsubscribeCh chan struct{}
}

func NewCompass() *Compass {
	return &Compass{
		config: &config.Config{
			Analysis: make(analyzer.AnalyzerMap),
		},

		eventCh:       make(chan event.Event, 1),
		unsubscribeCh: make(chan struct{}),
	}
}

func (c *Compass) RegisterAnalyzer(alz analyzer.Analyzer) error {
	if _, ok := c.config.Analysis[alz.Name()]; ok {
		return fmt.Errorf("analyzer %s already exists", alz.Name())
	}

	c.config.Analysis[alz.Name()] = alz
	return nil
}

func (c *Compass) DeleteAnalyzer(name string) {
	delete(c.config.Analysis, name)
}

func (r *Compass) Subscribe() (<-chan event.Event, chan struct{}) {
	return r.eventCh, r.unsubscribeCh
}

func (c *Compass) Run(ctx context.Context) {
	res := c.pipeline.Run(ctx, &plugin.CollectorPluginInput{
		RootDir:     ".",
		IgnoredList: config.DefaultIgnoredMap,
	})

	for _, dir := range res.([]string) {
		fmt.Println(dir)
	}
}

func (c *Compass) ParseDir(info *collector.Info, targetDir string) error {
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
				Package: c.parsePkg(fset, pkgAst, targetDir, info.ModuleName),
			}

			wg.Done()
		}(pkgAst)
	}

	wg.Wait()
	return nil
}

func (c *Compass) parsePkg(fset *token.FileSet, pkgAst *ast.Package, targetDir, moduleName string) *obj.PackageObj {
	pkgObj := obj.NewPackageObj(pkgAst, targetDir)

	var wg sync.WaitGroup
	for fileName, fileAst := range pkgAst.Files {
		wg.Add(1)
		go func(fileAst *ast.File, fileName string) {
			fileObj := c.parseFile(fset, fileAst, moduleName, fileName)
			pkgObj.AppendFile(fileObj)

			wg.Done()
		}(fileAst, fileName)
	}

	wg.Wait()
	return pkgObj
}

func (c *Compass) parseFile(fset *token.FileSet, fileAst *ast.File, moduleName, fileName string) *obj.FileObj {
	fileObj := obj.NewFileObj(fset, moduleName, filepath.Base(fileName))
	visitor := NewVisitor(fileObj, c.config.Analysis)
	ast.Walk(visitor, fileAst)

	return fileObj
}
