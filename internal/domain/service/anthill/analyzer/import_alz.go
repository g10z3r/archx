package analyzer

import (
	"go/ast"
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

func (a *ImportAnalyzer) Analyze(vtx *VisitorMetadata, node ast.Node) Object {
	_import, _ := node.(*ast.ImportSpec)

	if _import.Path == nil && _import.Path.Value == "" {
		return nil
	}

	path := strings.Trim(_import.Path.Value, `"`)
	if !strings.HasPrefix(path, vtx.ModName) {
		return nil
	}

	if _import.Name != nil && _import.Name.Name == "_" {
		return obj.NewImportObj(_import, obj.ImportTypeSideEffect)
	}

	if !strings.HasPrefix(_import.Path.Value, vtx.ModName) {
		return obj.NewImportObj(_import, obj.ImportTypeRegular)
	}

	return obj.NewImportObj(_import, obj.ImportTypeInternal)
}
