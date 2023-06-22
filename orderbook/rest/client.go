package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/susruth/wbtc-garden/orderbook/model"
)

type Client interface {
	FillOrder(orderID uint, sendAddress, recieveAddress string) error
	CreateOrder(sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) error
}

type client struct {
	url string
}

func NewClient(url string) Client {
	return &client{url: url}
}

func (c *client) FillOrder(orderID uint, sendAddress, recieveAddress string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(FillOrder{SendAddress: sendAddress, RecieveAddress: recieveAddress}); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/orders/%d", c.url, orderID), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fill order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to fill order: %v", resp.Status)
	}
	return nil
}

func (c *client) CreateOrder(sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(CreateOrder{SendAddress: sendAddress, RecieveAddress: recieveAddress, OrderPair: orderPair, SendAmount: sendAmount, RecieveAmount: recieveAmount, SecretHash: secretHash}); err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/orders", c.url), "application/json", &buf)
	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create order: %v", resp.Status)
	}
	return nil
}

func (c *client) GetFollowerInitiateOrders(id string) ([]model.Order, error) {
	resp, err := http.Get(fmt.Sprintf("%s/orders?taker=%s&status=3&verbose=true", c.url, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer resp.Body.Close()
	var orders []model.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}
	return orders, nil
}

func (c *client) GetFollowerRedeemOrders(id string) ([]model.Order, error) {
	resp, err := http.Get(fmt.Sprintf("%s/orders?taker=%s&status=3&verbose=true", c.url, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer resp.Body.Close()
	var orders []model.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}
	return orders, nil
}

func (c *client) GetInitiatorInitiateOrders(id string) ([]model.Order, error) {
	resp, err := http.Get(fmt.Sprintf("%s/orders?maker=%s&status=3&verbose=true", c.url, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer resp.Body.Close()
	var orders []model.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}
	return orders, nil
}

func (c *client) GetInitiatorRedeemOrders(id string) ([]model.Order, error) {
	resp, err := http.Get(fmt.Sprintf("%s/orders?maker=%s&status=5&verbose=true", c.url, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	defer resp.Body.Close()
	var orders []model.Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %v", err)
	}
	return orders, nil
}
