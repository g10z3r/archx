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

func (dto *ImportDTO) Trim(basePath string) {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	dto.Path = "/" + strings.TrimPrefix(dto.Path, basePath)
}

func NewImportDTO(importSpec *ast.ImportSpec) *ImportDTO {
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
