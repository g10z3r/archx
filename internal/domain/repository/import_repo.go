package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type ImportRepository interface {
	Append(ctx context.Context, _import *entity.ImportEntity, packagePath string) error
	AppendSideEffectImport(ctx context.Context, _import *entity.ImportEntity, packagePath string) error
}
