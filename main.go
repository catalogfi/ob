package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/price"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"github.com/susruth/wbtc-garden/watcher"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	PORT            string `binding:"required"`
	PSQL_DB         string `binding:"required"`
	PRICE_FEED_URL  string `binding:"required"`
	BTC_RPC         string
	ETH_RPC         string
	BTC_TESTNET_RPC string
	ETH_SEPOLIA_RPC string
	ETH_OPTIMISM_RPC string
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
func main() {
	// psql db
	envConfig := LoadConfiguration("./config.json")
	// fmt.Println(envConfig.PSQL_DB)
	store, err := store.New(postgres.Open(envConfig.PSQL_DB), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	config := model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinTestnet:  envConfig.BTC_TESTNET_RPC,
			model.EthereumSepolia: envConfig.ETH_SEPOLIA_RPC,
			model.Ethereum:        envConfig.ETH_RPC,
			model.Bitcoin:         envConfig.BTC_RPC,
			model.EthereumOptimism: envConfig.ETH_OPTIMISM_RPC,
		},
	}

	watcher := watcher.NewWatcher(store, config)
	price := price.NewPriceChecker(store, envConfig.PRICE_FEED_URL)
	go price.Run()
	go watcher.Run()
	server := rest.NewServer(store, config, "SECRET")
	if err := server.Run(fmt.Sprintf(":%s", envConfig.PORT)); err != nil {
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
