package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/cobi"
	"github.com/susruth/wbtc-garden/model"
)

type Config struct {
	PORT    string
	PSQL_DB string
	BTC_RPC string
	ETH_RPC string
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
	envConfig := LoadConfiguration("./config.json")
	if err := cobi.Run(model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinTestnet:  envConfig.BTC_RPC,
			model.EthereumSepolia: envConfig.ETH_RPC,
		},
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
