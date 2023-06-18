package main

import (
	"fmt"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/susruth/wbtc-garden/executor"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// psql db
	store, err := store.New(postgres.Open(os.Getenv("PSQL_DB")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	testnet := TestnetConfig()

	testnetSwapper, err := executor.New(testnet, store.SubStore(testnet.Name))
	if err != nil {
		panic(err)
	}
	go testnetSwapper.Run()
	server := rest.NewServer(map[string]rest.Swapper{
		testnet.Name: testnetSwapper,
	})

	if err := server.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}

func TestnetConfig() executor.Config {
	return executor.Config{
		Name:            "testnet",
		Params:          &chaincfg.TestNet3Params,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		WBTCAddress:     os.Getenv("SEPOLIA_WBTC"),
		BitcoinURL:      os.Getenv("BTC_TESTNET_RPC"),
		EthereumURL:     os.Getenv("SEPOLIA_RPC"),
		DeployerAddress: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
	}
}

func EthereumConfig() executor.Config {
	return executor.Config{
		Name:            "ethereum",
		Params:          &chaincfg.MainNetParams,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		WBTCAddress:     os.Getenv("ETHEREUM_WBTC"),
		BitcoinURL:      os.Getenv("BTC_RPC"),
		EthereumURL:     os.Getenv("ETHEREUM_RPC"),
		DeployerAddress: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
	}
}

func OptimismConfig() executor.Config {
	return executor.Config{
		Name:            "optimism",
		Params:          &chaincfg.MainNetParams,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		WBTCAddress:     os.Getenv("OPTIMISM_WBTC"),
		BitcoinURL:      os.Getenv("BTC_RPC"),
		EthereumURL:     os.Getenv("OPTIMISM_RPC"),
		DeployerAddress: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
	}
}

func ArbitrumConfig() executor.Config {
	return executor.Config{
		Name:            "arbitrum",
		Params:          &chaincfg.MainNetParams,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		WBTCAddress:     os.Getenv("ARBITRUM_WBTC"),
		BitcoinURL:      os.Getenv("BTC_RPC"),
		EthereumURL:     os.Getenv("ARBITRUM_RPC"),
		DeployerAddress: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
	}
}
