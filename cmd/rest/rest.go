package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/catalogfi/orderbook/internal/path"
	"github.com/catalogfi/orderbook/model"
	"github.com/catalogfi/orderbook/rest"
	"github.com/catalogfi/orderbook/screener"
	"github.com/catalogfi/orderbook/store"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	PORT          string       `binding:"required"`
	PSQL_DB       string       `binding:"required"`
	SERVER_SECRET string       `binding:"required"`
	CONFIG        model.Config `binding:"required"`
	TRM_KEY       string
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
	envConfig := LoadConfiguration(path.ConfigPath)
	store, err := store.New(postgres.Open(envConfig.PSQL_DB), path.SQLSetupPath, &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	screener := screener.NewScreener(store.Gorm(), envConfig.TRM_KEY)
	socketPool := rest.NewSocketPool()
	listener := rest.NewDBListener(envConfig.PSQL_DB, socketPool, logger, store)
	go listener.Start("updates_to_orders", "updates_to_atomic_swaps", "added_to_orders")

	server := rest.NewServer(store, envConfig.CONFIG, logger, envConfig.SERVER_SECRET, socketPool, screener)
	if err := server.Run(context.Background(), fmt.Sprintf(":%s", envConfig.PORT)); err != nil {
		panic(err)
	}
}
