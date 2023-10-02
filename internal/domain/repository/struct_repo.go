package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type StructRepository interface {
	Append(ctx context.Context, structEntity *entity.StructEntity, structIndex int, pkgPath string) error
}
