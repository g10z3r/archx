package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"

	"github.com/g10z3r/archx/pkg/bloom"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type packageCache struct {
	mu sync.RWMutex

	importsFilter     bloom.BloomFilter
	sideEffectImports bloom.BloomFilter

	Imports      []string
	ImportsIndex map[string]map[string]int

	StructsIndex map[string]int
}

func (p *packageCache) GetStructsIndexLength() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.StructsIndex)
}

func (pc *packageCache) GetImports() []string {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	return pc.Imports
}

func (pc *packageCache) AddStructIndex(structName string) int {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	index := len(pc.StructsIndex)
	pc.StructsIndex[structName] = index
	fmt.Println("Saved to cache", structName, index)
	return index
}

func (pc *packageCache) GetStructIndex(structName string) int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	index, exists := pc.StructsIndex[structName]
	if !exists {
		return -1
	}

	return index
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

func (pc *packageCache) AddSideEffectImport(_import *entity.ImportEntity) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.sideEffectImports.Put([]byte(_import.Path))
}

func (pc *packageCache) AddImport(_import *entity.ImportEntity) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.importsFilter.Put([]byte(_import.Path))

	pc.Imports = append(pc.Imports, _import.Path)

	if _, exists := pc.ImportsIndex[_import.File]; !exists {
		pc.ImportsIndex[_import.File] = make(map[string]int)
	}

	pc.ImportsIndex[_import.File][getAlias(_import)] = len(pc.Imports) - 1
}

func (pc *packageCache) AddImportAlias(_import *entity.ImportEntity, index int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.ImportsIndex[_import.File][getAlias(_import)] = index
}

func getAlias(_import *entity.ImportEntity) string {
	if _import.WithAlias {
		return _import.Alias
	}

	return path.Base(_import.Path)
}

func (pc *packageCache) GetImportIndex(fileName, alias string) int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	index, exists := pc.ImportsIndex[fileName][alias]
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
		ImportsIndex:      make(map[string]map[string]int),
		StructsIndex:      make(map[string]int),
	}
}

func (pc *packageCache) Debug() {
	jsonData, err := json.Marshal(pc)
	if err != nil {
		log.Println("Error marshalling data:", err)
		return
	}

	log.Println(string(jsonData))
}
