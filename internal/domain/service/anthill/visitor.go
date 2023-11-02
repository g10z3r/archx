package anthill

import (
	"context"
	"fmt"
	"go/ast"
	"sync"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Visitor interface {
	VisitWithContext(ctx context.Context, node ast.Node) (w Visitor)
}

type visitor struct {
	noCopy noCopy

	fileObj *obj.FileObj

	// Used to determine the type of an ast.Node.
	// This function helps identify the specific type of a node within the abstract syntax tree (AST).
	determinator func(ast.Node) uint
	analyzerMap  analyzer.AnalyzerMap[uint, ast.Node, obj.Object]

	once sync.Once
}

type visitorConfig struct {
	file         *obj.FileObj
	alzMap       analyzer.AnalyzerMap[uint, ast.Node, obj.Object]
	determinator func(ast.Node) uint
}

func NewVisitor(cfg visitorConfig) *visitor {
	v := new(visitor)
	v.once.Do(func() {
		v.fileObj = cfg.file
		v.analyzerMap = cfg.alzMap
		v.determinator = cfg.determinator

	})

	return v
}

func (v *visitor) VisitWithContext(ctx context.Context, node ast.Node) Visitor {
	if node == nil {
		return v
	}

	analyzer, ok := v.analyzerMap[v.determinator(node)]
	if !ok {
		return v
	}

	object, err := analyzer.Analyze(ctx, node)
	if err != nil {
		fmt.Println(err) // TODO: decide later how to handle the error
		return v
	}

	if err := v.fileObj.Save(object); err != nil {
		fmt.Println(err) // TODO: decide later how to handle the error
	}

	return v
}
