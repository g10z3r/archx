package scanner

import (
	"context"
	"go/ast"
	"go/token"

	"github.com/g10z3r/archx/internal/domain/entity"
)

func (s *ScanService) processGenDecl(ctx context.Context, genDecl *ast.GenDecl, pkgCache packageCache, pkgPath string) error {
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
			s.processStructType(ctx, pkgCache, structProcessingParams{
				typeSpec:   typeSpec,
				structType: t,
				pkgPath:    pkgPath,
				structName: typeSpec.Name.Name,
			})
		}

	}

	return nil
}

type structProcessingParams struct {
	typeSpec   *ast.TypeSpec
	structType *ast.StructType
	pkgPath    string
	structName string
}

func (s *ScanService) processStructType(ctx context.Context, pkgCache packageCache, params structProcessingParams) error {
	structEntity, usedPackages, err := entity.NewStructEntity(s.getFileSet(), params.structType, entity.NotEmbedded, &params.structName)
	if err != nil {
		return err
	}

	for i := 0; i < len(usedPackages); i++ {
		if index := pkgCache.GetImportIndex(usedPackages[i].Alias); index >= 0 {
			structEntity.AddDependency(index, usedPackages[i].Element)
		}
	}

	if index := pkgCache.GetStructIndex(params.structName); index < 0 {
		indexInCache := pkgCache.AddStructIndex(params.structName)
		if err := s.db.PackageRepo().StructRepo().Append(ctx, structEntity, indexInCache, params.pkgPath); err != nil {
			return err
		}

		return nil
	}

	return nil
}
