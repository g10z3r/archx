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

func (pb *PackageBuffer) AddStruct(structType *entity.StructInfo, structName string) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	pb.Structs = append(pb.Structs, structType)
	pb.StructsIndex[structName] = len(pb.Structs) - 1
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
