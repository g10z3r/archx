package buffer

import (
	"fmt"
	"log"

	"github.com/g10z3r/archx/internal/scaner/entity"
)

type UpsertStructEvent struct {
	StructInfo *entity.StructInfo
	StructName string
}

func (e *UpsertStructEvent) ToBuffer() int {
	return toStructsBuffer
}

func (e *UpsertStructEvent) Execute(buffer bufferBus, errChan chan<- error) {
	buf, ok := buffer.(*StructBuffer)
	if !ok {
		errChan <- errIncorrectStructBufferType
		return
	}

	buf.mutex.RLock()
	defer buf.mutex.RUnlock()

	e.StructInfo.Mutex.Lock()
	defer e.StructInfo.Mutex.Unlock()

	var index int

	if existingIndex, exists := buf.StructsIndex[e.StructName]; exists {
		existingStruct := buf.Structs[existingIndex]

		if !existingStruct.Incomplete && !e.StructInfo.Incomplete {
			existingStruct.SyncMethods(e.StructInfo)
			existingStruct.SyncDependencies(e.StructInfo)

			log.Printf("Updating struct %s, index %d", e.StructName, existingIndex)
			buf.Structs[existingIndex] = existingStruct
			index = existingIndex
		}

		if !existingStruct.Incomplete && e.StructInfo.Incomplete {
			e.StructInfo.SyncMethods(existingStruct)
			e.StructInfo.SyncDependencies(existingStruct)

			log.Printf("Rewriting struct %s, index %d", e.StructName, existingIndex)
			buf.Structs[existingIndex] = e.StructInfo
			index = existingIndex
		}
	} else {
		buf.Structs = append(buf.Structs, e.StructInfo)
		index = len(buf.Structs) - 1
		log.Printf("Creating struct %s, index %d", e.StructName, index)
		buf.StructsIndex[e.StructName] = index
	}
}

type AddMethodEvent struct {
	StructIndex int
	MethodName  string
	Method      *entity.Method
}

func (e *AddMethodEvent) ToBuffer() int {
	return toStructsBuffer
}

func (e *AddMethodEvent) Execute(buffer bufferBus, errChan chan<- error) {
	buf, ok := buffer.(*StructBuffer)
	if !ok {
		errChan <- fmt.Errorf("buffer is not of type *StructBuffer")
		return
	}

	buf.mutex.RLock()
	defer buf.mutex.RUnlock()

	buf.Structs[e.StructIndex].Mutex.Lock()
	defer buf.Structs[e.StructIndex].Mutex.Unlock()

	log.Printf("Adding method %s to struct %d start method ", e.MethodName, e.StructIndex)
	buf.Structs[e.StructIndex].AddMethod(e.Method, e.MethodName)
	log.Printf("Adding method %s to struct %d end method ", e.MethodName, e.StructIndex)
	log.Printf("Methods len %d", len(buf.Structs[e.StructIndex].Methods))
}
