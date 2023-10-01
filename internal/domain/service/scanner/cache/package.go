package cache

import (
	"path"
	"sync"

	"github.com/g10z3r/archx/pkg/bloom"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type packageCache struct {
	mu sync.RWMutex

	importsFilter     bloom.BloomFilter
	sideEffectImports bloom.BloomFilter

	Imports      []string
	ImportsIndex map[string]int
}

func (pc *packageCache) ImportsLen() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return len(pc.Imports)
}

func (pc *packageCache) CheckImport(b []byte) (bool, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return pc.importsFilter.MightContain(b)
}

func (pc *packageCache) CheckSideEffectImport(b []byte) (bool, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return pc.sideEffectImports.MightContain(b)
}

func (pc *packageCache) AddSideEffectImport(_import *domainDTO.ImportDTO) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.sideEffectImports.Put([]byte(_import.Path))
}

func (pc *packageCache) AddImport(_import *domainDTO.ImportDTO, index int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.importsFilter.Put([]byte(_import.Path))
	pc.Imports = append(pc.Imports, _import.Path)
	pc.ImportsIndex[getAlias(_import)] = index
}

func (pc *packageCache) AddImportIndex(_import *domainDTO.ImportDTO, index int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.ImportsIndex[getAlias(_import)] = index
}

func getAlias(_import *domainDTO.ImportDTO) string {
	if _import.WithAlias {
		return _import.Alias
	}

	return path.Base(_import.Path)
}

func (pc *packageCache) GetImportIndex(importAlias string) int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	index, exists := pc.ImportsIndex[importAlias]
	if !exists {
		return -1
	}

	return index
}

func NewPackageCache(cfg bloom.FilterConfig) *packageCache {
	m, _ := bloom.CalculateFilterParams(
		cfg.ExpectedItemCount,
		float64(cfg.DesiredFalsePositiveRate),
	)

	return &packageCache{
		mu:                sync.RWMutex{},
		importsFilter:     bloom.NewBloomFilter(m),
		sideEffectImports: bloom.NewBloomFilter(m),
		Imports:           make([]string, 0, cfg.ExpectedItemCount),
		ImportsIndex:      make(map[string]int, cfg.ExpectedItemCount),
	}
}