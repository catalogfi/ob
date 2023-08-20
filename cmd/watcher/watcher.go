package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/catalogfi/wbtc-garden/logger"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/store"
	"github.com/catalogfi/wbtc-garden/watcher"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	SENTRY_DSN string
	PSQL_DB    string       `binding:"required"`
	CONFIG     model.Config `binding:"required"`
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
	env := LoadConfiguration("./config.json")
	store, err := store.New(postgres.Open(env.PSQL_DB), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		panic(err)
	}

	var log *zap.Logger
	if env.SENTRY_DSN != "" {
		log = zap.New(logger.NewSentryCore(env.SENTRY_DSN, zapcore.ErrorLevel))
	} else {
		log, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	}
	defer log.Sync()

	watcher := watcher.NewWatcher(log, store, env.CONFIG)
	watcher.Run()
}
