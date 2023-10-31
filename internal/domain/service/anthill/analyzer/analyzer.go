package analyzer

import (
	"context"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Object interface {
	Type() string
}

type AnalyzerOld interface {
	Name() string
	Check(node ast.Node) bool
	Analyze(f *obj.FileObj, node ast.Node) Object
	Save(f *obj.FileObj, obj Object)
}

type AnalyzerMapOld map[string]AnalyzerOld

// New

type Analyzer[Input ast.Node, Output Object] interface {
	Analyze(ctx context.Context, i Input) (Output, error)
	Check(node ast.Node) bool
}

type CheckFunc func(node ast.Node) bool
type AnalyzeFunc[Input ast.Node, Output Object] func(ctx context.Context, f *obj.FileObj, i Input) (Output, error)

func NewAnalyzer[Input ast.Node, Output Object](
	file *obj.FileObj,
	analyze AnalyzeFunc[Input, Output],
	check CheckFunc,
) Analyzer[Input, Output] {
	return &analyzer[Input, Output]{
		file,
		analyze,
		check,
	}
}

type analyzer[Input ast.Node, Output Object] struct {
	file    *obj.FileObj
	analyze AnalyzeFunc[Input, Output]
	check   CheckFunc
}

func (a *analyzer[Input, Output]) Analyze(ctx context.Context, i Input) (Output, error) {
	return a.analyze(ctx, a.file, i)
}

func (a *analyzer[Input, Output]) Check(node ast.Node) bool {
	return a.check(node)
}
