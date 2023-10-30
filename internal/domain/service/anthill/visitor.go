package anthill

import (
	"go/ast"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
)

type Visitor struct {
	fileObj     *obj.FileObj
	analyzerMap analyzer.AnalyzerMap
}

func NewVisitor(f *obj.FileObj, analyzers analyzer.AnalyzerMap) *Visitor {
	return &Visitor{
		analyzerMap: analyzers,
		fileObj:     f,
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	for _, analyzer := range v.analyzerMap {
		if ok := analyzer.Check(node); ok {
			obj := analyzer.Analyze(v.fileObj, node)
			if obj != nil {
				analyzer.Save(v.fileObj, obj) // Add ok return
				break
			}

			log.Fatal("got nil object")
			break

		}
	}
	return v
}