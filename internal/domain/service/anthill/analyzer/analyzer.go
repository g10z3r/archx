package analyzer

import (
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Object interface {
	Type() string
}

type Context interface{}

type Analyzer interface {
	Name() string
	Check(node ast.Node) bool
	Analyze(f *obj.FileObj, node ast.Node) Object
	Save(f *obj.FileObj, obj Object)
}
