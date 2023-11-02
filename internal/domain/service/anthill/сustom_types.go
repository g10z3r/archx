package anthill

import (
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

const (
	StructNodeType uint = iota + 1
	FuncNodeType
	ImportNodeType
)

type EngineAFMap analyzer.AnalyzerFactoryMap[uint, ast.Node, obj.Object]

func (afm EngineAFMap) Initialize(f *obj.FileObj) analyzer.AnalyzerMap[uint, ast.Node, obj.Object] {
	result := make(analyzer.AnalyzerMap[uint, ast.Node, obj.Object], len(afm))
	for k, v := range afm {
		result[k] = v(f)
	}

	return result
}
