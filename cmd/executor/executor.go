package main

import (
	"os"

	"github.com/susruth/wbtc-garden/bot"
	"github.com/susruth/wbtc-garden/rest"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	store, err := bot.NewStore(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	bot.NewExecutor(os.Args[1], store, rest.NewClient("http://localhost:8080")).Run()
}
