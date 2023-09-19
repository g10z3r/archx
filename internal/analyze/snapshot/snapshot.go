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

	for _, structInfo := range fm.Structs {
		packageManifest.Structs = append(packageManifest.Structs, structInfo)
	}
	packageManifest.StructsIndex = fm.StructsIndex

	for _, interfaceInfo := range fm.Interfaces {
		packageManifest.Interfaces = append(packageManifest.Interfaces, interfaceInfo)
	}
	packageManifest.InterfacesIndex = fm.InterfacesIndex

	return nil
}
func NewSnapshot() *SnapshotManifest {
	return &SnapshotManifest{
		PackageMap: make(map[string]*PackageManifest),
	}
}
