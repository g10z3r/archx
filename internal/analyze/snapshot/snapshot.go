package snapshot

type SnapshotManifest struct {
	PackageMap map[string]*PackageManifest
}

func (sm *SnapshotManifest) UpdateFromFileManifest(fm *FileManifest) error {
	if sm.PackageMap == nil {
		sm.PackageMap = make(map[string]*PackageManifest)
	}

	packageManifest, exists := sm.PackageMap[fm.BelongToPackage]
	if !exists {
		packageManifest = NewPackageManifest()
		sm.PackageMap[fm.BelongToPackage] = packageManifest
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
		PackageMap: make(map[string]*PackageManifest),
	}
}
