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
			model.EthereumSepolia: "https://gateway.tenderly.co/public/sepolia",
			model.EthereumOptimism: "https://opt-mainnet.g.alchemy.com/v2/lM_wORHU7fDVp_SSYJPCCO-erSffgpX9",
		},
	}

	watcher := watcher.NewWatcher(store, config)
	price := price.NewPriceChecker(store, "https://api.coincap.io/v2/assets/bitcoin")
	go price.Run()
	go watcher.Run()
	server := rest.NewServer(store, config, "SECRET")
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
