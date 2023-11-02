package anthill

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"sync"

	"github.com/g10z3r/archx/internal/domain/service/anthill/collector"
	"github.com/g10z3r/archx/internal/domain/service/anthill/event"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Engine struct {
	noCopy noCopy

	afMap EngineAFMap

	// Used to determine the type of an ast.Node.
	// This function helps identify the specific type of a node within the abstract syntax tree (AST).
	determinator func(ast.Node) uint

	mutex sync.RWMutex
	once  sync.Once

	eventCh       chan event.Event
	unsubscribeCh chan struct{}
}

func NewEngine(alzFactoryMap EngineAFMap, determinator func(ast.Node) uint) *Engine {
	engine := new(Engine)
	engine.once.Do(func() {
		engine.afMap = alzFactoryMap
		engine.determinator = determinator
	})

	return engine
}

func (e *Engine) Subscribe(eventCh chan event.Event) chan struct{} {
	e.once.Do(func() {
		e.eventCh = eventCh
		e.unsubscribeCh = make(chan struct{})
	})

	return e.unsubscribeCh
}

// func (c *Compass) Run(ctx context.Context) {
// 	res := c.pipeline.Run(ctx, &plugin.CollectorPluginInput{
// 		RootDir:     ".",
// 		IgnoredList: config.DefaultIgnoredMap,
// 	})

// 	for _, dir := range res.([]string) {
// 		fmt.Println(dir)
// 	}
// }

func (e *Engine) ParseDir(info *collector.Info, targetDir string) ([]*obj.PackageObj, error) {
	fset := token.NewFileSet()
	pkg, err := parser.ParseDir(fset, targetDir, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	type resultDTO struct {
		sync.Mutex
		data []*obj.PackageObj
	}

	result := new(resultDTO)

	var wg sync.WaitGroup
	for _, pkgAst := range pkg {
		wg.Add(1)
		go func(pkgAst *ast.Package) {
			result.Lock()
			result.data = append(result.data, e.parsePkg(fset, pkgAst, targetDir, info.ModuleName))
			result.Unlock()

			wg.Done()
		}(pkgAst)
	}

	wg.Wait()

	return result.data, nil
}

func (e *Engine) parsePkg(fset *token.FileSet, pkgAst *ast.Package, targetDir, moduleName string) *obj.PackageObj {
	pkgObj := obj.NewPackageObj(pkgAst, targetDir)

	var wg sync.WaitGroup
	for fileName, fileAst := range pkgAst.Files {
		wg.Add(1)
		go func(fileAst *ast.File, fileName string) {
			fileObj := e.parseFile(fset, fileAst, moduleName, fileName)
			pkgObj.AppendFile(fileObj)

			wg.Done()
		}(fileAst, fileName)
	}

	wg.Wait()
	return pkgObj
}

func (e *Engine) parseFile(fset *token.FileSet, fileAst *ast.File, moduleName, fileName string) *obj.FileObj {
	fileObj := obj.NewFileObj(fset, moduleName, filepath.Base(fileName))
	visitor := NewVisitor(visitorConfig{
		file:         fileObj,
		alzMap:       e.afMap.Initialize(fileObj),
		determinator: e.determinator,
	})

	WalkWithContext(context.Background(), visitor, fileAst)
	return fileObj
}
