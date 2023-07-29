package price

import (
	"encoding/json"
	"fmt"
	"math"
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

func GetPrice(asset string, chain model.Chain, amount float64, PriceInUSD float64) float64 {

	var decimals float64
	if chain.IsEVM() {
		decimals = float64(config.ConfigMap[string(chain)][asset].Decimals)
	} else if chain.IsBTC() {
		decimals = 8
	}

	return (amount * PriceInUSD) / math.Pow(10, decimals)
}
