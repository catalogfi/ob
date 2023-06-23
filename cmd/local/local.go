package main

import (
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"github.com/susruth/wbtc-garden/watcher"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// sqlite db
	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	watcher := watcher.NewWatcher(store)
	go watcher.Run()
	server := rest.NewServer(store)
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
