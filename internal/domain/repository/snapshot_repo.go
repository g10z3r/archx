package repository

import (
	"context"

	"github.com/g10z3r/archx/internal/domain/service/anthill/obj"
)

type SnapshotRepository interface {
	// Register(ctx context.Context, result *entity.SnapshotEntity) error
	PackageAcc() PackageAccessor
}

type PackageAccessor interface {
	Append(ctx context.Context, pkg *obj.PackageObj, packageIndex int) error

	StructAcc() StructAccessor
	ImportAcc() ImportAccessor
}

type ImportAccessor interface {
	Append(ctx context.Context, _import *obj.ImportObj, packagePath string) error
	AppendSideEffectImport(ctx context.Context, _import *obj.ImportObj, packagePath string) error
}

type StructAccessor interface {
	Append(ctx context.Context, structEntity *obj.StructObj, structIndex int, pkgPath string) error
	// AddMethod(ctx context.Context, methodEntity *entity.MethodEntity, structIndex int, pkgPath string) error
}
