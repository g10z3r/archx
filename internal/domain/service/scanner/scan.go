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

	"github.com/g10z3r/archx/internal/domain/entity"
)

type scannerCache interface {
	AddPackage(pkgPath string, index int)
	GetPackageIndex(pkgName string) int
	PackagesIndexLen() int
}

type Scanner struct {
	mu sync.RWMutex

	_fset *token.FileSet

	cache scannerCache
	db    repository.SnapshotRepository
}

func NewScanner(scanRepo repository.SnapshotRepository) *Scanner {
	return &Scanner{
		_fset: token.NewFileSet(),
		cache: cache.NewScannerCache(),
		db:    scanRepo,
	}
}

func (s *Scanner) Perform(ctx context.Context, dirPath string, basePath string) {
	if err := s.db.Register(ctx, entity.NewSnapshotEntity(basePath)); err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		newPkg, err := s.registerNewPackage(ctx, dirPath, pkg.Name)
		if err != nil {
			log.Fatal(err)
		}

		pkgActor := newPackageActor(s.db.PackageAcc(), fset, newPkg)
		pkgImports := pkgActor.IndexationImports(pkg)

		if err := pkgActor.ProcessAllImports(ctx, pkgImports, basePath); err != nil {
			log.Fatal(err)
		}

		var pkgWaitGroup sync.WaitGroup
		for fileName, file := range pkg.Files {
			pkgWaitGroup.Add(1)

			go func(file *ast.File, fileName string) {
				defer pkgWaitGroup.Done()

				if err := pkgActor.ScanFile(ctx, fileName, file); err != nil {
					log.Fatal(err)
				}

			}(file, fileName)
		}

		pkgWaitGroup.Wait()
	}
}

func (s *Scanner) registerNewPackage(ctx context.Context, dirPath, pkgName string) (*entity.PackageEntity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newPkg := entity.NewPackageEntity(dirPath, pkgName)
	if err := s.db.PackageAcc().Append(ctx, newPkg, s.cache.PackagesIndexLen()); err != nil {
		return nil, err
	}
	s.cache.AddPackage(newPkg.Path, s.cache.PackagesIndexLen())

	return newPkg, nil
}

// type pkgImportData struct {
// 	pkgCache packageCache
// 	newPkg   *entity.PackageEntity
// 	basePath string
// }

// func (s *Scanner) processPackageImport(ctx context.Context, data pkgImportData, _import *entity.ImportEntity) error {
// 	_import.Trim(data.basePath)

// 	if isSideEffectImport(_import) {
// 		contains, err := data.pkgCache.CheckSideEffectImport([]byte(_import.Path))
// 		if err != nil {
// 			return err
// 		}

// 		if contains {
// 			return nil
// 		}

// 		if err := s.db.PackageAcc().ImportAcc().AppendSideEffectImport(ctx, _import, data.newPkg.Path); err != nil {
// 			return err
// 		}

// 		data.pkgCache.AddSideEffectImport(_import)
// 		return nil
// 	}

// 	contains, err := data.pkgCache.CheckImport([]byte(_import.Path))
// 	if err != nil {
// 		return err
// 	}

// 	if !contains {
// 		if err := s.db.PackageAcc().ImportAcc().Append(ctx, _import, data.newPkg.Path); err != nil {
// 			return err
// 		}

// 		data.pkgCache.AddImport(_import)
// 		return nil
// 	}

// 	if index := data.pkgCache.GetImportIndex(_import.Alias); index < 0 {
// 		for i, imp := range data.pkgCache.GetImports() {
// 			if imp == _import.Path {
// 				data.pkgCache.AddImportIndex(_import, i)
// 			}
// 		}
// 	}

// 	return nil
// }

// func fetchPackageImports(files map[string]*ast.File) ([]*entity.ImportEntity, int) {
// 	var impTotal int
// 	var imports []*entity.ImportEntity

// 	for _, file := range files {
// 		impTotal = impTotal + len(file.Imports)

// 		for _, imp := range file.Imports {
// 			if imp.Path != nil && imp.Path.Value != "" {
// 				imports = append(imports, entity.NewImportEntity(imp))
// 			}
// 		}
// 	}

// 	return imports, impTotal
// }

// func isSideEffectImport(_import *entity.ImportEntity) bool {
// 	return _import.WithAlias && _import.Alias == "_"
// }
