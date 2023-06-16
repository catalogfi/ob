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

	// Local ENV
	// os.Setenv("BITCOIN_URL", "http://localhost:30000")
	// os.Setenv("ETHEREUM_URL", "http://localhost:8545")
	// os.Setenv("WBTC_ADDRESS", "0x85495222Fd7069B987Ca38C2142732EbBFb7175D")
	// os.Setenv("PRIVATE_KEY", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	// os.Setenv("IS_MAINNET", "false")

	// sqlite db
	// store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }

	config := executor.Config{}
	config.BitcoinURL = os.Getenv("BITCOIN_URL")
	config.EthereumURL = os.Getenv("ETHEREUM_URL")
	config.WBTCAddress = os.Getenv("WBTC_ADDRESS")
	privKey := os.Getenv("PRIVATE_KEY")

	switch os.Getenv("NETWORK") {
	case "mainnet":
		config.Network = &chaincfg.MainNetParams
	case "testnet":
		config.Network = &chaincfg.TestNet3Params
	case "regtest":
		config.Network = &chaincfg.RegressionNetParams
	default:
		panic("invalid network")
	}

	swapper, err := executor.New(privKey, config, store)
	if err != nil {
		panic(err)
	}
	go swapper.Run()
	server := rest.NewServer(store, swapper)
	if err := server.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		panic(err)
	}
}