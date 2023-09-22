package buffer

import (
	"fmt"

	"github.com/g10z3r/archx/internal/analyze/entity"
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
		errChan <- fmt.Errorf("buffer is not of type *StructBuffer")
		return
	}

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	if existingIndex, exists := buf.structsIndex[e.StructName]; exists {
		existingStruct := buf.structs[existingIndex]

		if !existingStruct.IsFull && !e.StructInfo.IsFull {
			existingStruct.SyncMethods(e.StructInfo)
			buf.structs[existingIndex] = existingStruct
		}

		if !existingStruct.IsFull && e.StructInfo.IsFull {
			e.StructInfo.SyncMethods(existingStruct)
			buf.structs[existingIndex] = e.StructInfo
		}
	} else {
		buf.structs = append(buf.structs, e.StructInfo)
		index := len(buf.structs) - 1
		buf.structsIndex[e.StructName] = index
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

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	buf.structs[e.StructIndex].AddMethod(e.Method, e.MethodName)
}
