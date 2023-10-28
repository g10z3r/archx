package obj

import (
	"go/token"
	"path"
	"sync"
)

type FileObjImports struct {
	InternalImports     []string
	ExternalImports     []string
	SideEffectImports   []string
	InternalImportsMeta map[string]int
}

type FileobjEntities struct {
	Imports         *FileObjImports
	Structs         []*StructObj
	StructIndexes   map[string]int
	Functions       []*FuncObj
	FunctionIndexes map[string]int
}

type FileObjStats struct {
	Functions,
	Structs, Interfaces int
}

type FileObjMatadata struct {
	Module string
}

// TODO
// Add MethodIndexes, ... FieldIndexes (?)

type FileObj struct {
	mutex sync.RWMutex

	Name     string
	FileSet  *token.FileSet
	Entities *FileobjEntities
	Metadata *FileObjMatadata
	Stats    *FileObjStats
}

func (obj *FileObj) AppendImport(o *ImportObj) {
	obj.mutex.Lock()
	switch o.ImportType {
	case ImportTypeInternal:
		obj.Entities.Imports.InternalImportsMeta[getAlias(o)] = len(obj.Entities.Imports.InternalImports)
		obj.Entities.Imports.InternalImports = append(obj.Entities.Imports.InternalImports, o.Path[len(obj.Metadata.Module):])
	case ImportTypeExternal:
		obj.Entities.Imports.ExternalImports = append(obj.Entities.Imports.ExternalImports, o.Path)
	case ImportTypeSideEffect:
		obj.Entities.Imports.SideEffectImports = append(obj.Entities.Imports.SideEffectImports, o.Path)
	}
	obj.mutex.Unlock()
}

func (obj *FileObj) AppendStruct(o *StructObj) {
	obj.mutex.Lock()
	obj.Entities.StructIndexes[*o.Name] = len(obj.Entities.Structs)
	obj.Entities.Structs = append(obj.Entities.Structs, o)
	obj.mutex.Unlock()
}

func (obj *FileObj) AppendFunc(o *FuncObj) {
	obj.mutex.Lock()
	obj.Entities.FunctionIndexes[o.Name] = len(obj.Entities.Functions)
	obj.Entities.Functions = append(obj.Entities.Functions, o)
	obj.mutex.Unlock()
}

func NewFileObj(fset *token.FileSet, moduleName, fileName string) *FileObj {
	return &FileObj{
		Name:    fileName,
		FileSet: fset,
		Entities: &FileobjEntities{
			Imports: &FileObjImports{
				InternalImports:     make([]string, 0),
				ExternalImports:     make([]string, 0),
				SideEffectImports:   make([]string, 0),
				InternalImportsMeta: make(map[string]int),
			},
			Structs:       make([]*StructObj, 0),
			StructIndexes: make(map[string]int),
			Functions:     make([]*FuncObj, 0),
		},
		Metadata: &FileObjMatadata{
			Module: moduleName,
		},
	}
}

func getAlias(importObj *ImportObj) string {
	if importObj.WithAlias {
		return importObj.Alias
	}

	return path.Base(importObj.Path)
}
