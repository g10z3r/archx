package buffer

import (
	"sync"

	"github.com/g10z3r/archx/internal/scaner/entity"
)

type StructBuffer struct {
	mutex  sync.Mutex
	lenght int
	size   int

	structs      []*entity.StructInfo
	structsIndex map[string]int
}

func (buf *StructBuffer) HandleEvent(event bufferEvent, errChan chan<- error) {
	event.Execute(buf, errChan)
}

func (buf *StructBuffer) Size() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.size
}

func (buf *StructBuffer) Len() int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.lenght
}

func (buf *StructBuffer) IsPresent(key string) bool {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	_, exists := buf.structsIndex[key]
	return exists
}

func (buf *StructBuffer) GetByIndex(index int) *entity.StructInfo {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.structs[index]
}

func (buf *StructBuffer) GetIndex(name string) int {
	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	return buf.structsIndex[name]
}

func newStructBuffer() *StructBuffer {
	return &StructBuffer{
		mutex:        sync.Mutex{},
		structs:      make([]*entity.StructInfo, 0),
		structsIndex: make(map[string]int),
	}
}
