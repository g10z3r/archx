package obj

import (
	"go/token"
)

type FileObjImports struct {
	InternalImports     []string
	ExternalImports     []string
	SideEffectImports   []string
	InternalImportsMeta map[string]int
}

type FileobjEntities struct {
	Structs   []*StructObj
	Functions []*FuncObj
}
type FileObjStats struct {
	Functions,
	Structs, Interfaces int
}

type FileObjMatadata struct {
	Module string
}

// TODO
// Add MethodIndexes, FuncIndexes, StructIndexes, ... FieldIndexes (?)
// Move Imports to Entities

type FileObj struct {
	Name     string
	FileSet  *token.FileSet
	Imports  *FileObjImports
	Entities *FileobjEntities
	Metadata *FileObjMatadata
	Stats    *FileObjStats
}

func NewFileObj(fset *token.FileSet, moduleName, fileName string) *FileObj {
	return &FileObj{
		Name:    fileName,
		FileSet: fset,
		Imports: &FileObjImports{
			InternalImports:     make([]string, 0),
			ExternalImports:     make([]string, 0),
			SideEffectImports:   make([]string, 0),
			InternalImportsMeta: make(map[string]int),
		},
		Entities: &FileobjEntities{
			Structs:   make([]*StructObj, 0),
			Functions: make([]*FuncObj, 0),
		},
		Metadata: &FileObjMatadata{
			Module: moduleName,
		},
	}
}
