package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/entity"
)

type SnapshotRepository interface {
	Register(ctx context.Context, result *entity.SnapshotEntity) error
	PackageAcc() PackageAccessor
}

type PackageAccessor interface {
	Append(ctx context.Context, pkg *entity.PackageEntity, packageIndex int) error

	StructAcc() StructAccessor
	ImportAcc() ImportAccessor
}

type ImportAccessor interface {
	Append(ctx context.Context, _import *entity.ImportEntity, packagePath string) error
	AppendSideEffectImport(ctx context.Context, _import *entity.ImportEntity, packagePath string) error
}

type StructAccessor interface {
	Append(ctx context.Context, structEntity *entity.StructEntity, structIndex int, pkgPath string) error
}
