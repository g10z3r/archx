package analyze

import "github.com/g10z3r/archx/internal/analyze/snapshot"

type DependencyGraph struct {
	Nodes map[string]*GraphNode
}

type GraphNode struct {
	Name         string
	Dependencies []*GraphNode
	ReferencedBy []*GraphNode
}

func BuildDependencyGraph(manifests map[string]*snapshot.PackageManifest) *DependencyGraph {
	graph := &DependencyGraph{Nodes: make(map[string]*GraphNode)}

	// Creating nodes for each structure in each package
	for pkgName, pkgManifest := range manifests {
		for structName, _ := range pkgManifest.StructTypeMap {
			nodeName := pkgName + "." + structName
			graph.Nodes[nodeName] = &GraphNode{Name: nodeName}
		}
	}

	// Creating dependencies between nodes
	for pkgName, pkgManifest := range manifests {
		for structName, structType := range pkgManifest.StructTypeMap {
			nodeName := pkgName + "." + structName
			node := graph.Nodes[nodeName]

			for depPkg, depTypes := range structType.Dependencies {
				for _, depType := range depTypes {
					depNodeName := depPkg + "." + depType
					depNode, exists := graph.Nodes[depNodeName]

					// Create a node for the external dependency if it doesn't already exist
					if !exists {
						depNode = &GraphNode{Name: depNodeName}
						graph.Nodes[depNodeName] = depNode
					}

					node.Dependencies = append(node.Dependencies, depNode)
				}
			}
		}
	}

	// Create backlinks for afferent connectivity
	for _, node := range graph.Nodes {
		for _, dep := range node.Dependencies {
			dep.ReferencedBy = append(dep.ReferencedBy, node)
		}
	}

	return graph
}
