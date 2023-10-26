package analyzer

import (
	"go/ast"

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

func (a *StructAnalyzer) Analyze(vtx *VisitorMetadata, spec ast.Node) Object {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	t, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	structObj, usedPackages, err := obj.NewStructObj(vtx.fset, t, obj.NotEmbedded, &typeSpec.Name.Name)
	if err != nil {
		return nil
	}

	for _, pkg := range usedPackages {
		if index, exists := vtx.Imports.RegularImportsMeta[pkg.Alias]; exists {
			structObj.AddDependency(index, pkg.Element)
		}
	}

	return structObj
}
