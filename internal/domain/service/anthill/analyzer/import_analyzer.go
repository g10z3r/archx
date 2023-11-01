package analyzer

import (
	"context"
	"go/ast"
	"log"
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

func NewImportAnalyzer(file *obj.FileObj) Analyzer[ast.Node, obj.Object] {
	return NewAnalyzer[ast.Node, obj.Object](
		file,
		analyzeImportNode,
		checkImportNode,
	)
}

type ImportAnalyzer struct{}

func (a *ImportAnalyzer) Name() string {
	return "import"
}

func (a *ImportAnalyzer) Check(node ast.Node) bool {
	_, ok := node.(*ast.ImportSpec)
	return ok
}

func (a *ImportAnalyzer) Save(f *obj.FileObj, object obj.Object) {
	importObj, ok := object.(*obj.ImportObj)
	if !ok {
		log.Fatal("not a import objects")
	}

	f.AppendImport(importObj)
}

func getAlias(importObj *obj.ImportObj) string {
	if importObj.WithAlias {
		return importObj.Alias
	}

	return path.Base(importObj.Path)
}

func (a *ImportAnalyzer) Analyze(f *obj.FileObj, node ast.Node) obj.Object {
	importSpec, _ := node.(*ast.ImportSpec)

	if importSpec.Path == nil && importSpec.Path.Value == "" {
		return nil
	}

	path := strings.Trim(importSpec.Path.Value, `"`)
	if !strings.HasPrefix(path, f.Metadata.Module) {
		return obj.NewImportObj(importSpec, obj.ImportTypeExternal)
	}

	if importSpec.Name != nil && importSpec.Name.Name == "_" {
		return obj.NewImportObj(importSpec, obj.ImportTypeSideEffect)
	}

	return obj.NewImportObj(importSpec, obj.ImportTypeInternal)
}

func analyzeImportNode(ctx context.Context, f *obj.FileObj, node ast.Node) (obj.Object, error) {
	importSpec, _ := node.(*ast.ImportSpec)

	if importSpec.Path == nil && importSpec.Path.Value == "" {
		return nil, nil // TODO: add error return message
	}

	path := strings.Trim(importSpec.Path.Value, `"`)
	if !strings.HasPrefix(path, f.Metadata.Module) {
		return obj.NewImportObj(importSpec, obj.ImportTypeExternal), nil
	}

	if importSpec.Name != nil && importSpec.Name.Name == "_" {
		return obj.NewImportObj(importSpec, obj.ImportTypeSideEffect), nil
	}

	return obj.NewImportObj(importSpec, obj.ImportTypeInternal), nil
}

func checkImportNode(node ast.Node) bool {
	_, ok := node.(*ast.ImportSpec)
	return ok
}
