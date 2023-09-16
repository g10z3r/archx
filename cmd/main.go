package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/analyze"
)

func main() {
	data, err := analyze.ParseGoFile("./example.go")
	if err != nil {
		fmt.Println("Error analyzing Go file:", err)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

}
