package analyzer

import (
	"context"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
)

type Object interface {
	Type() string
}

type Analyzer interface {
	Name() string
	Check(node ast.Node) bool
	Analyze(f *obj.FileObj, node ast.Node) Object
	Save(f *obj.FileObj, obj Object)
}

type AnalyzerMap map[string]Analyzer

type Analyzer2[Input ast.Node, Output Object] interface {
	Analyze(ctx context.Context, f *obj.FileObj, i Input) (Output, error)
	Check(node ast.Node) bool
	// Cancel(i Input, err error)
}

func NewAnalyzer[Input ast.Node, Output Object](
	analyze func(ctx context.Context, f *obj.FileObj, i Input) (Output, error),
	check func(node ast.Node) bool,
	// cancel func(i Input, err error),
) Analyzer2[Input, Output] {
	return &analyzer[Input, Output]{analyze, check}
}

type analyzer[Input ast.Node, Output Object] struct {
	analyze func(ctx context.Context, f *obj.FileObj, i Input) (Output, error)
	check   func(node ast.Node) bool
	// cancel  func(i Input, err error)
}

func (a *analyzer[Input, Output]) Analyze(ctx context.Context, f *obj.FileObj, i Input) (Output, error) {
	return a.analyze(ctx, f, i)
}

func (a *analyzer[Input, Output]) Check(node ast.Node) bool {
	return a.check(node)
}

// func (a *analyzer[Input, Output]) Cancel(i Input, err error) {
// 	a.cancel(i, err)
// }
