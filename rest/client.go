package rest

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
	"github.com/susruth/wbtc-garden/model"
)

type Client interface {
	FillOrder(orderID uint, sendAddress, recieveAddress string) error
	CreateOrder(sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) (uint, error)
	GetFollowerInitiateOrders(id string) ([]model.Order, error)
	GetFollowerRedeemOrders(id string) ([]model.Order, error)
	GetInitiatorInitiateOrders(id string) ([]model.Order, error)
	GetInitiatorRedeemOrders(id string) ([]model.Order, error)
	SetJwt(token string)
	Health() (string, error)
	GetNonce() (string, error)
	Verify(address string) (string, error)
}

type client struct {
	url      string
	JwtToken string
}

func NewClient(url string) Client {
	return &client{url: url, JwtToken: ""}
}

type ErrorResponse struct {
	Error string `json:"error"`
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
	req.Header.Set("Authorization", c.JwtToken)

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

func (c *client) CreateOrder(sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) (uint, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(CreateOrder{SendAddress: sendAddress, RecieveAddress: recieveAddress, OrderPair: orderPair, SendAmount: sendAmount, RecieveAmount: recieveAmount, SecretHash: secretHash}); err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orders", c.url), &buf)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", c.JwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return 0, fmt.Errorf("failed to decode error response: %v", err)
		}
		return 0, fmt.Errorf("failed to create order: %v", errorResponse["error"])
	}

	var orderResponse struct {
		ID uint `json:"orderId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orderResponse); err != nil {
		return 0, fmt.Errorf("failed to decode order response: %v", err)
	}
	return orderResponse.ID, nil
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
	resp, err := http.Get(fmt.Sprintf("%s/orders?taker=%s&status=5&verbose=true", c.url, id))
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
	resp, err := http.Get(fmt.Sprintf("%s/orders?maker=%s&status=2&verbose=true", c.url, id))
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
	resp, err := http.Get(fmt.Sprintf("%s/orders?maker=%s&status=4&verbose=true", c.url, id))
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

func (c *client) Health() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/health", c.url))
	if err != nil {
		return "", fmt.Errorf("failed to get health: %v", err)
	}
	defer resp.Body.Close()

	var HealthPayload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&HealthPayload); err != nil {
		return "", fmt.Errorf("failed to decode health: %v", err)
	}
	return HealthPayload.Status, nil
}

func (c *client) SetJwt(jwt string) {
	c.JwtToken = jwt
}

func (c *client) GetNonce() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/nonce", c.url))
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	defer resp.Body.Close()

	var NoncePayload struct {
		Nonce string `json:"nonce"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&NoncePayload); err != nil {
		return "", fmt.Errorf("failed to decode nonce: %v", err)
	}
	return NoncePayload.Nonce, nil
}

func (c *client) Verify(address string) (string, error) {
	nonce, err := c.GetNonce()
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}
	message, err := CreateEip4361TestMessage(address, nonce)
	if err != nil {
		return "", fmt.Errorf("failed to create message: %v", err)
	}

	ethPrivKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %v", err)
	}

	signature, err := signMessage(message.String(), ethPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(model.VerifySiwe{
		Message:   message.String(),
		Signature: hexutil.Encode(signature),
	}); err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/verify", c.url), "application/json", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to verify: %v", err)
	}
	defer resp.Body.Close()

	var VerifyPayload struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&VerifyPayload); err != nil {
		return "", fmt.Errorf("failed to decode verify: %v", err)
	}

	return VerifyPayload.Token, nil
}

func signHash(data []byte) common.Hash {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256Hash([]byte(msg))
}

func signMessage(message string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	sign := signHash([]byte(message))
	signature, err := crypto.Sign(sign.Bytes(), privateKey)

	if err != nil {
		return nil, err
	}

	signature[64] += 27
	return signature, nil
}

func CreateEip4361TestMessage(
	publicAddress string,
	nonce string,
) (*siwe.Message, error) {
	options := make(map[string]interface{})
	options["chainId"] = 11155111 // for now sepolia later will use config to swtich based on chains
	options["statement"] = "Sign into Catalog and experience limitless cross chain"
	message, err := siwe.InitMessage(
		"localhost:8080",
		publicAddress,
		"https://localhost:3000",
		nonce,
		options,
	)

	if err != nil {
		return nil, err
	}

	return message, nil
}
