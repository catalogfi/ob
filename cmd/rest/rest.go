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
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	PORT           string       `binding:"required"`
	PSQL_DB        string       `binding:"required"`
	PRICE_FEED_URL string       `binding:"required"`
	SERVER_SECRET  string       `binding:"required"`
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
	jsonParser.Decode(&config)
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
	price := price.NewPriceChecker(store, envConfig.PRICE_FEED_URL)
	go price.Run()
	server := rest.NewServer(store, envConfig.CONFIG, logger, "SECRET")
	if err := server.Run(fmt.Sprintf(":%s", envConfig.PORT)); err != nil {
		panic(err)
	}
}
