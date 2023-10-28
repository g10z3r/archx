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

type FileObjMeta struct {
	Module string
}

type FileObj struct {
	mutex sync.Mutex

	Name     string
	FileSet  *token.FileSet
	Entities *FileobjEntities
	Metadata *FileObjMeta
	Stats    *FileObjStats
}

func (o *FileObj) AppendImport(obj *ImportObj) {
	o.mutex.Lock()
	switch obj.ImportType {
	case ImportTypeInternal:
		alias := obj.Alias
		if !obj.WithAlias {
			alias = path.Base(obj.Path)
		}

		o.Entities.Imports.InternalImportsMeta[alias] = len(o.Entities.Imports.InternalImports)
		o.Entities.Imports.InternalImports = append(o.Entities.Imports.InternalImports, obj.Path[len(o.Metadata.Module):])
	case ImportTypeExternal:
		o.Entities.Imports.ExternalImports = append(o.Entities.Imports.ExternalImports, obj.Path)
	case ImportTypeSideEffect:
		o.Entities.Imports.SideEffectImports = append(o.Entities.Imports.SideEffectImports, obj.Path)
	}
	o.mutex.Unlock()
}

func (o *FileObj) AppendStruct(obj *StructObj) {
	o.mutex.Lock()
	o.Entities.StructIndexes[*obj.Name] = len(o.Entities.Structs)
	o.Entities.Structs = append(o.Entities.Structs, obj)
	o.mutex.Unlock()
}

func (o *FileObj) AppendFunc(obj *FuncObj) {
	o.mutex.Lock()

	if obj.Receiver != nil {
		obj.Name = "$" + obj.Name
	}

	o.Entities.FunctionIndexes[obj.Name] = len(o.Entities.Functions)
	o.Entities.Functions = append(o.Entities.Functions, obj)
	o.mutex.Unlock()
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
			Structs:         make([]*StructObj, 0),
			StructIndexes:   make(map[string]int),
			Functions:       make([]*FuncObj, 0),
			FunctionIndexes: make(map[string]int),
		},
		Metadata: &FileObjMeta{
			Module: moduleName,
		},
	}
}
