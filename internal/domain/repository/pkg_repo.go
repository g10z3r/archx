package repository

import (
	"context"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type PackageRepository interface {
	Create(ctx context.Context, pkg *domainDTO.PackageDTO, packageIndex int) error
}
