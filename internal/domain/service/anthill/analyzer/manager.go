package analyzer

import (
	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
	"github.com/g10z3r/archx/pkg/bloom"
)

type ImportManager struct {
	importFilter bloom.BloomFilter
	seiFilter    bloom.BloomFilter
}

func (m *ImportManager) Operate(fileObj *obj.FileObj) {

}
