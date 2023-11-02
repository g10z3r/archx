package event

import "github.com/g10z3r/archx/internal/domain/service/anthill/obj"

// e.eventCh <- &event.PackageFormedEvent{
// 	Package: e.parsePkg(fset, pkgAst, targetDir, info.ModuleName),
// }

type PackageFormedEvent struct {
	Package *obj.PackageObj
}

func (d *PackageFormedEvent) Name() string {
	return "PackageFormed"
}
