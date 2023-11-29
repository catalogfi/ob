package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/catalogfi/orderbook/internal/path"
	"github.com/catalogfi/orderbook/model"
	"github.com/catalogfi/orderbook/store"
	"github.com/catalogfi/orderbook/watcher"
	watchers "github.com/catalogfi/orderbook/watcher"
	"github.com/catalogfi/orderbook/screener"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	SENTRY_DSN    string
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
	if err := jsonParser.Decode(&config); err != nil {
		panic(err)
	}
	return config
}

func main() {
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

	if envConfig.SENTRY_DSN != "" {
		client, err := sentry.NewClient(sentry.ClientOptions{Dsn: envConfig.SENTRY_DSN})
		if err != nil {
			panic(err)
		}
		cfg := zapsentry.Configuration{
			Level: zapcore.ErrorLevel,
		}
		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))
		if err != nil {
			panic(err)
		}
		logger = zapsentry.AttachCoreToLogger(core, logger)
		defer logger.Sync()
	}

	screener := screener.NewScreener(store.Gorm(), envConfig.TRM_KEY)
	for chain, Network := range envConfig.CONFIG.Network {
		if chain.IsBTC() {
			//interval is set to 10 seconds to detect iw tx's quicky
			btcWatcher := watchers.NewBTCWatcher(store, chain, envConfig.CONFIG, screener, 5*time.Second, logger)
			go btcWatcher.Watch(context.Background())
		} else if chain.IsEVM() {
			for asset, token := range Network.Assets {
				ethWatcher, err := watchers.NewEthereumWatcher(store, chain, Network, common.HexToAddress(string(asset)), token.StartBlock, uint64(Network.EventWindow), screener, logger)
				if err != nil {
					panic(err)
				}
				go ethWatcher.Watch()
			}
		}

	}

	watcher := watcher.NewWatcher(logger, store, 4)
	watcher.Run(context.Background())
}
