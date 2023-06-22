package main

import (
	"encoding/json"
	"os"
)

func main() {
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	values := map[string]string{}
	if err := json.Unmarshal(file, &values); err != nil {
		panic(err)
	}
	// for k, v := range values {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }
}
