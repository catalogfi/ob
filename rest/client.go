package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/susruth/wbtc-garden/model"
)

type Client interface {
	Health() error
	GetChains() ([]string, error)
	ChainClient(chain string) ChainClient
}

type ChainClient interface {
	GetAccount() (model.Account, error)
	GetAddresses(from, to, secretHash string, wbtcExpiry int64) (model.HTLCAddresses, error)
	PostTransaction(from, to, secretHash string, wbtcExpiry int64) error
	GetTransactions(address string) ([]model.Transaction, error)
}

type client struct {
	url string
}

func New(url string) Client {
	return &client{url: url}
}

func (c *client) Health() error {
	resp, err := http.Get(fmt.Sprintf("%s/health", c.url))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status code %d", resp.StatusCode)
	}
	return nil
}

func (c *client) GetChains() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/chains", c.url))
	if err != nil {
		return nil, err
	}
	chains := []string{}
	if err := json.NewDecoder(resp.Body).Decode(&chains); err != nil {
		return nil, err
	}
	return chains, nil
}

func (c *client) ChainClient(chain string) ChainClient {
	return NewChainClient(fmt.Sprintf("%s/%s", c.url, chain))
}

func NewChainClient(url string) ChainClient {
	return &chainClient{url: url}
}

type chainClient struct {
	url string
}

func (c *chainClient) GetAccount() (model.Account, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return model.Account{}, err
	}
	account := model.Account{}
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (c *chainClient) GetAddresses(from, to, secretHash string, wbtcExpiry int64) (model.HTLCAddresses, error) {
	resp, err := http.Get(fmt.Sprintf("%s/addresses", c.url))
	if err != nil {
		return model.HTLCAddresses{}, err
	}
	addrs := model.HTLCAddresses{}
	if err := json.NewDecoder(resp.Body).Decode(&addrs); err != nil {
		return model.HTLCAddresses{}, err
	}
	return addrs, nil
}

func (c *chainClient) PostTransaction(from, to, secretHash string, wbtcExpiry int64) error {
	reqBytes, err := json.Marshal(PostTransactionReq{
		From:       from,
		To:         to,
		SecretHash: secretHash,
		WBTCExpiry: float64(wbtcExpiry),
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/transactions", c.url), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create transaction: %s", resp.Status)
	}
	return nil
}

func (c *chainClient) GetTransactions(address string) ([]model.Transaction, error) {
	resp, err := http.Get(fmt.Sprintf("%s/transactions/%s", c.url, address))
	if err != nil {
		return nil, err
	}
	txs := []model.Transaction{}
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		return nil, err
	}
	return txs, nil
}
