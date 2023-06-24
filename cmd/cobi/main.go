package main

import (
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/cobi"
)

func main() {
	if err := cobi.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
