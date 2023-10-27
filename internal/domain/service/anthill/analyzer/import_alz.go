package analyzer

import (
	"go/ast"
	"log"
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type ImportAnalyzer struct{}

func (a *ImportAnalyzer) Name() string {
	return "import"
}

func (a *ImportAnalyzer) Check(node ast.Node) bool {
	_, ok := node.(*ast.ImportSpec)
	return ok
}

func (a *ImportAnalyzer) Save(f *obj.FileObj, object Object) {
	importObj, ok := object.(*obj.ImportObj)
	if !ok {
		log.Fatal("not a import objects")
	}

	switch importObj.ImportType {
	case obj.ImportTypeInternal:
		f.Entities.Imports.InternalImportsMeta[getAlias(importObj)] = len(f.Entities.Imports.InternalImports)
		f.Entities.Imports.InternalImports = append(f.Entities.Imports.InternalImports, importObj.Path[len(f.Metadata.Module):])

	case obj.ImportTypeExternal:
		f.Entities.Imports.ExternalImports = append(f.Entities.Imports.ExternalImports, importObj.Path)

	case obj.ImportTypeSideEffect:
		f.Entities.Imports.SideEffectImports = append(f.Entities.Imports.SideEffectImports, importObj.Path)
	}
}

func getAlias(importObj *obj.ImportObj) string {
	if importObj.WithAlias {
		return importObj.Alias
	}

	return path.Base(importObj.Path)
}

func (a *ImportAnalyzer) Analyze(f *obj.FileObj, node ast.Node) Object {
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
