package anthill

import (
	"go/ast"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"
)

type Visitor struct {
	fileObj      *obj.FileObj
	analyzerMap  analyzer.AnalyzerMap
	analyzer2Map map[string]analyzer.Analyzer2[ast.Node, analyzer.Object]
}

func NewVisitor(f *obj.FileObj, analyzerMap analyzer.AnalyzerMap, analyzers2 map[string]analyzer.Analyzer2[ast.Node, analyzer.Object]) *Visitor {
	return &Visitor{
		analyzerMap:  analyzerMap,
		analyzer2Map: analyzers2,
		fileObj:      f,
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

	// for _, analyzer := range v.analyzer2Map {
	// 	if ok := analyzer.Check(node); ok {
	// 		o, err := analyzer.Analyze(context.Background(), v.fileObj, node)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			break
	// 		}

	// 		if o != nil {
	// 			v.fileObj.AppendImport(o.(*obj.ImportObj))
	// 			break
	// 		}

	// 	}
	// }
	return v
}
