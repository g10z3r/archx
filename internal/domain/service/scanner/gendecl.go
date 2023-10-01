package scanner

import (
	"context"
	"go/ast"
	"go/token"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
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
	structDTO, usedPackages, err := domainDTO.NewStructDTO(s.getFileSet(), params.structType, domainDTO.NotEmbedded, &params.structName)
	if err != nil {
		return err
	}

	for i := 0; i < len(usedPackages); i++ {
		if index := pkgCache.GetImportIndex(usedPackages[i].Alias); index >= 0 {
			structDTO.AddDependency(index, usedPackages[i].Element)
		}
	}

	if index := pkgCache.GetStructIndex(params.structName); index < 0 {
		if err := s.db.PackageRepo().StructRepo().Append(ctx, structDTO, pkgCache.StructsIndexLen(), params.pkgPath); err != nil {
			return err
		}

		pkgCache.AddStructIndex(params.structName)
		return nil
	}

	return nil
}
