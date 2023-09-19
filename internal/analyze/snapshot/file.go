package snapshot

import (
	"errors"
	"fmt"
	"go/ast"
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/analyze/entity"
)

type FileManifest struct {
	StructTypeMap    map[string]*entity.StructType
	InterfaceTypeMap map[string]*entity.InterfaceType
	Imports          map[string]string
	BelongToPackage  string
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

func (fm *FileManifest) AddStructType(structName string, structType *entity.StructType) {
	if fm.StructTypeMap == nil {
		fm.StructTypeMap = make(map[string]*entity.StructType)
	}

	fm.StructTypeMap[structName] = structType
}

func (fm *FileManifest) AddInterfaceType(interfaceName string, it *entity.InterfaceType) {
	if fm.InterfaceTypeMap == nil {
		fm.InterfaceTypeMap = make(map[string]*entity.InterfaceType)
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

func (fm *FileManifest) AddMethodToStruct(structName, methodName, fieldName string, fieldUsage entity.FieldUsage) error {
	if fm.StructTypeMap == nil {
		return errors.New("StructTypeMap is not initialized")
	}

	structType, exists := fm.StructTypeMap[structName]
	if !exists {
		return fmt.Errorf("structure %s does not exist", structName)
	}

	if structType.Methods == nil {
		structType.Methods = make(map[string]map[string]entity.FieldUsage)
	}

	if structType.Methods[methodName] == nil {
		structType.Methods[methodName] = make(map[string]entity.FieldUsage)
	}

	structType.Methods[methodName][fieldName] = fieldUsage
	return nil
}

func NewFileManifest(bToPkg string) *FileManifest {
	return &FileManifest{
		StructTypeMap:    make(map[string]*entity.StructType),
		InterfaceTypeMap: make(map[string]*entity.InterfaceType),
		Imports:          make(map[string]string),
		BelongToPackage:  bToPkg,
	}
}
