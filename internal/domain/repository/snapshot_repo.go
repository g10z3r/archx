package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type SnapshotRepository interface {
	Register(ctx context.Context, result *entity.SnapshotEntity) error
	PackageRepo() PackageRepository
}
