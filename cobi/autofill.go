package cobi

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

type Strategy struct {
	MaxFillOrders   uint     `json:"maxFillOrders"`
	MaxFillDeadline uint64   `json:"maxFillDeadline"`
	FromMaker       string   `json:"fromMaker"`
	FromChain       string   `json:"fromChain"`
	ToChain         string   `json:"toChain"`
	FromAsset       []string `json:"fromAsset"`
	ToAsset         []string `json:"toAsset"`
	MinFillAmount   float64  `json:"minFillAmount"`
	MaxFillAmount   float64  `json:"maxFillAmount"`
	OrderBy         string   `json:"orderBy"`
	FilterByPage    int      `json:"filterByPage"`
}

func AutoFill(entropy []byte, store Store) *cobra.Command {
	var (
		url      string
		account  uint32
		strategy string
	)
	var cmd = &cobra.Command{
		Use:   "autofill",
		Short: "fills the Orders based on strategy provided",
		Run: func(c *cobra.Command, args []string) {

			vals, err := getKeys(entropy, model.Ethereum, account, []uint32{0})
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting the signing key: %v", err))
			}
			privKey := vals[0].(*ecdsa.PrivateKey)
			client := rest.NewClient(url, privKey.D.Text(16))
			token, err := client.Login()
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting the signing key: %v", err))

			}
			if err := client.SetJwt(token); err != nil {
				cobra.CheckErr(fmt.Sprintf("Error to parse signing key: %v", err))

			}

			data, err := ioutil.ReadFile(strategy)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("error while reading strategy.json: %v", err))

			}

			var strategy Strategy
			err = json.Unmarshal(data, &strategy)
			if err != nil {
				cobra.CheckErr(fmt.Errorf("error while unmarshalling strategy.json: %v", err))

			}

			if strategy.MaxFillOrders == 0 {
				strategy.MaxFillOrders = math.MaxUint
			}

			if strategy.MaxFillDeadline == 0 {
				strategy.MaxFillDeadline = math.MaxUint64 - 1
			}

			if strategy.FilterByPage == 0 {
				strategy.FilterByPage = 100
			}

			var orders []model.Order
			totalOrdersFilled := uint(0)
			for {
				fmt.Println("Fetching Orders....")
				if uint64(time.Now().Unix()) > (strategy.MaxFillDeadline) {
					cobra.CheckErr("Max fill deadline reached")
				}
				if (len(strategy.FromAsset) == 1 && strategy.FromAsset[0] == "any") || (len(strategy.ToAsset) == 1 && strategy.ToAsset[0] == "any") {
					orders, err = GetAllAssets(
						client,
						strategy.FromMaker,
						strategy.FromChain,
						strategy.ToChain,
						strategy.FromAsset,
						strategy.ToAsset,
						strategy.MinFillAmount,
						strategy.MaxFillAmount,
						strategy.FilterByPage,
						strategy.OrderBy,
					)
					if err != nil {
						cobra.CheckErr(fmt.Errorf("Error while fetching Order: %v", err))

					}
				} else {
					orders = make([]model.Order, 0)
					for _, fromasset := range strategy.FromAsset {
						for _, toAsset := range strategy.ToAsset {
							orderPair := fmt.Sprintf("%s:%s-%s:%s", strategy.FromChain, fromasset, strategy.ToChain, toAsset)
							order, err := client.GetOrders(rest.GetOrdersFilter{
								Maker:     strategy.FromMaker,
								OrderPair: orderPair,
								OrderBy:   strategy.OrderBy,
								Verbose:   false,
								MinPrice:  strategy.MinFillAmount,
								MaxPrice:  strategy.MaxFillAmount,
								Status:    int(model.OrderCreated),
								PerPage:   strategy.FilterByPage,
							})
							if err != nil {
								cobra.CheckErr(fmt.Sprintf("Error while fetching Order: %v", err))
							}
							orders = append(orders, order...)
						}
					}
				}

				for _, order := range orders {
					fromChain, toChain, _, _, err := model.ParseOrderPair(order.OrderPair)
					if err != nil {
						cobra.CheckErr(fmt.Sprintf("Error while parsing order pair: %v", err))

					}

					toAddress, err := getAddressString(entropy, fromChain, account, 0)
					if err != nil {
						cobra.CheckErr(fmt.Sprintf("Error while getting address string: %v", err))

					}

					fromAddress, err := getAddressString(entropy, toChain, account, 0)
					if err != nil {
						cobra.CheckErr(fmt.Sprintf("Error while getting address string: %v", err))
					}
					if err := client.FillOrder(order.ID, fromAddress, toAddress); err != nil {
						cobra.CheckErr(fmt.Sprintf("Error while Filling the Order: %v with OrderID %d cross ❌", err, order.ID))
					}
					if err = store.PutSecretHash(order.SecretHash, uint64(order.ID)); err != nil {
						cobra.CheckErr(fmt.Sprintf("Error while storing secret hash: %v", err))
						return
					}
					totalOrdersFilled++
					if totalOrdersFilled >= strategy.MaxFillOrders {
						cobra.CheckErr("MaxFillOrders reached")

					}
					fmt.Printf("Filled order %d ✅", order.ID)
				}

				time.Sleep(15 * time.Second)
			}

		}}

	cmd.Flags().StringVar(&url, "url", "", "config file (default is ./config.json)")
	cmd.MarkFlagRequired("url")
	cmd.Flags().Uint32Var(&account, "account", 0, "config file (default: 0)")
	cmd.Flags().StringVar(&strategy, "strategy", "../../strategy.json", "config file (default: ./strategy.json)")
	return cmd
}

func GetAllAssets(
	client rest.Client,
	maker string,
	fromChain string,
	toChain string,
	fromAsset []string,
	toAsset []string,
	minPrice float64,
	maxPrice float64,
	fetchPerPage int,
	OrderBy string,
) ([]model.Order, error) {
	orders, err := client.GetOrders(rest.GetOrdersFilter{
		Maker:    maker,
		OrderBy:  OrderBy,
		Verbose:  false,
		Status:   int(model.OrderCreated),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		PerPage:  fetchPerPage,
	})
	if err != nil {
		fmt.Println("Error while getting orders:", err)
		return nil, err
	}

	filteredOrders := make([]model.Order, 0)
	for _, order := range orders {
		if len(fromAsset) == 1 && len(toAsset) == 1 {
			filteredOrders = append(filteredOrders, order)
			continue
		}

		orderPair := order.OrderPair
		FromChain, ToChain, FromAsset, ToAsset, err := model.ParseOrderPair(orderPair)
		if err != nil {
			fmt.Println("Error while parsing order pair:", err)
			return nil, err
		}
		if (contains(fromAsset, string(FromAsset)) || contains(toAsset, string(ToAsset))) && string(FromChain) == fromChain && string(ToChain) == toChain {
			filteredOrders = append(filteredOrders, order)
		}
	}
	return filteredOrders, nil

}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
