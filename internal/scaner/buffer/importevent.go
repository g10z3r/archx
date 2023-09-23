package buffer

import (
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/scaner/entity"
)

type AddImportEvent struct {
	Import *entity.Import
}

func (e *AddImportEvent) ToBuffer() int {
	return toImportsBuffer
}

func (e *AddImportEvent) Execute(buffer bufferBus, errChan chan<- error) {
	buf, ok := buffer.(*ImportBuffer)
	if !ok {
		errChan <- errIncorrectImportBufferType
		return
	}

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	if !strings.HasPrefix(e.Import.Path, buf.Module) {
		return
	}

	e.Import.Path = trimBasePath(e.Import.Path, buf.Module)

	contains, err := buf.filter.MightContain([]byte(e.Import.Path))
	if err != nil {
		errChan <- err
	} else if !contains {
		buf.addImport(e.Import)
	} else {
		buf.updateIndexOrSideEffects(e.Import)
	}
}

func (buf *ImportBuffer) addImport(imp *entity.Import) {
	buf.filter.Put([]byte(imp.Path))
	buf.Imports = append(buf.Imports, imp.Path)
	index := len(buf.Imports) - 1

	alias := getAlias(imp)
	if isSideEffectImport(imp) {
		buf.SideEffectImports = append(buf.SideEffectImports, index)
	} else {
		buf.ImportsIndex[alias] = index
	}
}

func (buf *ImportBuffer) updateIndexOrSideEffects(imp *entity.Import) {
	targetMap := buf.ImportsIndex
	if isSideEffectImport(imp) {
		targetMap = make(map[string]int)
	}

	alias := getAlias(imp)
	for index, path := range buf.Imports {
		if path == imp.Path {
			targetMap[alias] = index
			if isSideEffectImport(imp) {
				buf.SideEffectImports = append(buf.SideEffectImports, index)
			}
			return
		}
	}
}

func isSideEffectImport(imp *entity.Import) bool {
	return imp.WithAlias && imp.Alias == "_"
}

func getAlias(imp *entity.Import) string {
	if imp.WithAlias {
		return imp.Alias
	}

	return path.Base(imp.Path)
}

func trimBasePath(fullPath, basePath string) string {
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	return "/" + strings.TrimPrefix(fullPath, basePath)
}
