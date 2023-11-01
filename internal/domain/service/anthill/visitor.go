package anthill

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/analyzer"
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type Visitor struct {
	noCopy noCopy

	fileObj      *obj.FileObj
	analyzerMap  analyzer.AnalyzerMapOld
	analyzer2Map map[string]analyzer.Analyzer[ast.Node, obj.Object]
}

func NewVisitor(f *obj.FileObj, analyzerMap analyzer.AnalyzerMapOld, analyzers2 map[string]analyzer.Analyzer[ast.Node, obj.Object]) *Visitor {
	return &Visitor{
		analyzerMap:  analyzerMap,
		analyzer2Map: analyzers2,
		fileObj:      f,
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	// for _, analyzer := range v.analyzerMap {
	// 	if ok := analyzer.Check(node); ok {

	// 		obj := analyzer.Analyze(v.fileObj, node)
	// 		if obj != nil {
	// 			analyzer.Save(v.fileObj, obj) // Add ok return
	// 			break
	// 		}

	// 		log.Fatal("got nil object")
	// 		break

	// 	}
	// }

	for _, analyzer := range v.analyzer2Map {
		if ok := analyzer.Check(node); ok {
			object, err := analyzer.Analyze(context.Background(), node)
			if err != nil {
				fmt.Println(err)
				break
			}

			if object != nil {
				switch o := object.(type) {
				case *obj.ImportObj:
					v.fileObj.AppendImport(o)
				case *obj.FuncObj:
					v.fileObj.AppendFunc(o)
				case *obj.StructObj:
					v.fileObj.AppendStruct(o)
				}

				break
			}

		}
	}
	return v
}
