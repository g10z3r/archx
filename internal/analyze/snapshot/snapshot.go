package snapshot

type SnapshotManifest struct {
	Package map[string]*PackageManifest
}

func (sm *SnapshotManifest) UpdateFromFileManifest(fm *FileManifest) error {
	if sm.Package == nil {
		sm.Package = make(map[string]*PackageManifest)
	}

	packageManifest, exists := sm.Package[fm.BelongToPackage]
	if !exists {
		packageManifest = NewPackageManifest()
		sm.Package[fm.BelongToPackage] = packageManifest
	}

	for k, v := range fm.StructTypeMap {
		packageManifest.StructTypeMap[k] = v
	}

	for k, v := range fm.InterfaceTypeMap {
		packageManifest.InterfaceTypeMap[k] = v
	}

	return nil
}

func NewSnapshot() *SnapshotManifest {
	return &SnapshotManifest{
		Package: make(map[string]*PackageManifest),
	}
}
