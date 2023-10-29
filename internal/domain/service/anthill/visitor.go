package anthill

import (
	"go/ast"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
	"github.com/g10z3r/archx/internal/domain/service/anthill/common"
)

type Visitor struct {
	fileObj     *obj.FileObj
	analyzerMap common.AnalyzerMap
}

func NewVisitor(f *obj.FileObj, analyzers common.AnalyzerMap) *Visitor {
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
