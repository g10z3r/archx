package scaner

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/scaner/buffer"
	"github.com/g10z3r/archx/internal/scaner/entity"
)

func processGenDecl(buf *buffer.BufferEventBus, fs *token.FileSet, genDecl *ast.GenDecl) {
	if genDecl.Tok != token.TYPE {
		return
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		sType, err := processStructType(buf, fs, typeSpec, structType)
		if err != nil {
			errChan <- err
			continue
		}

		notifyStructUpsert(buf, typeSpec.Name.Name, sType)
	}
}

func processStructType(buf *buffer.BufferEventBus, fs *token.FileSet, typeSpec *ast.TypeSpec, structType *ast.StructType) (*entity.StructInfo, error) {
	sType, usedPackages, err := entity.NewStructType(fs, structType, entity.NotEmbedded)
	if err != nil {
		return nil, fmt.Errorf("failed to create new struct type: %w", err)
	}

	for _, p := range usedPackages {
		if importIndex, exists := buf.ImportBuffer.GetIndexByAlias(p.Alias); exists {
			sType.AddDependency(importIndex, p.Element)
		}
	}

	return sType, nil
}

func updateDependencies(buf *buffer.BufferEventBus, sType *entity.StructInfo, usedPackages []entity.UsedPackage) {
	for _, p := range usedPackages {
		if importIndex, exists := buf.ImportBuffer.GetIndexByAlias(p.Alias); exists {
			sType.AddDependency(importIndex, p.Element)
		}
	}
}

func notifyStructUpsert(buf *buffer.BufferEventBus, structName string, sType *entity.StructInfo) {
	buf.SendEvent(
		&buffer.UpsertStructEvent{
			StructInfo: sType,
			StructName: structName,
		},
	)
}
