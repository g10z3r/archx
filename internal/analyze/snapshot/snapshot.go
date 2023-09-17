package snapshot

import (
	"github.com/g10z3r/archx/internal/analyze/types"
)

type SnapshotManifest struct {
	Package map[string]*PackageManifest
}

func (sm *SnapshotManifest) UpdateFromFileManifest(fm *FileManifest) error {
	if sm.Package == nil {
		sm.Package = make(map[string]*PackageManifest)
	}

	packageManifest, exists := sm.Package[fm.BelongToPackage]
	if !exists {
		packageManifest = &PackageManifest{
			StructTypeMap: make(map[string]*types.StructType),
		}
		sm.Package[fm.BelongToPackage] = packageManifest
	}

	for k, v := range fm.StructTypeMap {
		packageManifest.StructTypeMap[k] = v
	}

	return nil
}

// func (s *SnapshotManifest) UpsertPackageManifest(filePath, pkg string) string {
// 	dir, _ := filepath.Split(filePath)
// 	pkgPath := makePkgPath(dir, pkg)
// 	if _, exists := s.Package[pkgPath]; !exists {
// 		s.Package[pkgPath] = NewPackageManifest()
// 	}

// 	return pkgPath
// }

// func (s *SnapshotManifest) AddStructType(packagePath, structTypeName string, structType *types.StructType) {
// 	pkgManifest, exists := s.Package[packagePath]
// 	if !exists {
// 		pkgManifest = NewPackageManifest()
// 		s.Package[packagePath] = pkgManifest
// 	}

// 	pkgManifest.StructTypeMap[structTypeName] = structType
// }

func NewSnapshot() *SnapshotManifest {
	return &SnapshotManifest{
		Package: make(map[string]*PackageManifest),
	}
}

// func makePkgPath(pkgDir, pkgName string) string {
// 	return fmt.Sprintf("%s%s", pkgDir, pkgName)
// }
