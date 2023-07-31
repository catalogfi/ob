package price

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/susruth/wbtc-garden/config"
	"github.com/susruth/wbtc-garden/model"
)

type Store interface {
	SetPrice(fromChain string, toChain string, price float64) error
	Price(fromChain string, toChain string) (float64, error)
}

type PriceChecker struct {
	store Store
	url   string
}

func NewPriceChecker(store Store, url string) *PriceChecker {
	return &PriceChecker{store: store, url: url}
}

type ApiResponse struct {
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func (p *PriceChecker) Run() error {
	for {
		resp, err := http.Get(p.url)
		if err != nil {
			fmt.Println("failed to get prices", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		var apiResponse ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			fmt.Println("failed to decode response", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Convert priceUsd to float64
		priceUsdStr, ok := apiResponse.Data["priceUsd"].(string)
		if !ok {
			fmt.Println("failed to convert priceUsd to string")
			time.Sleep(5 * time.Second)
			continue
		}

		priceUsd, err := strconv.ParseFloat(priceUsdStr, 64)
		if err != nil {
			fmt.Println("failed to convert priceUsd to float64", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Println(priceUsd)

		if err := p.store.SetPrice("bitcoin", "ethereum", priceUsd); err != nil {
			return err
		}

		time.Sleep(10 * time.Second)
	}
	// return p.store.SetPrice("bitcoin", "ethereum", float64(30000))
}

func GetPrice(asset model.Asset, chain model.Chain, amount *big.Int, PriceInUSD *big.Int) *big.Int {

	var decimals int64
	var assetId string
	if asset == model.Primary {
		assetId = "primary"
	} else {
		assetId = asset.SecondaryID()
	}
	if chain.IsEVM() {
		decimals = int64(config.ConfigMap[string(chain)][assetId].Decimals)
	} else if chain.IsBTC() {
		decimals = 8
	}

	return new(big.Int).Div(new(big.Int).Mul(PriceInUSD, amount), new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}
