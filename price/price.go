package price

import (
	"math"

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

func (p *PriceChecker) Run() error {
	// for {
	// 	resp, err := http.Get(p.url)
	// 	if err != nil {
	// 		fmt.Println("failed to get prices", err)
	// 		time.Sleep(5 * time.Second)
	// 		continue
	// 	}
	// 	defer resp.Body.Close()

	// 	var data map[string]map[string]float64
	// 	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
	// 		fmt.Println("failed to decode response", err)
	// 		time.Sleep(5 * time.Second)
	// 		continue
	// 	}

	// 	fmt.Println(data["bitcoin"]["usd"])
	// 	if err := p.store.SetPrice("bitcoin", "ethereum", float64(data["bitcoin"]["usd"])); err != nil {
	// 		return err
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }
	return p.store.SetPrice("bitcoin", "ethereum", float64(30000))
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
