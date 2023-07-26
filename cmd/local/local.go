package main

import (
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/price"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"github.com/susruth/wbtc-garden/watcher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// sqlite db
	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	config := model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinTestnet:  "https://mempool.space/testnet/api",
			model.EthereumSepolia: "http://localhost:8545",
		},
	}

	watcher := watcher.NewWatcher(store, config)
	price := price.NewPriceChecker(store, "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd")
	go price.Run()
	go watcher.Run()
	server := rest.NewServer(store, config, "SECRET")
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
