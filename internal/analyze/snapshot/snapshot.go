package snapshot

type SnapshotManifest struct {
	PackageMap map[string]*PackageManifest
}

func (sm *SnapshotManifest) UpdateFromFileManifest(fm *FileManifest) error {
	// Initialize the package map if it doesn't exist
	if sm.PackageMap == nil {
		sm.PackageMap = make(map[string]*PackageManifest)
	}

	// Get the existing package manifest or create a new one if it doesn't exist
	packageManifest, exists := sm.PackageMap[fm.BelongToPackage]
	if !exists {
		packageManifest = NewPackageManifest()
		sm.PackageMap[fm.BelongToPackage] = packageManifest
	}

	// Get the next index for struct and interface to avoid index conflict
	nextStructIndex := len(packageManifest.Structs)
	nextInterfaceIndex := len(packageManifest.Interfaces)

	// Add new struct info from file manifest to package manifest
	for _, structInfo := range fm.Structs {
		packageManifest.Structs = append(packageManifest.Structs, structInfo)
	}

	// Update the structs index map with new indexes from the file manifest
	for structName, structIndex := range fm.StructsIndex {
		// If structName does not exist in the package manifest, add it with an updated index
		if _, exists := packageManifest.StructsIndex[structName]; !exists {
			packageManifest.StructsIndex[structName] = structIndex + nextStructIndex
		}
	}

	// Add new interface info from file manifest to package manifest
	for _, interfaceInfo := range fm.Interfaces {
		packageManifest.Interfaces = append(packageManifest.Interfaces, interfaceInfo)
	}

	// Update the interfaces index map with new indexes from the file manifest
	for interfaceName, interfaceIndex := range fm.InterfacesIndex {
		// If interfaceName does not exist in the package manifest, add it with an updated index
		if _, exists := packageManifest.InterfacesIndex[interfaceName]; !exists {
			packageManifest.InterfacesIndex[interfaceName] = interfaceIndex + nextInterfaceIndex
		}
	}

	return nil
}

func NewSnapshot() *SnapshotManifest {
	return &SnapshotManifest{
		PackageMap: make(map[string]*PackageManifest),
	}
}
