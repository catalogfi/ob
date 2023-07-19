package price

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/susruth/wbtc-garden/model"
)

func GetPriceInUSD(asset string, chain model.Chain) (float64, error) {

	var tokenId string
	if chain.IsEVM() {
		tokenId = ConfigMap[string(chain)][asset].ID
	} else if chain.IsBTC() {
		tokenId = "bitcoin"
	}
	resp, err := http.Get("https://api.coingecko.com/api/v3/simple/price?ids=" + tokenId + "&vs_currencies=usd")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, nil
	}
	return data[tokenId]["usd"], nil

}

func GetPrice(asset string, chain model.Chain, amount float64, PriceInUSD float64) float64 {

	var decimals float64
	if chain.IsEVM() {
		decimals = float64(ConfigMap[string(chain)][asset].Decimals)
	} else if chain.IsBTC() {
		decimals = 8
	}

	return (amount * PriceInUSD) / math.Pow(10, decimals)
}
