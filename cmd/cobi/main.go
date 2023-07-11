package main

import (
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/cobi"
	"github.com/susruth/wbtc-garden/model"
)

func main() {
	if err := cobi.Run(model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinRegtest:   "http://localhost:30000",
			model.EthereumLocalnet: "http://localhost:8545",
		},
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
