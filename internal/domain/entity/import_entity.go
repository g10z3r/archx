package entity

import (
	"go/ast"
	"path/filepath"
	"strings"
)

type ImportEntity struct {
	File      string
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

func NewImportEntity(fileName string, importSpec *ast.ImportSpec) *ImportEntity {
	var isWithAlias bool
	var alias string

	if importSpec.Name != nil {
		alias = importSpec.Name.Name
		isWithAlias = true
	}
	return &ImportEntity{
		File:      filepath.Base(fileName),
		Path:      strings.Trim(importSpec.Path.Value, `"`),
		Alias:     alias,
		WithAlias: isWithAlias,
	}
}
