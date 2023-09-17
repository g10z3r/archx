package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze"
	"github.com/g10z3r/archx/internal/analyze/snapshot"
)

func main() {
	snapshot := snapshot.NewSnapshot()
	fileManifest, err := analyze.ParseGoFile("./example/main.go")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := snapshot.UpdateFromFileManifest(fileManifest); err != nil {
		fmt.Println(err)
		return
	}

	// for nodeName, node := range data {
	// 	fmt.Printf("LCOM96 for %s = %f\n", nodeName, analyze.CalculateLCOM96B(node))
	// 	fmt.Printf("LCOM for %s = %f\n", nodeName, analyze.CalculateLCOM(node))
	// }

	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

}
