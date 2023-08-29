package main

import (
	"context"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/rest"
	"github.com/catalogfi/wbtc-garden/screener"
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
			RPC: map[string]string{
				"mempool": "https://mempool.space/testnet/api",
			},
			Oracles: map[model.Asset]model.Token{
				model.NewSecondary(""): {
					TokenAddress: "",
					Decimals:     0,
				},
			},
			Expiry: 144,
		},
		model.EthereumSepolia: {
			RPC: map[string]string{
				"ethrpc": "https://gateway.tenderly.co/public/sepolia",
			},
			Oracles: map[model.Asset]model.Token{
				model.NewSecondary(""): {
					PriceUrl:     "",
					TokenAddress: "",
					Decimals:     0,
				},
			},
			Expiry: 6542,
		},
		model.EthereumOptimism: {
			RPC: map[string]string{
				"ethrpc": "https://opt-mainnet.g.alchemy.com/v2/lM_wORHU7fDVp_SSYJPCCO-erSffgpX9",
			},
			Oracles: map[model.Asset]model.Token{
				model.NewSecondary(""): {
					PriceUrl:     "",
					TokenAddress: "",
					Decimals:     0,
				},
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

	// Screen is not doing sanction check in this case
	screener := screener.NewScreener(nil, "")
	server := rest.NewServer(store, model.Config{Network: config}, logger, "SECRET", screener)
	if err := server.Run(context.Background(), ":8080"); err != nil {
		panic(err)
	}
}
