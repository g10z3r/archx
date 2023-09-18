package snapshot

import (
	"errors"
	"fmt"
	"path"

	"github.com/g10z3r/archx/internal/analyze/types"
)

type FileManifest struct {
	StructTypeMap    map[string]*types.StructType
	InterfaceTypeMap map[string]*types.InterfaceType
	Imports          map[string]string
	BelongToPackage  string
}

func (fm *FileManifest) AddImport(fullPath string) {
	if fm.Imports == nil {
		fm.Imports = make(map[string]string)
	}

	fm.Imports[path.Base(fullPath)] = fullPath
}

func (fm *FileManifest) AddStructType(structName string, structType *types.StructType) {
	if fm.StructTypeMap == nil {
		fm.StructTypeMap = make(map[string]*types.StructType)
	}

	fm.StructTypeMap[structName] = structType
}

func (fm *FileManifest) AddInterfaceType(interfaceName string, it *types.InterfaceType) {
	if fm.InterfaceTypeMap == nil {
		fm.InterfaceTypeMap = make(map[string]*types.InterfaceType)
	}

	fm.InterfaceTypeMap[interfaceName] = it
}

func (fm *FileManifest) HasStructType(structName string) bool {
	if fm.StructTypeMap == nil {
		return false
	}
	_, exists := fm.StructTypeMap[structName]
	return exists
}

func (fm *FileManifest) IsFieldPresent(structName, fieldName string) (bool, error) {
	if fm.StructTypeMap == nil {
		return false, errors.New("StructTypeMap is not initialized")
	}

	structType, exists := fm.StructTypeMap[structName]
	if !exists {
		return false, fmt.Errorf("structure %s does not exist", structName)
	}

	if structType.Fields == nil {
		return false, errors.New("field map is not initialized for the structure")
	}

	_, exists = structType.Fields[fieldName]
	return exists, nil
}

func (fm *FileManifest) AddMethodToStruct(structName, methodName, fieldName string) error {
	if fm.StructTypeMap == nil {
		return errors.New("StructTypeMap is not initialized")
	}

	structType, exists := fm.StructTypeMap[structName]
	if !exists {
		return fmt.Errorf("structure %s does not exist", structName)
	}

	if structType.Methods == nil {
		structType.Methods = make(map[string]map[string]struct{})
	}

	if structType.Methods[methodName] == nil {
		structType.Methods[methodName] = make(map[string]struct{})
	}

	structType.Methods[methodName][fieldName] = struct{}{}
	return nil
}

func NewFileManifest(bToPkg string) *FileManifest {
	return &FileManifest{
		StructTypeMap:    make(map[string]*types.StructType),
		InterfaceTypeMap: make(map[string]*types.InterfaceType),
		Imports:          make(map[string]string),
		BelongToPackage:  bToPkg,
	}
}
