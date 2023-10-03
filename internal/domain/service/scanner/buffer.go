package scanner

import (
	"sync"

	"github.com/g10z3r/archx/internal/domain/entity"
)

// Just a tmp solution to implement a main logic
// TODO: implement better buffer
type Buffer struct {
	mu      sync.Mutex
	methods map[string][]*entity.MethodEntity
}

func (mb *Buffer) AddMethod(structName string, method *entity.MethodEntity) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	mb.methods[structName] = append(mb.methods[structName], method)
}

func (mb *Buffer) GetAndClearMethods(structName string) []*entity.MethodEntity {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	methods := mb.methods[structName]
	delete(mb.methods, structName)
	return methods
}

func NewBuffer() *Buffer {
	return &Buffer{
		mu:      sync.Mutex{},
		methods: make(map[string][]*entity.MethodEntity),
	}
}
