package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/cmd/bot"
	"github.com/susruth/wbtc-garden/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	var cmd = &cobra.Command{
		Use: "OrderBook",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
	}

	store, err := store.New(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	cmd.AddCommand(bot.NewCreateCommand(store))
	cmd.AddCommand(bot.NewFillCommand(store))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
