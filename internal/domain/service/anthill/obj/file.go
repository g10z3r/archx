package obj

import (
	"go/token"
)

type FileObjImports struct {
	SideEffectImports  []string
	RegularImports     []string
	RegularImportsMeta map[string]int
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
			SideEffectImports:  make([]string, 0),
			RegularImports:     make([]string, 0),
			RegularImportsMeta: make(map[string]int),
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
