package buffer

import (
	"sync"

	"github.com/g10z3r/archx/pkg/bloom"
)

type ImportBuffer struct {
	mutex  sync.Mutex
	lenght int
	size   int

	filter bloom.BloomFilter

	Module            string
	Imports           []string
	ImportsIndex      map[string]int
	SideEffectImports []int
}

func (buf *ImportBuffer) HandleEvent(event bufferEvent, errChan chan<- error) {
	event.Execute(buf, errChan)
}

func (buf *ImportBuffer) Size() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.size
}

func (buf *ImportBuffer) Len() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.lenght
}

func (buf *ImportBuffer) IsPresent(key string) (string, bool) {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	index, exists := buf.ImportsIndex[key]
	return buf.Imports[index], exists
}

func newImportBuffer(mod string, filterConfig bloom.FilterConfig) *ImportBuffer {
	m, _ := bloom.CalculateFilterParams(
		filterConfig.ExpectedItemCount,
		float64(filterConfig.DesiredFalsePositiveRate),
	)
	return &ImportBuffer{
		mutex:             sync.Mutex{},
		filter:            bloom.NewBloomFilter(m),
		Module:            mod,
		Imports:           make([]string, 0),
		ImportsIndex:      make(map[string]int),
		SideEffectImports: make([]int, 0),
	}
}
