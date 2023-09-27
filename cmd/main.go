package main

import (
	"encoding/json"
	"fmt"

	"github.com/g10z3r/archx/internal/scaner"
)

func main() {
	total := 100
	var ts string
	var counter int
	for i := 0; i < total; i++ {
		buf, err := scaner.ScanPackage("./example/cmd", "github.com/g10z3r/archx")
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonData, err := json.Marshal(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(jsonData))

		if ts != "" && len(ts) != len(string(jsonData)) {
			fmt.Println(i)
			break
		}

		ts = string(jsonData)
		counter++

	}

	if counter == total {
		fmt.Println("\n\nDone!")
	}

}
