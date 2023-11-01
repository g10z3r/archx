package analyzer

import (
	"context"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type AnalyzerOld interface {
	Name() string
	Check(node ast.Node) bool
	Analyze(f *obj.FileObj, node ast.Node) obj.Object
	Save(f *obj.FileObj, obj obj.Object)
}

type AnalyzerMapOld map[string]AnalyzerOld

// New

type Analyzer[Input, Output any] interface {
	Analyze(ctx context.Context, i Input) (Output, error)
	Check(i Input) bool // TODO: add context as arg
}

type CheckFunc[Input any] func(node Input) bool
type AnalyzeFunc[Input, Output any] func(ctx context.Context, f *obj.FileObj, i Input) (Output, error)

func NewAnalyzer[Input, Output any](
	file *obj.FileObj,
	analyze AnalyzeFunc[Input, Output],
	check CheckFunc[Input],
) Analyzer[Input, Output] {
	return &analyzer[Input, Output]{
		file,
		analyze,
		check,
	}
}

type analyzer[Input, Output any] struct {
	file    *obj.FileObj
	analyze AnalyzeFunc[Input, Output]
	check   CheckFunc[Input]
}

func (a *analyzer[Input, Output]) Analyze(ctx context.Context, i Input) (Output, error) {
	return a.analyze(ctx, a.file, i)
}

func (a *analyzer[Input, Output]) Check(node Input) bool {
	return a.check(node)
}
