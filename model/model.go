package model

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Config struct {
	RPC map[Chain]string `json:"rpc"`
}

type Chain string

const (
	Bitcoin          Chain = "bitcoin"
	BitcoinTestnet   Chain = "bitcoin_testnet"
	BitcoinRegtest   Chain = "bitcoin_regtest"
	Ethereum         Chain = "ethereum"
	EthereumSepolia  Chain = "ethereum_sepolia"
	EthereumLocalnet Chain = "ethereum_localnet"
)

func ParseChain(c string) (Chain, error) {
	switch strings.ToLower(c) {
	case "bitcoin":
		return Bitcoin, nil
	case "bitcoin_testnet", "bitcoin-testnet", "bitcoin-testnet3":
		return BitcoinTestnet, nil
	case "bitcoin_regtest", "bitcoin-regtest", "bitcoin-localnet":
		return BitcoinRegtest, nil
	case "ethereum":
		return Ethereum, nil
	case "ethereum_sepolia", "sepolia", "ethereum-sepolia":
		return EthereumSepolia, nil
	case "ethereum_localnet", "ethereum-localnet":
		return EthereumLocalnet, nil
	default:
		return Chain(""), fmt.Errorf("unknown chain %v", c)
	}
}

func (c Chain) IsEVM() bool {
	return c == Ethereum || c == EthereumSepolia || c == EthereumLocalnet
}

func (c Chain) IsBTC() bool {
	return c == Bitcoin || c == BitcoinTestnet || c == BitcoinRegtest
}

type Asset string

const (
	Primary Asset = "primary"
)

func NewSecondary(address string) Asset {
	return Asset("secondary" + address)
}

func (a Asset) SecondaryID() string {
	if string(a) == "primary" || string(a[:9]) != "secondary" {
		return ""
	}
	return string(a[9:])
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

type VerifySiwe struct {
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type Order struct {
	gorm.Model

	Maker     string `json:"maker"`
	Taker     string `json:"taker"`
	OrderPair string `json:"orderPair"`

	InitiatorAtomicSwapID uint
	FollowerAtomicSwapID  uint
	InitiatorAtomicSwap   *AtomicSwap `json:"initiatorAtomicSwap" gorm:"foreignKey:InitiatorAtomicSwapID"`
	FollowerAtomicSwap    *AtomicSwap `json:"followerAtomicSwap" gorm:"foreignKey:FollowerAtomicSwapID"`

	SecretHash           string  `json:"secretHash" gorm:"unique;not null"`
	Secret               string  `json:"secret" gorm:"unique;not null"`
	Price                float64 `json:"price"`
	Status               Status  `json:"status"`
	SecretNonce          uint64  `json:"secretNonce"`
	UserBtcWalletAddress string  `json:"userBtcWalletAddress"`
}

type AtomicSwap struct {
	gorm.Model

	InitiatorAddress string `json:"initiatorAddress"`
	RedeemerAddress  string `json:"redeemerAddress"`
	Timelock         string `json:"timelock"`
	Chain            Chain  `json:"chain"`
	Asset            Asset  `json:"asset"`
	Amount           string `json:"amount"`
	InitiateTxHash   string `json:"initiateTxHash" gorm:"unique"`
	RedeemTxHash     string `json:"redeemTxHash" gorm:"unique"`
	RefundTxHash     string `json:"refundTxHash" gorm:"unique"`
}

type StringArray []string

func (sa StringArray) Value() (driver.Value, error) {
	return strings.Join(sa, ","), nil
}

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = make([]string, 0)
		return nil
	}

	switch v := value.(type) {
	case string:
		*sa = strings.Split(v, ",")
	case []byte:
		*sa = strings.Split(string(v), ",")
	default:
		return fmt.Errorf("unsupported data type for StringArray: %T", value)
	}

	return nil
}

func ParseOrderPair(orderPair string) (Chain, Chain, Asset, Asset, error) {
	chainAssets := strings.Split(orderPair, "-")
	if len(chainAssets) != 2 {
		return "", "", "", "", fmt.Errorf("failed to parse the order pair, should be of the format <chain>:<asset>-<chain>:<asset>. got: %v", orderPair)
	}
	sendChain, sendAsset, err := ParseChainAsset(chainAssets[0])
	if err != nil {
		return "", "", "", "", err
	}
	recieveChain, recieveAsset, err := ParseChainAsset(chainAssets[1])
	if err != nil {
		return "", "", "", "", err
	}
	return sendChain, recieveChain, sendAsset, recieveAsset, nil
}

func ParseChainAsset(chainAsset string) (Chain, Asset, error) {
	chainAndAsset := strings.Split(chainAsset, ":")
	if len(chainAndAsset) > 2 {
		return "", "", fmt.Errorf("failed to parse the chain and asset, should be of the format <chain>:<asset>. got: %v", chainAsset)
	}
	if len(chainAndAsset) == 1 {
		return Chain(chainAndAsset[0]), Primary, nil
	}
	return Chain(chainAndAsset[0]), NewSecondary(chainAndAsset[1]), nil
}
