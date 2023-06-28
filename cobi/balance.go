package cobi

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
)

func Balances(entropy []byte, config model.Config) *cobra.Command {
	var (
		chain   string
		user    uint32
		asset   string
		page    int
		perPage int
	)
	cmd := &cobra.Command{
		Use: "balance",
		Run: func(c *cobra.Command, args []string) {
			ch, err := model.ParseChain(chain)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while generating secret: %v", err))
				return
			}

			selectors := []uint32{}
			for i := perPage*page - perPage; i < perPage*page; i++ {
				selectors = append(selectors, uint32(i))
			}

			var units string
			if ch.IsBTC() {
				units = "satoshi"
			} else {
				units = "wei"
			}

			addrs, balances, err := getBalances(entropy, ch, user, selectors, config, model.Asset(asset))
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting addresses: %v", err))
				return
			}
			for i, addr := range addrs {
				fmt.Printf("[%d] %s -> Balance %d %s\n ", selectors[i], addr, balances[i], units)
			}
		},
		DisableAutoGenTag: true,
	}
	cmd.Flags().StringVarP(&chain, "chain", "c", "", "user should provide the chain")
	cmd.MarkFlagRequired("chain")
	cmd.Flags().Uint32Var(&user, "account", 0, "user can provide the user id (default: 0)")
	cmd.Flags().IntVar(&perPage, "per-page", 10, "User can provide number of accounts to display per page (default: 10)")
	cmd.Flags().IntVar(&page, "page", 1, "User can provide which page to display (default: 1)")
	cmd.Flags().StringVarP(&asset, "asset", "a", "primary", "user should provide the asset token address (default: primary)")
	return cmd
}
