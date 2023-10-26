package event

import "github.com/g10z3r/archx/internal/domain/service/anthill/obj"

type PackageFormedEvent struct {
	Package *obj.PackageObj
}

func (d *PackageFormedEvent) Name() string {
	return "pkgFormedEvent"
}
