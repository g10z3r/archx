package buffer

import (
	"fmt"

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
		errChan <- fmt.Errorf("buffer is not of type *StructBuffer")
		return
	}

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	if existingIndex, exists := buf.StructsIndex[e.StructName]; exists {
		existingStruct := buf.Structs[existingIndex]

		if !existingStruct.Incompplete && !e.StructInfo.Incompplete {
			existingStruct.SyncMethods(e.StructInfo)
			buf.Structs[existingIndex] = existingStruct
		}

		if !existingStruct.Incompplete && e.StructInfo.Incompplete {
			e.StructInfo.SyncMethods(existingStruct)
			buf.Structs[existingIndex] = e.StructInfo
		}
	} else {
		buf.Structs = append(buf.Structs, e.StructInfo)
		index := len(buf.Structs) - 1
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

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	buf.Structs[e.StructIndex].AddMethod(e.Method, e.MethodName)
}
