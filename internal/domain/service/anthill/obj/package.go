package obj

import (
	"go/ast"
	"sync"
)

type PackageObj struct {
	lock sync.Mutex

	// Name holds the name of the package.
	Name string

	// Path stores the file system path where the package is located.
	Path string

	// Files is an array of FileObj pointers, each representing a file within the package.
	Files []*FileObj

	// FileIndexes maps a file name to its corresponding index in the Files array.
	// The key represents the file name, while the value corresponds to the index of the file in the `Files` array.
	FileIndexes map[string]int
}

func NewPackageObj(pkgAst *ast.Package, path string) *PackageObj {
	return &PackageObj{
		Name:        pkgAst.Name,
		Path:        path,
		FileIndexes: make(map[string]int, len(pkgAst.Files)),
	}
}
func (obj *PackageObj) AppendFile(f *FileObj) {
	obj.lock.Lock()
	obj.FileIndexes[f.Name] = len(obj.Files)
	obj.Files = append(obj.Files, f)
	obj.lock.Unlock()
}
