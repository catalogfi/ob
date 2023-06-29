package cobi

import (
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

func List() *cobra.Command {
	var (
		url        string
		maker      string
		taker      string
		orderPair  string
		secretHash string
		orderBy    string
		minPrice   float64
		maxPrice   float64
		page       int
		perPage    int
	)

	var cmd = &cobra.Command{
		Use: "list",
		Run: func(c *cobra.Command, args []string) {
			privKey, err := crypto.GenerateKey()
			if err != nil {
				cobra.CheckErr(err)
				return
			}
			orders, err := rest.NewClient(url, privKey.D.Text(16)).GetOrders(rest.GetOrdersFilter{
				Maker:      maker,
				Taker:      taker,
				OrderPair:  orderPair,
				SecretHash: secretHash,
				OrderBy:    orderBy,
				MinPrice:   minPrice,
				MaxPrice:   maxPrice,
				Page:       page,
				PerPage:    perPage,
				Verbose:    true,
				Status:     int(model.OrderCreated),
			})
			if err != nil {
				cobra.CheckErr(err)
				return
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Order ID", "From Asset", "To Asset", "Price", "Amount"})
			rows := make([]table.Row, len(orders))
			for i, order := range orders {
				assets := strings.Split(order.OrderPair, "-")
				rows[i] = table.Row{order.ID, assets[0], assets[1], order.Price, order.FollowerAtomicSwap.Amount}
			}
			t.AppendRows(rows)
			t.Render()
		},
		DisableAutoGenTag: true,
	}

	cmd.Flags().StringVar(&url, "url", "", "config file (default is ./config.json)")
	cmd.MarkFlagRequired("url")
	cmd.Flags().StringVar(&maker, "maker", "", "maker address to filter with (default: any)")
	cmd.Flags().StringVar(&taker, "taker", "", "taker address to filter with (default: any)")
	cmd.Flags().StringVar(&orderPair, "order-pair", "", "order pair to filter with (default: any)")
	cmd.Flags().StringVar(&secretHash, "secret-hash", "", "secret-hash to filter with (default: any)")
	cmd.Flags().StringVar(&orderBy, "order-by", "", "order by (default: creation time)")
	cmd.Flags().Float64Var(&minPrice, "min-price", 0, "minimum price to filter with (default: any)")
	cmd.Flags().Float64Var(&maxPrice, "max-price", 0, "maximum price to filter with (default: any)")
	cmd.Flags().IntVar(&page, "page", 1, "page number (default: 0)")
	cmd.Flags().IntVar(&perPage, "per-page", 10, "per page number (default: 10)")
	return cmd
}
