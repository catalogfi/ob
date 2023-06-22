package model

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Chain string

const (
	Bitcoin  Chain = "bitcoin"
	Ethereum Chain = "ethereum"
)

type Asset string

const (
	Primary Asset = "primary"
)

func NewSecondary(address string) Asset {
	return Asset("secondary" + address)
}

type Status uint

const (
	Unknown Status = iota
	OrderCreated
	OrderFilled
	InitiatorAtomicSwapInitiated
	FollowerAtomicSwapInitiated
	FollowerAtomicSwapRedeemed
	InitiatorAtomicSwapRedeemed
	InitiatorAtomicSwapRefunded
	FollowerAtomicSwapRefunded
	OrderExecuted
	OrderFailedSoft
	OrderFailedHard
)

type Order struct {
	gorm.Model

	Maker     string `json:"maker"`
	Taker     string `json:"taker"`
	OrderPair string `json:"orderPair"`

	InitiatorAtomicSwapID uint
	FollowerAtomicSwapID  uint
	InitiatorAtomicSwap   *AtomicSwap `json:"initiatorAtomicSwap" gorm:"foreignKey:InitiatorAtomicSwapID"`
	FollowerAtomicSwap    *AtomicSwap `json:"followerAtomicSwap" gorm:"foreignKey:FollowerAtomicSwapID"`

	SecretHash string  `json:"secretHash"`
	Secret     string  `json:"secret"`
	Price      float64 `json:"price"`
	Status     Status  `json:"status"`
}

type AtomicSwap struct {
	gorm.Model

	InitiatorAddress string `json:"initiatorAddress"`
	RedeemerAddress  string `json:"redeemerAddress"`
	Timelock         string `json:"timelock"`
	Chain            Chain  `json:"chain"`
	Asset            Asset  `json:"asset"`
	Amount           string `json:"amount"`
	InitiateTxHash   string `json:"initiateTxHash"`
	RedeemTxHash     string `json:"redeemTxHash"`
	RefundTxHash     string `json:"refundTxHash"`
}

func ParseOrderPair(orderPair string) (Chain, Chain, Asset, Asset, error) {
	chainAssets := strings.Split(orderPair, "-")
	if len(chainAssets) != 2 {
		return "", "", "", "", fmt.Errorf("failed to parse the order pair, should be of the format <chain>:<asset>-<chain>:<asset>. got: %v", orderPair)
	}
	sendChain, sendAsset, err := parseChainAsset(chainAssets[0])
	if err != nil {
		return "", "", "", "", err
	}
	recieveChain, recieveAsset, err := parseChainAsset(chainAssets[1])
	if err != nil {
		return "", "", "", "", err
	}
	return sendChain, recieveChain, sendAsset, recieveAsset, nil
}

func parseChainAsset(chainAsset string) (Chain, Asset, error) {
	chainAndAsset := strings.Split(chainAsset, ":")
	if len(chainAndAsset) != 2 {
		return "", "", fmt.Errorf("failed to parse the chain and asset, should be of the format <chain>:<asset>. got: %v", chainAsset)
	}
	return Chain(chainAndAsset[0]), Asset(chainAndAsset[1]), nil
}

func NewOrderPair(from Chain, fromAsset Asset, to Chain, toAsset Asset) string {
	return fmt.Sprintf("%s:%s-%s:%s", from, fromAsset, to, toAsset)
}
