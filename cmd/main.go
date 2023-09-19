package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze"
	"github.com/g10z3r/archx/internal/analyze/snapshot"
)

func main() {
	snapshot := snapshot.NewSnapshot()
	fileManifest, err := analyze.ParseGoFile("./example/cmd/main.go", "github.com/g10z3r/archx")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := snapshot.UpdateFromFileManifest(fileManifest); err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range snapshot.PackageMap {
		for sn, si := range p.StructsIndex {
			fmt.Printf("LCOM96 for %s = %f\n", sn, analyze.CalculateLCOM96B(p.Structs[si]))
			fmt.Printf("LCOM for %s = %f\n", sn, analyze.CalculateLCOM(p.Structs[si]))
		}
	}

	jsonData, err := json.Marshal(snapshot)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}
