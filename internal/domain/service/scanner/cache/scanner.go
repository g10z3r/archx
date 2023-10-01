package cache

import "sync"

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

func (sc *scannerCache) PackagesIndexLen() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return len(sc.packagesIndex)
}

func NewScannerCache() *scannerCache {
	return &scannerCache{
		mu:            sync.RWMutex{},
		packagesIndex: make(map[string]int),
	}
}
