package common

import (
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
