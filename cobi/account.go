package cobi

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
)

func Accounts(entropy []byte) *cobra.Command {
	var (
		chain   string
		user    uint32
		page    int
		perPage int
	)
	cmd := &cobra.Command{
		Use: "accounts",
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

			addrs, err := getAddresses(entropy, ch, user, selectors)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting addresses: %v", err))
				return
			}
			for i, addr := range addrs {
				fmt.Printf("[%d] %s\n", selectors[i], addr)
			}
		},
		DisableAutoGenTag: true,
	}
	cmd.Flags().StringVarP(&chain, "chain", "c", "", "user should provide the chain")
	cmd.MarkFlagRequired("chain")
	cmd.Flags().Uint32Var(&user, "account", 0, "user can provide the user id (default: 0)")
	cmd.Flags().IntVar(&perPage, "per-page", 10, "User can provide number of accounts to display per page (default: 10)")
	cmd.Flags().IntVar(&page, "page", 1, "User can provide which page to display (default: 1)")
	return cmd
}