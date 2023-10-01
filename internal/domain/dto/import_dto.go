package dto

import (
	"go/ast"
	"strings"
)

type ImportDTO struct {
	Path      string
	Alias     string
	WithAlias bool
}

func NewImport(importSpec *ast.ImportSpec) *ImportDTO {
	var isWithAlias bool
	var alias string

	if importSpec.Name != nil {
		alias = importSpec.Name.Name
		isWithAlias = true
	}
	return &ImportDTO{
		Path:      strings.Trim(importSpec.Path.Value, `"`),
		Alias:     alias,
		WithAlias: isWithAlias,
	}
}
