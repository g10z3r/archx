package repository

import (
	"context"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type PackageRepository interface {
	Append(ctx context.Context, pkg *domainDTO.PackageDTO, packageIndex int) error

	StructRepo() StructRepository
	ImportRepo() ImportRepository
}
