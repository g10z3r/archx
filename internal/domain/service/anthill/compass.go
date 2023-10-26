package anthill

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Manager struct {
	analyzers map[string]analyzer.Analyzer
}

func (m *Manager) Register(a analyzer.Analyzer) {
	m.analyzers[a.Name()] = a
}

type compassEvent interface {
	Name() string
}

type Compass struct {
	manager *Manager

	eventCh       chan compassEvent
	unsubscribeCh chan struct{}
}

func NewCompass() *Compass {
	return &Compass{
		manager: &Manager{
			analyzers: map[string]analyzer.Analyzer{},
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

	for _, pkgAst := range pkg {
		pkgObj := obj.NewPackageObj(pkgAst, targetDir)

		for fileName, fileAst := range pkgAst.Files {
			fileObj := obj.NewFileObj(fset, info.ModuleName, filepath.Base(fileName))
			pkgObj.AppendFile(fileObj)

			vis := analyzer.NewVisitor(fileObj, c.manager.analyzers)
			ast.Walk(vis, fileAst)
		}

		c.eventCh <- &event.PackageFormedEvent{
			Package: pkgObj,
		}
	}
}
