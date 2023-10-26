package obj

import (
	"sync"
)

type PackageObj struct {
	lock  sync.Mutex
	Name  string
	Path  string
	Files []*FileObj
}

func NewPackageObj(path, name string) *PackageObj {
	return &PackageObj{Path: path,
		Name: name}
}
func (obj *PackageObj) AppendFile(f *FileObj) {
	obj.lock.Lock()
	obj.Files = append(obj.Files, f)
	obj.lock.Unlock()
}
