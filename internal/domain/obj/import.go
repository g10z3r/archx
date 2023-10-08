package obj

import (
	"go/ast"
	"path/filepath"
	"strings"
)

type ImportObj struct {
	File      string
	Path      string
	Alias     string
	WithAlias bool
}

func (e *ImportObj) CheckAndTrim(modName string) bool {
	if len(e.Path) < len(modName) {
		return false
	}

	if !strings.HasPrefix(e.Path, modName) {
		return false
	}

	e.Path = e.Path[len(modName):]
	return true
}

func NewImportObj(fileName string, importSpec *ast.ImportSpec) *ImportObj {
	var isWithAlias bool
	var alias string

	if importSpec.Name != nil {
		alias = importSpec.Name.Name
		isWithAlias = true
	}
	return &ImportObj{
		File:      filepath.Base(fileName),
		Path:      strings.Trim(importSpec.Path.Value, `"`),
		Alias:     alias,
		WithAlias: isWithAlias,
	}
}
