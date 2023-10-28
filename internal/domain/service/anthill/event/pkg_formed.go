package event

import "github.com/g10z3r/archx/internal/domain/service/anthill/analyzer/obj"

type PackageFormedEvent struct {
	Package *obj.PackageObj
}

func (d *PackageFormedEvent) Name() string {
	return "pkgFormedEvent"
}
