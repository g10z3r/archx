package analyze

import (
	"go/ast"
	"path"
	"strings"
	"sync"

	"github.com/g10z3r/archx/internal/analyze/entity"
)

type PackageBuffer struct {
	mutex sync.Mutex
	name  string

	Imports map[string]string

	Structs      []*entity.StructInfo
	StructsIndex map[string]int

	Interfaces      []*entity.InterfaceType
	InterfacesIndex map[string]int
}

// TODO - remove ast IMPORT
func (fm *PackageBuffer) AddImport(importSpec *ast.ImportSpec, mod string) {
	// Remove quotes around the imported string
	importPath := strings.Trim(importSpec.Path.Value, `"`)

	if !strings.HasPrefix(importPath, mod) {
		return
	}

	if importSpec.Name != nil {
		fm.Imports[importSpec.Name.Name] = importPath
		return
	}

	fm.Imports[path.Base(importPath)] = importPath
}

// func (pb *PackageBuffer) AddStruct(structType *entity.StructInfo, structName string) {
// 	pb.mutex.Lock()
// 	defer pb.mutex.Unlock()

// 	pb.Structs = append(pb.Structs, structType)
// 	pb.StructsIndex[structName] = len(pb.Structs) - 1
// }

func (pb *PackageBuffer) AddStruct(incomingStructType *entity.StructInfo, incomingStructName string) int {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	if existingIndex, exists := pb.StructsIndex[incomingStructName]; exists {
		// Если структура уже существует в буфере
		existingStruct := pb.Structs[existingIndex]

		if !existingStruct.IsFull && !incomingStructType.IsFull {
			// Просто обновляем методы и заменяем структуру в буфере
			existingStruct.SyncMethods(incomingStructType)
			pb.Structs[existingIndex] = existingStruct
			return existingIndex
		}

		if !existingStruct.IsFull && incomingStructType.IsFull {
			// Синхронизируем методы и заменяем структуру в буфере
			incomingStructType.SyncMethods(existingStruct)
			pb.Structs[existingIndex] = incomingStructType
			return existingIndex
		}
	}

	// Если структуры нет в буфере, просто добавляем её
	pb.Structs = append(pb.Structs, incomingStructType)
	index := len(pb.Structs) - 1
	pb.StructsIndex[incomingStructName] = index

	return index
}

func (pb *PackageBuffer) GetStructByName(name string) (*entity.StructInfo, int, bool) {
	index, exist := pb.StructsIndex[name]
	if !exist {
		return nil, 0, false
	}

	return pb.Structs[index], index, true
}

func (fm *PackageBuffer) HasStructType(structName string) bool {
	if fm.StructsIndex == nil {
		return false
	}

	_, exists := fm.StructsIndex[structName]
	return exists
}

func (pb *PackageBuffer) AddInterface(interfaceType *entity.InterfaceType, interfaceName string) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.Interfaces = append(pb.Interfaces, interfaceType)
	pb.InterfacesIndex[interfaceName] = len(pb.Interfaces) - 1
}

func NewPackageBuffer(packageName string) *PackageBuffer {
	return &PackageBuffer{
		mutex:           sync.Mutex{},
		name:            packageName,
		Imports:         make(map[string]string),
		Structs:         make([]*entity.StructInfo, 0),
		StructsIndex:    make(map[string]int),
		Interfaces:      make([]*entity.InterfaceType, 0),
		InterfacesIndex: make(map[string]int),
	}
}
