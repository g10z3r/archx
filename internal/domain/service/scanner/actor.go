package scanner

import (
	"context"
	"go/ast"
	"go/token"
	"log"
	"path"
	"sync"

	"github.com/g10z3r/archx/internal/domain/entity"
	"github.com/g10z3r/archx/internal/domain/repository"
	"github.com/g10z3r/archx/internal/domain/service/scanner/cache"
	"github.com/g10z3r/archx/pkg/bloom"
)

type packageCache interface {
	CheckImport(b []byte) (bool, error)
	AddImport(_import *entity.ImportEntity)
	AddImportAlias(_import *entity.ImportEntity, index int)
	GetImportIndex(fileName, alias string) int
	GetImports() []string
	CheckSideEffectImport(b []byte) (bool, error)
	AddSideEffectImport(_import *entity.ImportEntity)

	AddStructIndex(structName string) int
	GetStructIndex(structName string) int
	GetStructsIndexLength() int

	Debug()
}

type task struct {
	key   string
	value func() error
}

type packageActor struct {
	mu sync.RWMutex

	_fset *token.FileSet

	pkg *entity.PackageEntity

	db    repository.PackageAccessor
	cache packageCache
	buf   packageBuffer
}

func (pa *packageActor) ScanFile(ctx context.Context, fileName string, file *ast.File) error {
	log.Printf("Processing file: %s", fileName)

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if err := pa.processGenDecl(ctx, d, fileName); err != nil {
				log.Fatal(err)
			}
		case *ast.FuncDecl:
			if err := pa.processFuncDecl(ctx, d, fileName); err != nil {
				log.Fatal(err)
			}
		}
	}

	pa.cache.Debug()
	return nil
}

func (pa *packageActor) FileSet() *token.FileSet {
	pa.mu.RLock()
	defer pa.mu.RUnlock()

	return pa._fset
}

func (pa *packageActor) IndexationImports(pkg *ast.Package) []*entity.ImportEntity {
	pkgImports, total := fetchPackageImports(pkg.Files)
	pa.cache = cache.NewPackageCache(bloom.FilterConfig{
		ExpectedItemCount:        uint64(total),
		DesiredFalsePositiveRate: 0.01,
	})

	return pkgImports
}

func (pa *packageActor) ProcessAllImports(ctx context.Context, imports []*entity.ImportEntity, modName string) error {
	for _, _import := range imports {
		if err := pa.processImport(ctx, _import, modName); err != nil {
			return err
		}
	}

	return nil
}

type pkgImportData struct {
	pkgCache packageCache
	newPkg   *entity.PackageEntity
	basePath string
}

func (pa *packageActor) processImport(ctx context.Context, _import *entity.ImportEntity, modName string) error {
	_import.Trim(modName)

	if isSideEffectImport(_import) {
		contains, err := pa.cache.CheckSideEffectImport([]byte(_import.Path))
		if err != nil {
			return err
		}

		if contains {
			return nil
		}

		if err := pa.db.ImportAcc().AppendSideEffectImport(ctx, _import, pa.pkg.Path); err != nil {
			return err
		}

		pa.cache.AddSideEffectImport(_import)
		return nil
	}

	contains, err := pa.cache.CheckImport([]byte(_import.Path))
	if err != nil {
		return err
	}

	if !contains {
		//  TODO: save to DB in the end, not one by one
		if err := pa.db.ImportAcc().Append(ctx, _import, pa.pkg.Path); err != nil {
			return err
		}

		pa.cache.AddImport(_import)
		return nil
	}

	if index := pa.cache.GetImportIndex(_import.File, getAlias(_import)); index < 0 {
		for i, imp := range pa.cache.GetImports() {
			if imp == _import.Path {
				pa.cache.AddImportAlias(_import, i)
			}
		}
	}

	return nil
}

func getAlias(_import *entity.ImportEntity) string {
	if _import.WithAlias {
		return _import.Alias
	}

	return path.Base(_import.Path)
}

func fetchPackageImports(files map[string]*ast.File) ([]*entity.ImportEntity, int) {
	var impTotal int
	var imports []*entity.ImportEntity

	for fileName, file := range files {
		impTotal = impTotal + len(file.Imports)

		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				imports = append(imports, entity.NewImportEntity(fileName, imp))
			}
		}
	}

	return imports, impTotal
}

func isSideEffectImport(_import *entity.ImportEntity) bool {
	return _import.WithAlias && _import.Alias == "_"
}

func newPackageActor(dbPkgAcc repository.PackageAccessor, fset *token.FileSet, pkg *entity.PackageEntity) *packageActor {
	return &packageActor{
		mu:    sync.RWMutex{},
		_fset: fset,
		pkg:   pkg,
		db:    dbPkgAcc,
		buf:   NewBuffer(),
	}
}
