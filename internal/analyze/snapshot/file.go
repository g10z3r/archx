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
	StructTypeMap    map[string]*entity.StructInfo
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

func (fm *FileManifest) AddStructType(structName string, structType *entity.StructInfo) {
	if fm.StructTypeMap == nil {
		fm.StructTypeMap = make(map[string]*entity.StructInfo)
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
		return false, errors.New("structTypeMap is not initialized")
	}

	structType, exists := fm.StructTypeMap[structName]
	if !exists {
		return false, fmt.Errorf("structure %s does not exist", structName)
	}

	if structType.FieldsIndex == nil {
		return false, errors.New("field index is not initialized for the structure")
	}

	_, exists = structType.FieldsIndex[fieldName]
	return exists, nil
}

func (fm *FileManifest) AddMethodToStruct(structName, methodName, fieldName string, fieldUsage entity.Usage) error {
	if fm.StructTypeMap == nil {
		return errors.New("structTypeMap is not initialized")
	}

	structInfo, exists := fm.StructTypeMap[structName]
	if !exists {
		return fmt.Errorf("structure %s does not exist", structName)
	}

	methodIndex, exists := structInfo.MethodsIndex[methodName]
	var methodInfo *entity.MethodInfo
	if exists {
		// Если метод существует, получаем его информацию
		methodInfo = structInfo.Methods[methodIndex]
	} else {
		// Если метода нет, создаем новую информацию о методе
		methodInfo = &entity.MethodInfo{
			Pos:      token.NoPos, // TODO: Set correct position
			End:      token.NoPos, // TODO: Set correct end
			Usages:   make(map[string]entity.Usage),
			IsPublic: false, // TODO: Set correct visibility
		}
		// Добавляем новую информацию о методе в слайс и индекс
		structInfo.Methods = append(structInfo.Methods, methodInfo)
		structInfo.MethodsIndex[methodName] = len(structInfo.Methods) - 1
	}

	// Обновляем информацию об использовании поля для этого метода
	methodInfo.Usages[fieldName] = fieldUsage

	return nil
}

func NewFileManifest(bToPkg string) *FileManifest {
	return &FileManifest{
		StructTypeMap:    make(map[string]*entity.StructInfo),
		InterfaceTypeMap: make(map[string]*entity.InterfaceType),
		Imports:          make(map[string]string),
		BelongToPackage:  bToPkg,
	}
}
