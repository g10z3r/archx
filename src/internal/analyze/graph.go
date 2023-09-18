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

	// Шаг 1: Создайте узлы для каждой структуры в каждом пакете
	for pkgName, pkgManifest := range manifests {
		for structName, _ := range pkgManifest.StructTypeMap {
			nodeName := pkgName + "." + structName
			graph.Nodes[nodeName] = &GraphNode{Name: nodeName}
		}
	}

	// Шаг 2: Создайте зависимости между узлами
	for pkgName, pkgManifest := range manifests {
		for structName, structType := range pkgManifest.StructTypeMap {
			nodeName := pkgName + "." + structName
			node := graph.Nodes[nodeName]

			for depPkg, depTypes := range structType.Dependencies {
				for _, depType := range depTypes {
					depNodeName := depPkg + "." + depType
					depNode, exists := graph.Nodes[depNodeName]

					// Создайте узел для внешней зависимости, если он еще не существует
					if !exists {
						depNode = &GraphNode{Name: depNodeName}
						graph.Nodes[depNodeName] = depNode
					}

					node.Dependencies = append(node.Dependencies, depNode)
				}
			}
		}
	}

	// Шаг 3: Создайте обратные ссылки для афферентной связанности
	for _, node := range graph.Nodes {
		for _, dep := range node.Dependencies {
			dep.ReferencedBy = append(dep.ReferencedBy, node)
		}
	}

	return graph
}
