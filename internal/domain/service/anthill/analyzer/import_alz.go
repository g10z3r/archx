package analyzer

import (
	"go/ast"
	"log"
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
		// f.Imports.RegularImportsMeta[importObj.Alias] = len(f.Imports.RegularImports)
		// f.Imports.RegularImports = append(f.Imports.RegularImports, importObj.Path)
	case obj.ImportTypeRegular:
		f.Imports.RegularImportsMeta[importObj.Alias] = len(f.Imports.RegularImports)
		f.Imports.RegularImports = append(f.Imports.RegularImports, importObj.Path)
	case obj.ImportTypeSideEffect:
		f.Imports.SideEffectImports = append(f.Imports.SideEffectImports, importObj.Path)
	}
}

func (a *ImportAnalyzer) Analyze(f *obj.FileObj, node ast.Node) Object {
	_import, _ := node.(*ast.ImportSpec)

	if _import.Path == nil && _import.Path.Value == "" {
		return nil
	}

	path := strings.Trim(_import.Path.Value, `"`)
	if !strings.HasPrefix(path, f.Metadata.Module) {
		return nil
	}

	if _import.Name != nil && _import.Name.Name == "_" {
		return obj.NewImportObj(_import, obj.ImportTypeSideEffect)
	}

	if !strings.HasPrefix(_import.Path.Value, f.Metadata.Module) {
		return obj.NewImportObj(_import, obj.ImportTypeRegular)
	}

	return obj.NewImportObj(_import, obj.ImportTypeInternal)
}
