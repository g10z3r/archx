package analyzer

import (
	"fmt"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Visitor struct {
	fileObj     *obj.FileObj
	analyzerMap map[string]Analyzer
}

func NewVisitor(f *obj.FileObj, analyzers map[string]Analyzer) *Visitor {
	return &Visitor{
		analyzerMap: analyzers,
		fileObj:     f,
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	for _, analyzer := range v.analyzerMap {
		if ok := analyzer.Check(node); ok {

			obj := analyzer.Analyze(v.fileObj, node)
			fmt.Println(analyzer.Name())
			analyzer.Save(v.fileObj, obj) // Add ok return

			break
		}
	}
	return v
}
