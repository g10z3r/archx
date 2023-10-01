package scanner

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"sync"

	"github.com/g10z3r/archx/internal/domain/repository"
	"github.com/g10z3r/archx/internal/domain/service/scanner/cache"
	"github.com/g10z3r/archx/pkg/bloom"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type scannerCache interface {
	AddPackage(pkgPath string, index int)
	GetPackageIndex(pkgName string) int
	PackagesIndexLen() int
}

type packageCache interface {
	ImportsLen() int
	CheckImport(b []byte) (bool, error)
	AddImport(_import *domainDTO.ImportDTO, index int)
	AddImportIndex(_import *domainDTO.ImportDTO, index int)
	GetImportIndex(importAlias string) int
	GetImports() []string
	CheckSideEffectImport(b []byte) (bool, error)
	AddSideEffectImport(_import *domainDTO.ImportDTO)

	StructsIndexLen() int
	AddStructIndex(structName string)
	GetStructIndex(structName string) int

	Debug()
}

type ScanService struct {
	mu sync.RWMutex

	_fset *token.FileSet

	cache scannerCache
	db    repository.ScannerRepository
}

func (s *ScanService) getFileSet() *token.FileSet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s._fset
}

func NewScanService(scanRepo repository.ScannerRepository) *ScanService {
	return &ScanService{
		_fset: token.NewFileSet(),
		cache: cache.NewScannerCache(),
		db:    scanRepo,
	}
}

func (s *ScanService) Perform(ctx context.Context, dirPath string, basePath string) {
	if err := s.db.Register(ctx, domainDTO.NewScanResultDTO(basePath)); err != nil {
		log.Fatal(err)
	}

	pkgs, err := parser.ParseDir(s._fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		newPkg, err := s.registerNewPackage(ctx, dirPath, pkg.Name)
		if err != nil {
			log.Fatal(err)
		}

		pkgImports, total := fetchPackageImports(pkg.Files)
		pkgCache := cache.NewPackageCache(bloom.FilterConfig{
			ExpectedItemCount:        uint64(total),
			DesiredFalsePositiveRate: 0.01,
		})

		for _, _import := range pkgImports {
			pid := pkgImportData{
				pkgCache: pkgCache,
				newPkg:   newPkg,
				basePath: basePath,
			}

			if err := s.processPackageImport(ctx, pid, _import); err != nil {
				log.Fatal(err)
			}
		}

		var wg sync.WaitGroup

		for fileName, file := range pkg.Files {
			wg.Add(1)
			go func(file *ast.File, fileName string) {
				defer wg.Done()

				log.Printf("Processing file: %s", fileName)

				for _, decl := range file.Decls {
					switch d := decl.(type) {
					case *ast.GenDecl:
						if err = s.processGenDecl(ctx, d, pkgCache, newPkg.Path); err != nil {
							log.Fatal(err)
						}
					}
				}
			}(file, fileName)
		}

		wg.Wait()
		pkgCache.Debug()
	}
}

func (s *ScanService) registerNewPackage(ctx context.Context, dirPath, pkgName string) (*domainDTO.PackageDTO, error) {
	newPkg := domainDTO.NewPackageDTO(dirPath, pkgName)
	if err := s.db.PackageRepo().Append(ctx, newPkg, s.cache.PackagesIndexLen()); err != nil {
		return nil, err
	}
	s.cache.AddPackage(newPkg.Path, s.cache.PackagesIndexLen())

	return newPkg, nil
}

type pkgImportData struct {
	pkgCache packageCache
	newPkg   *domainDTO.PackageDTO
	basePath string
}

func (s *ScanService) processPackageImport(ctx context.Context, data pkgImportData, _import *domainDTO.ImportDTO) error {
	_import.Trim(data.basePath)

	if isSideEffectImport(_import) {
		contains, err := data.pkgCache.CheckSideEffectImport([]byte(_import.Path))
		if err != nil {
			return err
		}

		if contains {
			return nil
		}

		if err := s.db.PackageRepo().ImportRepo().AppendSideEffectImport(ctx, _import, data.newPkg.Path); err != nil {
			return err
		}

		data.pkgCache.AddSideEffectImport(_import)
		return nil
	}

	contains, err := data.pkgCache.CheckImport([]byte(_import.Path))
	if err != nil {
		return err
	}

	if !contains {
		if err := s.db.PackageRepo().ImportRepo().Append(ctx, _import, data.newPkg.Path); err != nil {
			return err
		}

		data.pkgCache.AddImport(_import, data.pkgCache.ImportsLen())
		return nil
	}

	if index := data.pkgCache.GetImportIndex(_import.Alias); index < 0 {
		for i, imp := range data.pkgCache.GetImports() {
			if imp == _import.Path {
				data.pkgCache.AddImportIndex(_import, i)
			}
		}
	}

	return nil
}

func fetchPackageImports(files map[string]*ast.File) ([]*domainDTO.ImportDTO, int) {
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

func isSideEffectImport(_import *domainDTO.ImportDTO) bool {
	return _import.WithAlias && _import.Alias == "_"
}
