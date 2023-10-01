package repository

import (
	"context"

	domainDTO "github.com/g10z3r/archx/internal/domain/dto"
)

type ImportRepository interface {
	Append(ctx context.Context, _import *domainDTO.ImportDTO, packagePath string) error
	AppendSideEffectImport(ctx context.Context, _import *domainDTO.ImportDTO, packagePath string) error
}
