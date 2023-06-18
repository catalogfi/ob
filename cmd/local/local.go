package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/susruth/wbtc-garden/executor"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// sqlite db
	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	confFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("error reading config file (%s): %v", os.Args[1], err))
	}

	config := executor.Config{}
	if err := json.Unmarshal(confFile, &config); err != nil {
		panic(fmt.Sprintf("error parsing config file (%s): %v", os.Args[1], err))
	}
	config.Params = &chaincfg.RegressionNetParams

	swapper, err := executor.New(config, store.SubStore(config.Name))
	if err != nil {
		panic(err)
	}
	go swapper.Run()
	server := rest.NewServer(map[string]rest.Swapper{
		config.Name: swapper,
	})

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
