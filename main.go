package main

import (
	"os"

	"github.com/susruth/wbtc-garden-server/executor"
	"github.com/susruth/wbtc-garden-server/rest"
	"github.com/susruth/wbtc-garden-server/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// psql db
	// db, err := store.New(postgres.Open(os.Getenv("PSQL_DB")), &gorm.Config{})

	// sqlite db
	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	config := executor.Config{}
	config.BitcoinURL = os.Getenv("BITCOIN_URL")
	config.EthereumURL = os.Getenv("ETHEREUM_URL")
	config.WBTCAddress = os.Getenv("WBTC_ADDRESS")
	if os.Getenv("IS_MAINNET") == "" {
		config.IsMainnet = false
	} else {
		config.IsMainnet = true
	}

	privKey := os.Getenv("PRIVATE_KEY")
	swapper, err := executor.New(privKey, config, store)
	if err != nil {
		panic(err)
	}
	go swapper.Run()
	server := rest.NewServer(store, swapper)
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
