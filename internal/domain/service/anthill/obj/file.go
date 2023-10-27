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
	Imports       *FileObjImports
	Structs       []*StructObj
	StructIndexes map[string]int
	Functions     []*FuncObj
}

func (obj *FileobjEntities) AppendStruct(o *StructObj) {
	obj.StructIndexes[*o.Name] = len(obj.Structs)
	obj.Structs = append(obj.Structs, o)
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

type FileObj struct {
	Name     string
	FileSet  *token.FileSet
	Entities *FileobjEntities
	Metadata *FileObjMatadata
	Stats    *FileObjStats
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
