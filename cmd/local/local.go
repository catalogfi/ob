package main

import (
	"github.com/susruth/wbtc-garden/model"
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
			model.BitcoinRegtest:   "http://localhost:30000",
			model.EthereumLocalnet: "http://localhost:8545",
		},
		DEPLOYERS: map[model.Chain]string{
			model.EthereumLocalnet: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
		},
	}

	watcher := watcher.NewWatcher(store, config)
	go watcher.Run()
	server := rest.NewServer(store, config, "SECRET")
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
