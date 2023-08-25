package main

import (
	"context"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/rest"
	"github.com/catalogfi/wbtc-garden/store"
	"github.com/catalogfi/wbtc-garden/watcher"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// sqlite db
	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	config := model.Network{
		model.BitcoinTestnet: {
			RPC: "https://mempool.space/testnet/api",
			Oracles: map[model.Asset]string{
				model.Primary: "https://api.coincap.io/v2/assets/bitcoin",
			},
			Expiry: 144,
		},
		model.EthereumSepolia: {
			RPC: "https://gateway.tenderly.co/public/sepolia",
			Oracles: map[model.Asset]string{
				model.NewSecondary(""): "https://api.coincap.io/v2/assets/bitcoin",
			},
			Expiry: 6542,
		},
		model.EthereumOptimism: {
			RPC: "https://opt-mainnet.g.alchemy.com/v2/lM_wORHU7fDVp_SSYJPCCO-erSffgpX9",
			Oracles: map[model.Asset]string{
				model.NewSecondary(""): "https://api.coincap.io/v2/assets/bitcoin",
			},
			Expiry: 10000,
		},
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	watcher := watcher.NewWatcher(logger, store, config, 4)
	go watcher.Run(context.Background())
	server := rest.NewServer(store, model.Config{Network: config}, logger, "SECRET")
	if err := server.Run(context.Background(), ":8080"); err != nil {
		panic(err)
	}
}
