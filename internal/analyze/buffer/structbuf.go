package buffer

import (
	"sync"

	"github.com/g10z3r/archx/internal/analyze/entity"
)

type StructBuffer struct {
	mutex  sync.Mutex
	lenght int
	size   int

	// TODO: make private
	Structs []*entity.StructInfo
	// TODO: make private
	StructsIndex map[string]int
}

func (buf *StructBuffer) HandleEvent(event Event, errChan chan<- error) {
	event.Execute(buf, errChan)
}

func (buf *StructBuffer) Size() int {
	return buf.size
}

func (buf *StructBuffer) Len() int {
	return buf.lenght
}

func (buf *StructBuffer) IsPresent(key string) bool {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	if buf.StructsIndex == nil {
		return false
	}

	_, exists := buf.StructsIndex[key]
	return exists
}

func newStructBuffer() *StructBuffer {
	return &StructBuffer{
		mutex:        sync.Mutex{},
		Structs:      make([]*entity.StructInfo, 0),
		StructsIndex: make(map[string]int),
	}
}
