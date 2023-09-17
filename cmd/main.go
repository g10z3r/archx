package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze"
)

func main() {
	data, err := analyze.MustParseGoFile("./example/main.go")
	if err != nil {
		fmt.Println("Error analyzing Go file:", err)
		return
	}

	for nodeName, node := range data {
		fmt.Printf("LCOM96 for %s = %f\n", nodeName, analyze.CalculateLCOM96B(node))
		fmt.Printf("LCOM for %s = %f\n", nodeName, analyze.CalculateLCOM(node))
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

}
