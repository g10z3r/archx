package anthill

import (
	"go/ast"
	"reflect"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

const (
	StructNodeType uint = iota + 1
	FuncNodeType
	ImportNodeType
)

type EngineAFMap analyzer.AnalyzerFactoryMap[reflect.Type, ast.Node, obj.Object]

func (afm EngineAFMap) Initialize(f *obj.FileObj) analyzer.AnalyzerMap[reflect.Type, ast.Node, obj.Object] {
	result := make(analyzer.AnalyzerMap[reflect.Type, ast.Node, obj.Object], len(afm))
	for k, v := range afm {
		result[k] = v(f)
	}

	return result
}
