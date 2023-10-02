package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type PackageRepository interface {
	Append(ctx context.Context, pkg *entity.PackageEntity, packageIndex int) error

	StructRepo() StructRepository
	ImportRepo() ImportRepository
}
