package entity

import (
	"go/ast"
	"strings"
)

type Import struct {
	Path      string
	Alias     string
	WithAlias bool
}

func NewImport(importSpec *ast.ImportSpec) *Import {
	var isWithAlias bool
	var alias string

	if importSpec.Name != nil {
		alias = importSpec.Name.Name
		isWithAlias = true
	}
	return &Import{
		Path:      strings.Trim(importSpec.Path.Value, `"`),
		Alias:     alias,
		WithAlias: isWithAlias,
	}
}
