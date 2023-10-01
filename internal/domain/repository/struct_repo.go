package repository

import (
	"context"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type StructRepository interface {
	Append(ctx context.Context, structDTO *domainDTO.StructDTO, structIndex int, pkgPath string) error
}
