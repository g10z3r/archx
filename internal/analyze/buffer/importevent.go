package buffer

import (
	"fmt"
	"path"
	"strings"

	"github.com/g10z3r/archx/internal/analyze/entity"
)

type AddImportEvent struct {
	Import *entity.Import
	Mod    string
}

func (e *AddImportEvent) ToBuffer() int {
	return toImportsBuffer
}

func (e *AddImportEvent) Execute(buffer bufferBus, errChan chan<- error) {
	buf, ok := buffer.(*ImportBuffer)
	if !ok {
		errChan <- fmt.Errorf("buffer is not of type *StructBuffer")
		return
	}

	buf.mutex.Lock()
	defer buf.mutex.Unlock()

	if !strings.HasPrefix(e.Import.Path, e.Mod) {
		return
	}

	if e.Import.WithAlias {
		buf.imports[e.Import.Alias] = e.Import.Path
		return
	}

	buf.imports[path.Base(e.Import.Path)] = e.Import.Path

}
