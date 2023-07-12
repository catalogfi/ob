package main

import (
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"github.com/susruth/wbtc-garden/watcher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// psql db
	store, err := store.New(sqlite.Open("wbtc_garden.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	config := model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinTestnet:   os.Getenv("BTC_RPC"),
			model.EthereumLocalnet: os.Getenv("ETH_RPC"),
		},
	}

	watcher := watcher.NewWatcher(store, config)
	go watcher.Run()
	server := rest.NewServer(store, config, "SECRET")
	if err := server.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}

// func TestnetConfig() executor.Config {
// 	return executor.Config{
// 		Name:            "testnet",
// 		Params:          &chaincfg.TestNet3Params,
// 		PrivateKey:      os.Getenv("PRIVATE_KEY"),
// 		WBTCAddress:     os.Getenv("SEPOLIA_WBTC"),
// 		BitcoinURL:      os.Getenv("BTC_TESTNET_RPC"),
// 		EthereumURL:     os.Getenv("SEPOLIA_RPC"),
// 		DeployerAddress: "0xf8fC386f964a380007a54D04Ce74E13A2033f26B",
// 	}
// }

// func EthereumConfig() executor.Config {
// 	return executor.Config{
// 		Name:            "ethereum",
// 		Params:          &chaincfg.MainNetParams,
// 		PrivateKey:      os.Getenv("PRIVATE_KEY"),
// 		WBTCAddress:     os.Getenv("ETHEREUM_WBTC"),
// 		BitcoinURL:      os.Getenv("BTC_RPC"),
// 		EthereumURL:     os.Getenv("ETHEREUM_RPC"),
// 		DeployerAddress: "0xf8fC386f964a380007a54D04Ce74E13A2033f26B",
// 	}
// }

// func OptimismConfig() executor.Config {
// 	return executor.Config{
// 		Name:            "optimism",
// 		Params:          &chaincfg.MainNetParams,
// 		PrivateKey:      os.Getenv("PRIVATE_KEY"),
// 		WBTCAddress:     os.Getenv("OPTIMISM_WBTC"),
// 		BitcoinURL:      os.Getenv("BTC_RPC"),
// 		EthereumURL:     os.Getenv("OPTIMISM_RPC"),
// 		DeployerAddress: "0xf8fC386f964a380007a54D04Ce74E13A2033f26B",
// 	}
// }

// func ArbitrumConfig() executor.Config {
// 	return executor.Config{
// 		Name:            "arbitrum",
// 		Params:          &chaincfg.MainNetParams,
// 		PrivateKey:      os.Getenv("PRIVATE_KEY"),
// 		WBTCAddress:     os.Getenv("ARBITRUM_WBTC"),
// 		BitcoinURL:      os.Getenv("BTC_RPC"),
// 		EthereumURL:     os.Getenv("ARBITRUM_RPC"),
// 		DeployerAddress: "0xf8fC386f964a380007a54D04Ce74E13A2033f26B",
// 	}
// }
