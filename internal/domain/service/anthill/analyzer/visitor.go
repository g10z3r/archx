package analyzer

import (
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/pkg/dsl"
)

type VisitorImports struct {
	SideEffectImports  []string
	RegularImports     []string
	RegularImportsMeta map[string]int
}

type VisitorStats struct {
	Functions,
	Structs,
	Interfaces int
}

type VisitorContext struct {
	fset     *token.FileSet
	ModName  string
	FileName string
	Imports  *VisitorImports
	Stats    *VisitorStats
	Errors   []error
}

type AnalyzerMap map[string]Analyzer

type Visitor struct {
	context   *VisitorContext
	analyzers AnalyzerMap
	bucket    dsl.Map[string, []Object]
}

func NewVisitor(fset *token.FileSet, analyzers AnalyzerMap, fileName string) *Visitor {
	return &Visitor{
		context: &VisitorContext{
			fset:     fset,
			FileName: fileName,
			Imports: &VisitorImports{
				SideEffectImports:  make([]string, 0),
				RegularImports:     make([]string, 0),
				RegularImportsMeta: make(map[string]int),
			},
		},
		analyzers: analyzers,
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	for _, analyzer := range v.analyzers {
		if ok := analyzer.Check(node); ok {
			res := analyzer.Analyze(v.context, node)

			bucket, _ := v.bucket.Load(analyzer.Name())
			v.bucket.Store(analyzer.Name(), append(bucket, res))
			break
		}
	}
	return v
}
