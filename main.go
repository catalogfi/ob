package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/catalogfi/ob/feehub"
	"github.com/catalogfi/ob/internal/path"
	"github.com/catalogfi/ob/model"
	"github.com/catalogfi/ob/price"
	"github.com/catalogfi/ob/rest"
	"github.com/catalogfi/ob/screener"
	"github.com/catalogfi/ob/store"
	watchers "github.com/catalogfi/ob/watcher"
	"github.com/ethereum/go-ethereum/common"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	SENTRY_DSN    string
	PORT          string       `binding:"required"`
	PSQL_DB       string       `binding:"required"`
	SERVER_SECRET string       `binding:"required"`
	CONFIG        model.Config `binding:"required"`
	FEEHUB_URL    string       `binding:"required"`
	PRICE_URL     string       `binding:"required"`
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
	// psql db
	envConfig := LoadConfiguration(path.ConfigPath)
	store, err := store.New(postgres.Open(envConfig.PSQL_DB), path.SQLSetupPath, &gorm.Config{
		NowFunc: func() time.Time { return time.Now().UTC() },
		Logger:  logger.Default.LogMode(logger.Silent),
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

	watcher := watchers.NewWatcher(logger, store, 4)
	go watcher.Run(context.Background())

	screener := screener.NewScreener(store.Gorm(), envConfig.TRM_KEY)
	for chain, Network := range envConfig.CONFIG.Network {
		if chain.IsBTC() {
			//interval is set to 10 seconds to detect iw tx's quicky
			btcWatcher := watchers.NewBTCWatcher(store, chain, envConfig.CONFIG, screener, 5*time.Second, logger)
			go btcWatcher.Watch(context.Background())
		} else if chain.IsEVM() {
			if chain == model.EthereumArbitrum {
				for asset, token := range Network.Assets {
					ethl2Watcher, err := watchers.NewEthereumL2Watcher(store, chain, Network, common.HexToAddress(string(asset)), token.StartBlock, uint64(Network.EventWindow), screener, logger)
					if err != nil {
						panic(err)
					}
					go ethl2Watcher.Watch()
				}
			} else {
				for asset, token := range Network.Assets {
					ethWatcher, err := watchers.NewEthereumWatcher(store, chain, Network, common.HexToAddress(string(asset)), token.StartBlock, uint64(Network.EventWindow), screener, logger)
					if err != nil {
						panic(err)
					}
					go ethWatcher.Watch()
				}
			}
		}

	}
	socketPool := rest.NewSocketPool()
	listener := rest.NewDBListener(envConfig.PSQL_DB, socketPool, logger, store)
	go listener.Start("updates_to_orders", "updates_to_atomic_swaps", "added_to_orders")

	priceFetcher := price.NewPriceFetcher(price.Options{
		URL: envConfig.PRICE_URL,
	})
	feehubClient := feehub.NewFeehubClient(envConfig.FEEHUB_URL)
	server := rest.NewServer(store, envConfig.CONFIG, logger, envConfig.SERVER_SECRET, socketPool, screener, feehubClient, priceFetcher)
	if err := server.Run(context.Background(), fmt.Sprintf(":%s", envConfig.PORT)); err != nil {
		panic(err)
	}
}
