package repository

import (
	"context"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type ScannerRepository interface {
	Create(ctx context.Context, result *domainDTO.ScannerResultDTO) error
	PackageRepo() PackageRepository
}
