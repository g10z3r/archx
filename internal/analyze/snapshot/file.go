package snapshot

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/analyze/entity"
)

type FileManifest struct {
	Structs      []*entity.StructInfo
	StructsIndex map[string]int

	Interfaces      []*entity.InterfaceType
	InterfacesIndex map[string]int

	Imports         map[string]string
	BelongToPackage string
}

func (fm *FileManifest) AddImport(t *ast.ImportSpec, mod string) {
	// Remove quotes around the imported string
	importPath := t.Path.Value[1 : len(t.Path.Value)-1]

	if !strings.HasPrefix(importPath, mod) {
		return
	}

	if t.Name != nil {
		fm.Imports[t.Name.Name] = importPath
		return
	}

	fm.Imports[path.Base(importPath)] = importPath
}

func (fm *FileManifest) AddStruct(structName string, structType *entity.StructInfo) {
	if fm.Structs == nil {
		fm.Structs = make([]*entity.StructInfo, 0)
	}
	if fm.StructsIndex == nil {
		fm.StructsIndex = make(map[string]int)
	}

	fm.Structs = append(fm.Structs, structType)
	fm.StructsIndex[structName] = len(fm.Structs) - 1
}

func (fm *FileManifest) AddInterface(interfaceName string, it *entity.InterfaceType) {
	if fm.Interfaces == nil {
		fm.Interfaces = make([]*entity.InterfaceType, 0)
	}
	if fm.InterfacesIndex == nil {
		fm.InterfacesIndex = make(map[string]int)
	}

	fm.Interfaces = append(fm.Interfaces, it)
	fm.InterfacesIndex[interfaceName] = len(fm.Interfaces) - 1
}

func (fm *FileManifest) HasStructType(structName string) bool {
	if fm.StructsIndex == nil {
		return false
	}
	_, exists := fm.StructsIndex[structName]
	return exists
}

func (fm *FileManifest) IsFieldPresent(structName, fieldName string) (bool, error) {
	if fm.StructsIndex == nil {
		return false, errors.New("structs index is not initialized")
	}

	structIndex, exists := fm.StructsIndex[structName]
	if !exists {
		return false, fmt.Errorf("structure %s does not exist", structName)
	}

	structType := fm.Structs[structIndex]

	if structType.FieldsIndex == nil {
		return false, errors.New("field index is not initialized for the structure")
	}

	_, exists = structType.FieldsIndex[fieldName]
	return exists, nil
}

func (fm *FileManifest) AddMethodToStruct(structName, methodName, fieldName string, fieldUsage entity.Usage) error {
	if fm.StructsIndex == nil {
		return errors.New("structs index is not initialized")
	}

	structIndex, exists := fm.StructsIndex[structName]
	if !exists {
		return fmt.Errorf("structure %s does not exist", structName)
	}

	structInfo := fm.Structs[structIndex]

	var methodInfo *entity.MethodInfo
	methodIndex, exists := structInfo.MethodsIndex[methodName]
	if exists {
		methodInfo = structInfo.Methods[methodIndex]
	} else {
		methodInfo = &entity.MethodInfo{
			Pos:      token.NoPos, // TODO: Set correct position
			End:      token.NoPos, // TODO: Set correct end
			Usages:   make(map[string]entity.Usage),
			IsPublic: false, // TODO: Set correct visibility
		}

		// Add new information about the method to the slice and index
		structInfo.Methods = append(structInfo.Methods, methodInfo)
		structInfo.MethodsIndex[methodName] = len(structInfo.Methods) - 1
	}

	// Update field usage information for this method
	methodInfo.Usages[fieldName] = fieldUsage

	return nil
}

func NewFileManifest(bToPkg string) *FileManifest {
	return &FileManifest{
		Structs:         []*entity.StructInfo{},
		StructsIndex:    make(map[string]int),
		Interfaces:      []*entity.InterfaceType{},
		InterfacesIndex: make(map[string]int),
		Imports:         make(map[string]string),
		BelongToPackage: bToPkg,
	}
}
