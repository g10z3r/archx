package entity

import (
	"go/ast"
	"strings"
)

type ImportEntity struct {
	Path      string
	Alias     string
	WithAlias bool
}

func (e *ImportEntity) Trim(basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	e.Path = "/" + strings.TrimPrefix(e.Path, basePath)
}

func NewImportEntity(importSpec *ast.ImportSpec) *ImportEntity {
	var isWithAlias bool
	var alias string

	if importSpec.Name != nil {
		alias = importSpec.Name.Name
		isWithAlias = true
	}
	return &ImportEntity{
		Path:      strings.Trim(importSpec.Path.Value, `"`),
		Alias:     alias,
		WithAlias: isWithAlias,
	}
}
