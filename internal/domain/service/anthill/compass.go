package anthill

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

var alz = map[string]analyzer.Analyzer{}

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
	fset    *token.FileSet
	manager *Manager

	eventCh       chan compassEvent
	unsubscribeCh chan struct{}
}

func NewCompass() *Compass {
	return &Compass{
		fset: token.NewFileSet(),
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

func (c *Compass) Parse() {
	importAlz := &analyzer.ImportAnalyzer{}
	c.manager.Register(importAlz)

	fmt.Println(len(c.manager.analyzers))

	structAlz := &analyzer.StructAnalyzer{}
	c.manager.Register(structAlz)

	// funcAlz := &analyzer.FunctionAnalyzer{}
	// c.manager.Register(funcAlz)

	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, "./example/cmd", nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	pkgObj := &obj.PackageObj{}
	for _, p := range pkg {
		for _, f := range p.Files {
			fileObj := obj.NewFileObj(fset, "github.com/g10z3r/archx", f.Name.Name)
			pkgObj.AppendFile(fileObj)

			vis := analyzer.NewVisitor(fileObj, c.manager.analyzers)
			ast.Walk(vis, f)
		}
	}

	c.eventCh <- &event.PackageFormedEvent{
		Package: pkgObj,
	}
}

func toPkgStructs(data []analyzer.Object) []*obj.StructObj {
	dataOutput := make([]*obj.StructObj, 0, len(data))

	for _, structObj := range data {
		dataOutput = append(dataOutput, structObj.(*obj.StructObj))
	}

	return dataOutput
}

func toPkgImports(data []analyzer.Object) []string {
	dataOutput := make([]string, 0, len(data))
	for _, importObj := range data {
		dataOutput = append(dataOutput, importObj.(*obj.ImportObj).Path)
	}

	return dataOutput
}
