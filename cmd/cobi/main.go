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
		DEPLOYERS: map[model.Chain]string{
			model.EthereumLocalnet: "0x2D914D5a5F9d66c16475B2F5cA9a820308F8B35a",
		},
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
