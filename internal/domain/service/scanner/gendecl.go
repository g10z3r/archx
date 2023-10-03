package scanner

import (
	"context"
	"go/ast"
	"go/token"
	"log"

	"github.com/g10z3r/archx/internal/domain/entity"
)

func (pa *packageActor) processGenDecl(ctx context.Context, genDecl *ast.GenDecl) error {
	if genDecl.Tok != token.TYPE {
		return nil
	}

	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch t := typeSpec.Type.(type) {
		case *ast.StructType:
			return pa.processStructType(ctx, &structProcessingParams{
				typeSpec:   typeSpec,
				structType: t,
				structName: typeSpec.Name.Name,
			})
		}

	}

	return nil
}

type structProcessingParams struct {
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
	// pkgPath    string
	structName string
}

func (pa *packageActor) processStructType(ctx context.Context, params *structProcessingParams) error {
	structEntity, usedPackages, err := entity.NewStructEntity(pa.FileSet(), params.structType, entity.NotEmbedded, &params.structName)
	if err != nil {
		return err
	}

	for _, pkg := range usedPackages {
		if index := pa.cache.GetImportIndex(pkg.Alias); index >= 0 {
			structEntity.AddDependency(index, pkg.Element)
		}
	}

	pa.mu.Lock()
	defer pa.mu.Unlock()

	if index := pa.cache.GetStructIndex(params.structName); index < 0 {
		// sync with buffer
		for _, method := range pa.buf.GetAndClearMethods(params.structName) {
			structEntity.AddMethod(method, method.Name)
			log.Printf("Syncing method %s", method.Name)

			depsLen := len(structEntity.DependenciesIndex)
			for dep, i := range method.DependenciesIndex {
				structEntity.Dependencies = append(structEntity.Dependencies, method.Dependencies[i])
				structEntity.DependenciesIndex[dep] = depsLen + i
			}
		}

		index := pa.cache.AddStructIndex(params.structName)
		return pa.db.StructAcc().Append(ctx, structEntity, index, pa.pkg.Path)
	}

	return nil
}
