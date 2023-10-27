package analyzer

import (
	"go/ast"
	"log"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type StructAnalyzer struct{}

func (a *StructAnalyzer) Name() string {
	return "struct"
}

func (a *StructAnalyzer) Check(node ast.Node) bool {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return false
	}

	_, ok = typeSpec.Type.(*ast.StructType)
	if !ok {
		return false
	}

	return true
}

func (a *StructAnalyzer) Save(f *obj.FileObj, object Object) {
	structObj, ok := object.(*obj.StructObj)
	if !ok {
		log.Fatal("not a struct objects")
	}

	f.Entities.AppendStruct(structObj)
}

func (a *StructAnalyzer) Analyze(f *obj.FileObj, spec ast.Node) Object {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	t, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	structObj, usedPackages, err := obj.NewStructObj(f.FileSet, t, obj.NotEmbedded, &typeSpec.Name.Name)
	if err != nil {
		return nil
	}

	for _, pkg := range usedPackages {
		if index, exists := f.Entities.Imports.InternalImportsMeta[pkg.Alias]; exists {
			structObj.AddDependency(index, pkg.Element)
		}
	}

	return structObj
}
