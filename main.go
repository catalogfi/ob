package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/price"
	"github.com/catalogfi/wbtc-garden/rest"
	"github.com/catalogfi/wbtc-garden/store"
	"github.com/catalogfi/wbtc-garden/watcher"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	PORT           string       `binding:"required"`
	PSQL_DB        string       `binding:"required"`
	PRICE_FEED_URL string       `binding:"required"`
	CONFIG         model.Config `binding:"required"`
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func main() {
	// psql db
	envConfig := LoadConfiguration("./config.json")
	// fmt.Println(envConfig.PSQL_DB)
	store, err := store.New(postgres.Open(envConfig.PSQL_DB), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	watcher := watcher.NewWatcher(logger, store, envConfig.CONFIG)
	price := price.NewPriceChecker(store, envConfig.PRICE_FEED_URL)
	go price.Run()
	go watcher.Run()
	server := rest.NewServer(store, envConfig.CONFIG, logger, "SECRET")
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
