package analyzer

import (
	"context"
	"errors"
	"go/ast"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

func NewStructAnalyzer(file *obj.FileObj) Analyzer[ast.Node, obj.Object] {
	return NewAnalyzer[ast.Node, obj.Object](
		file,
		analyzeStructNode,
	)
}

// func checkStructNode(node ast.Node) bool {
// 	typeSpec, ok := node.(*ast.TypeSpec)
// 	if !ok {
// 		return false
// 	}

// 	_, ok = typeSpec.Type.(*ast.StructType)
// 	if !ok {
// 		return false
// 	}

// 	return true
// }

func analyzeStructNode(ctx context.Context, f *obj.FileObj, node ast.Node) (obj.Object, error) {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return nil, errors.New("some error from analyzeStructNode 1") // TODO: add normal error return message
	}

	t, ok := node.(*ast.TypeSpec).Type.(*ast.StructType)
	if !ok {
		return nil, errors.New("some error from analyzeStructNode 2") // TODO: add normal error return message
	}

	structObj, usedPackages, err := obj.NewStructObj(f.FileSet, t, obj.NotEmbedded, &typeSpec.Name.Name)
	if err != nil {
		return nil, errors.New("some error from analyzeStructNode 3") // TODO: add normal error return message
	}

	for _, pkg := range usedPackages {
		if index, exists := f.Entities.Imports.InternalImportsMeta[pkg.Alias]; exists {
			structObj.AddDependency(index, pkg.Element)
		}
	}

	return structObj, nil
}
