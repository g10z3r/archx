package snapshot

import (
	"errors"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze/types"
)

type FileManifest struct {
	StructTypeMap   map[string]*types.StructType
	BelongToPackage string
}

func (fm *FileManifest) AddStructType(structName string, structType *types.StructType) {
	if fm.StructTypeMap == nil {
		fm.StructTypeMap = make(map[string]*types.StructType)
	}
	fm.StructTypeMap[structName] = structType
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

	if structType.Field == nil {
		return false, errors.New("field map is not initialized for the structure")
	}

	_, exists = structType.Field[fieldName]
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

	if structType.Method == nil {
		structType.Method = make(map[string]map[string]struct{})
	}

	if structType.Method[methodName] == nil {
		structType.Method[methodName] = make(map[string]struct{})
	}

	structType.Method[methodName][fieldName] = struct{}{}
	return nil
}

func NewFileManifest(bToPkg string) *FileManifest {
	return &FileManifest{
		StructTypeMap:   make(map[string]*types.StructType),
		BelongToPackage: bToPkg,
	}
}
