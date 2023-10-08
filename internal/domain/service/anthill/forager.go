package anthill

import (
	"go/ast"
	"go/token"
	"log"
	"path"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/g10z3r/archx/internal/domain/obj"
	"github.com/g10z3r/archx/pkg/bloom"
	"github.com/g10z3r/archx/pkg/dsl"
)

type bucket[K comparable, V any] struct {
	dsl.Map[K, V]
}

type nest struct {
	importFilter bloom.BloomFilter
	seiFilter    bloom.BloomFilter

	structBucket bucket[string, []*obj.StructObj]
	funcBucket   bucket[string, []*obj.FuncObj]
}

type forager struct {
	storage *nest
	fset    *token.FileSet
	frozen  bool
	count   int64
}

func newForager(fset *token.FileSet) *forager {
	return &forager{
		storage: &nest{},
		fset:    fset,
		frozen:  false,
		count:   0,
	}
}

func (f *forager) process(pkg *ast.Package, pkgPath, modName string) *obj.PackageObj {
	pkgObj := obj.NewPackageObj(pkgPath, pkg.Name)
	head := f.processHead(pkg.Files, modName)

	pkgObj.Imports = append(pkgObj.Imports, head.RegularImports...)
	pkgObj.SideEffectImports = append(pkgObj.SideEffectImports, head.SideEffectImports...)

	f.processBody(pkg.Files, head)
	f.frozen = true

	f.storage.structBucket.Range(func(key string, shard []*obj.StructObj) bool {
		pkgObj.Structs = append(pkgObj.Structs, shard...)
		return true
	})

	f.storage.funcBucket.Range(func(key string, shard []*obj.FuncObj) bool {
		pkgObj.Functions = append(pkgObj.Functions, shard...)
		return true
	})

	return pkgObj
}

type headDTO struct {
	SideEffectImports  []string
	RegularImports     []string
	RegularImportsMeta map[string]map[string]int
}

func makeHeadDTO(total int) *headDTO {
	return &headDTO{
		SideEffectImports:  make([]string, 0),
		RegularImports:     make([]string, 0, total),
		RegularImportsMeta: make(map[string]map[string]int),
	}
}

func (f *forager) processHead(files map[string]*ast.File, modName string) *headDTO {
	imports, total := fetchPackageImports(files)
	f.storage.importFilter = calcAndCreateBloomFilter(total)
	f.storage.seiFilter = calcAndCreateBloomFilter(total / 2)
	dto := makeHeadDTO(len(imports))

	for _, _import := range imports {
		if internal := _import.CheckAndTrim(modName); !internal {
			continue
		}

		if isSideEffectImport(_import) {
			contains, err := f.storage.seiFilter.MightContain([]byte(_import.Path))
			if err != nil {
				continue
			}

			if contains {
				continue
			}

			dto.SideEffectImports = append(dto.SideEffectImports, _import.Path)
			continue
		}

		contains, err := f.storage.importFilter.MightContain([]byte(_import.Path))
		if err != nil {
			log.Fatal(err)
		}

		if _, exists := dto.RegularImportsMeta[_import.File]; !exists {
			dto.RegularImportsMeta[_import.File] = make(map[string]int)
		}

		if !contains {
			if err := f.storage.importFilter.Put([]byte(_import.Path)); err != nil {
				log.Fatal(err)
			}

			dto.RegularImportsMeta[_import.File][getAlias(_import)] = len(dto.RegularImports)
			dto.RegularImports = append(dto.RegularImports, _import.Path)
			continue
		}

		if _, exists := dto.RegularImportsMeta[_import.File][getAlias(_import)]; exists {
			continue
		}

		for i, importPath := range dto.RegularImports {
			if _import.Path == importPath {
				dto.RegularImportsMeta[_import.File][getAlias(_import)] = i
				break
			}
		}
	}

	return dto
}

func (f *forager) processBody(files map[string]*ast.File, head *headDTO) {
	var wg sync.WaitGroup

	for fileName, file := range files {
		atomic.AddInt64(&f.count, 1)
		wg.Add(1)

		fileName = filepath.Base(fileName)

		go func(fset *token.FileSet, file *ast.File, impMeta map[string]int, fileName string) {
			log.Printf("Processing file: %s", fileName)

			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.GenDecl:
					if err := f.processGenDecl(fset, d, impMeta, fileName); err != nil {
						log.Fatal(err)
					}
				case *ast.FuncDecl:
					if err := f.processFuncDecl(fset, d, impMeta, fileName); err != nil {
						log.Fatal(err)
					}
				}
			}

			wg.Done()
		}(f.fset, file, head.RegularImportsMeta[fileName], fileName)
	}

	wg.Wait()
}

func fetchPackageImports(files map[string]*ast.File) ([]*obj.ImportObj, int) {
	var impTotal int
	var imports []*obj.ImportObj

	for fileName, file := range files {
		impTotal = impTotal + len(file.Imports)

		for _, imp := range file.Imports {
			if imp.Path != nil && imp.Path.Value != "" {
				imports = append(imports, obj.NewImportObj(fileName, imp))
			}
		}
	}

	return imports, impTotal
}

func calcAndCreateBloomFilter(total int) bloom.BloomFilter {
	m, _ := bloom.CalcFilterParams(
		uint64(total),
		float64(0.01),
	)

	return bloom.NewBloomFilter(m)
}

func getAlias(_import *obj.ImportObj) string {
	if _import.WithAlias {
		return _import.Alias
	}

	return path.Base(_import.Path)
}

func isSideEffectImport(_import *obj.ImportObj) bool {
	return _import.WithAlias && _import.Alias == "_"
}
