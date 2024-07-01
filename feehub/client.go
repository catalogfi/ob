package feehub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/catalogfi/ob/model"
	"gorm.io/gorm"
)

type FeehubClient struct {
	BaseURL string
	Client  *http.Client
}
type ConditionalPayment struct {
	ChannelID          uint   `json:"channelId"`
	HTLC               HTLC   `json:"htlc"`
	UserStateSignature string `json:"userStateSignature"`
}

type SplitConditionalFee struct {
	OrderID     uint64             `json:"orderId"`
	SendPayment ConditionalPayment `json:"sendPayment"`
}
type HTLC struct {
	gorm.Model

	SecretHash    string       `json:"secretHash"`
	TimeLock      uint64       `json:"timeLock"`
	OnChainExpiry uint64       `json:"-"`
	SendAmount    model.BigInt `gorm:"type:decimal" default:"0" json:"sendAmount"`
	RecvAmount    model.BigInt `gorm:"type:decimal" default:"0" json:"receiveAmount"`
	OrderID       uint64       `json:"orderId"`
}

func NewFeehubClient(baseURL string) *FeehubClient {
	return &FeehubClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (fc *FeehubClient) PayFiller(ctx context.Context, ConditionalFee ConditionalPayment, fillerAddr string, auth string) error {
	var fee struct {
		status string `json:"status"`
	}

	err := fc.SubmitPost(ctx, "htlc", auth, ConditionalFee, &fee)
	if err != nil || fee.status != "ok" {
		return err
	}

	return nil
}
func (fc *FeehubClient) SubmitPost(ctx context.Context, endpoint string, auth string, req interface{}, resp interface{}) error {
	url := fmt.Sprintf("%s/%s", fc.BaseURL, endpoint)

	requestBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v,endpoint %s", err, endpoint)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v,endpoint %s", err, endpoint)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", auth)

	response, err := fc.Client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request: %v,endpoint %s", err, endpoint)
	}
	defer response.Body.Close()

	if err := parseResponse(response, &resp); err != nil {
		return fmt.Errorf("failed to parse response: %v,endpoint %s", err, endpoint)
	}

	return nil
}
func parseResponse(resp *http.Response, data interface{}) error {
	var res struct {
		Data  json.RawMessage `json:"data"`
		Error string          `json:"error"`
	}

	if resp.StatusCode != http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return fmt.Errorf("server returned non-200 ,failed to decode error, status: %d", resp.StatusCode)
		}
		return fmt.Errorf("server returned non-200, error: %s", res.Error)
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	if res.Error != "" {
		return fmt.Errorf("server returned error: %s", res.Error)
	}

	if data != nil {
		if err := json.Unmarshal(res.Data, data); err != nil {
			return err
		}
	}

	return nil
}
