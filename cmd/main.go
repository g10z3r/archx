package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze"
)

func main() {

	// analyze.ParseGoFile2("./example/cmd", "./example/cmd/config.go")
	// snapshot := snapshot.NewSnapshot()
	// files := []string{"./example/cmd/main.go", "./example/cmd/config.go"}

	// for _, v := range files {
	// 	fileManifest, err := analyze.ParseGoFile(v, "github.com/g10z3r/archx")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	if err := snapshot.UpdateFromFileManifest(fileManifest); err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// }

	buf, err := analyze.ParsePackage("./example/cmd", "github.com/g10z3r/archx")
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, p := range snapshot.PackageMap {
	// 	for sn, si := range p.StructsIndex {
	// 		fmt.Printf("LCOM96 for %s = %f\n", sn, analyze.CalculateLCOM96B(p.Structs[si]))
	// 		fmt.Printf("LCOM for %s = %f\n", sn, analyze.CalculateLCOM(p.Structs[si]))
	// 	}
	// }

	jsonData, err := json.Marshal(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))
}
