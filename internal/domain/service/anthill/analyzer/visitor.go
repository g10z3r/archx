package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
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

type VisitorMetadata struct {
	fset     *token.FileSet
	ModName  string
	FileName string
	Imports  *VisitorImports
	Stats    *VisitorStats
	Errors   []error
}

type AnalyzerMap map[string]Analyzer

type Visitor struct {
	metadata  *VisitorMetadata
	analyzers AnalyzerMap
	bucket    map[string][]Object
}

func NewVisitor(fset *token.FileSet, analyzers AnalyzerMap, fileName string) *Visitor {
	return &Visitor{
		analyzers: analyzers,
		metadata: &VisitorMetadata{
			fset:     fset,
			FileName: fileName,
			Imports: &VisitorImports{
				SideEffectImports:  make([]string, 0),
				RegularImports:     make([]string, 0),
				RegularImportsMeta: make(map[string]int),
			},
		},
		bucket: make(map[string][]Object),
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	for _, analyzer := range v.analyzers {
		if ok := analyzer.Check(node); ok {
			obj := analyzer.Analyze(v.metadata, node)
			v.bucket[analyzer.Name()] = append(v.bucket[analyzer.Name()], obj)
			break
		}
	}
	return v
}

func (v *Visitor) Unload(sectionKey string) []Object {
	fmt.Println(v.bucket)
	if data, ok := v.bucket[sectionKey]; ok {
		return data
	}

	return nil
}
