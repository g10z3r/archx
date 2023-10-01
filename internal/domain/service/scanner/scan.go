package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
	"github.com/g10z3r/archx/internal/domain/repository"
	"github.com/g10z3r/archx/internal/domain/service/scanner/cache"
	"github.com/g10z3r/archx/pkg/bloom"
)

type scannerCache interface {
	AddPackage(pkgPath string, index int)
	GetPackageIndex(pkgName string) int
	PackagesIndexLen() int
}

type ScanService struct {
	cache scannerCache
	db    repository.ScannerRepository
}

func NewScanService(scanRepo repository.ScannerRepository) *ScanService {
	return &ScanService{
		cache: cache.NewScannerCache(),
		db:    scanRepo,
	}
}

func (s *ScanService) Perform(ctx context.Context, dirPath string, basePath string) {
	if err := s.db.Create(ctx, domainDTO.NewScanResultDTO(basePath)); err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		newPkg := domainDTO.NewPackageDTO(dirPath, pkg.Name)
		if err := s.db.PackageRepo().Append(ctx, newPkg, s.cache.PackagesIndexLen()); err != nil {
			log.Fatal(err)
		}
		s.cache.AddPackage(newPkg.Path, s.cache.PackagesIndexLen())

		pkgImports, total := processImports(pkg.Files)
		pkgCache := cache.NewPackageCache(bloom.FilterConfig{
			ExpectedItemCount:        uint64(total),
			DesiredFalsePositiveRate: 0.01,
		})

		for _, _import := range pkgImports {
			_import.Trim(basePath)

			if isSideEffectImport(_import) {
				contains, err := pkgCache.CheckSideEffectImport([]byte(_import.Path))
				if err != nil {
					log.Fatal(err)
				}

				if contains {
					continue
				}

				if err := s.db.PackageRepo().ImportRepo().AppendSideEffectImport(ctx, _import, newPkg.Path); err != nil {
					log.Fatal(err)
				}

				pkgCache.AddSideEffectImport(_import)
				continue
			}

			contains, err := pkgCache.CheckImport([]byte(_import.Path))
			if err != nil {
				log.Fatal(err)
			}

			if !contains {
				if err := s.db.PackageRepo().ImportRepo().Append(ctx, _import, newPkg.Path); err != nil {
					log.Fatal(err)
				}

				pkgCache.AddImport(_import, pkgCache.ImportsLen())
				continue
			}

			if index := pkgCache.GetImportIndex(_import.Alias); index < 0 {
				for i, imp := range pkgCache.Imports {
					if imp == _import.Path {
						pkgCache.AddImportIndex(_import, i)
					}
				}
			}
		}

		jsonData, err := json.Marshal(pkgCache)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(jsonData))
	}
}

func isSideEffectImport(_import *domainDTO.ImportDTO) bool {
	return _import.WithAlias && _import.Alias == "_"
}

func processImports(files map[string]*ast.File) ([]*domainDTO.ImportDTO, int) {
	var impTotal int
	var imports []*domainDTO.ImportDTO

	for _, file := range files {
		impTotal = impTotal + len(file.Imports)

		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				imports = append(imports, domainDTO.NewImportDTO(imp))
			}
		}
	}

	return imports, impTotal
}
