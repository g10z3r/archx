package scanner

import (
	"sync"

	"github.com/g10z3r/archx/pkg/bloom"
)

type scannerCache struct {
	mu sync.RWMutex

	packagesIndex map[string]int
}

func (sc *scannerCache) AddPackage(pkgPath string, index int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.packagesIndex[pkgPath] = index
}

func (sc *scannerCache) GetPackageIndex(pkgName string) int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	index, exists := sc.packagesIndex[pkgName]
	if !exists {
		return -1
	}

	return index
}

func newScannerCache() *scannerCache {
	return &scannerCache{
		mu:            sync.RWMutex{},
		packagesIndex: make(map[string]int),
	}
}

type packageCache struct {
	mu sync.RWMutex

	importsFilter bloom.BloomFilter
	importsIndex  map[string]int
}

func (pc *packageCache) AddImport(importAlias string, index int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.importsIndex[importAlias] = index
}

func (pc *packageCache) GetImportIndex(importAlias string) int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	index, exists := pc.importsIndex[importAlias]
	if !exists {
		return -1
	}

	return index
}

func newPackageCache(mod string, filterCfg bloom.FilterConfig) *packageCache {
	m, _ := bloom.CalculateFilterParams(
		filterCfg.ExpectedItemCount,
		float64(filterCfg.DesiredFalsePositiveRate),
	)

	return &packageCache{
		mu:            sync.RWMutex{},
		importsFilter: bloom.NewBloomFilter(m),
	}
}
