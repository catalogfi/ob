package price

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PriceFetcher is an interface for fetching the price of btc in
type PriceFetcher interface {
	GetPrice(ctx context.Context) (PriceData, error)
}
type PriceData struct {
	Price         float64 `json:"amount"`
	L1BlockNumber uint64  `json:"ethBlock"`
}

type priceFetcher struct {
	options Options
}

type Options struct {
	URL          string
	DefaultPrice float64
}

func (options Options) WithDefaultPrice(defaultPrice float64) Options {
	options.DefaultPrice = defaultPrice
	return options
}

func NewPriceFetcher(options Options) PriceFetcher {
	return &priceFetcher{
		options: options,
	}
}

func (pf *priceFetcher) GetPrice(ctx context.Context) (PriceData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s?amount=1", pf.options.URL, "price"), nil)
	if err != nil {
		return PriceData{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PriceData{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PriceData{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res struct {
		Data  json.RawMessage `json:"data"`
		Error string          `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return PriceData{}, err
	}

	if res.Error != "" {
		return PriceData{}, fmt.Errorf("error: %s", res.Error)
	}

	var data PriceData

	if err := json.Unmarshal(res.Data, &data); err != nil {
		return PriceData{}, err
	}

	return data, nil
}
