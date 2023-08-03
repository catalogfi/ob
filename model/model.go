package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/susruth/wbtc-garden/config"
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
	EthereumOptimism Chain = "ethereum_optimism"
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
	case "ethereum_optimism" , "optimism" , "ethereum-optimism":
		return EthereumOptimism, nil
	default:
		return Chain(""), fmt.Errorf("unknown chain %v", c)
	}
}

func (c Chain) IsEVM() bool {
	return c == Ethereum || c == EthereumSepolia || c == EthereumLocalnet || c == EthereumOptimism
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
	OrderCancelled
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
	Secret               string  `json:"secret"`
	Price                float64 `json:"price"`
	Status               Status  `json:"status"`
	SecretNonce          uint64  `json:"secretNonce"`
	UserBtcWalletAddress string  `json:"userBtcWalletAddress"`

	Fee uint `json:"fee"`
}

type AtomicSwap struct {
	gorm.Model

	InitiatorAddress     string  `json:"initiatorAddress"`
	RedeemerAddress      string  `json:"redeemerAddress"`
	Timelock             string  `json:"timelock"`
	Chain                Chain   `json:"chain"`
	Asset                Asset   `json:"asset"`
	Amount               string  `json:"amount"`
	InitiateTxHash       string  `json:"initiateTxHash" `
	RedeemTxHash         string  `json:"redeemTxHash" `
	RefundTxHash         string  `json:"refundTxHash" `
	PriceByOracle        float64 `json:"priceByOracle"`
	MinimumConfirmations uint64  `josn:"minimumConfirmations"`
	IsInstantWallet      bool
}

type LockedAmount struct {
	Asset  string
	Amount sql.NullInt64
}

func CombineAndAddAmount(arr1, arr2 []LockedAmount) []LockedAmount {
	combinedMap := make(map[string]sql.NullInt64)

	if len(arr1) == 0 {
		return arr2
	} else if len(arr2) == 0 {
		return arr1
	} else if len(arr1) == 0 && len(arr2) == 0 {
		return nil
	}
	for _, item := range arr1 {
		if _, ok := combinedMap[item.Asset]; ok {
			combinedMap[item.Asset] = sql.NullInt64{
				Int64: combinedMap[item.Asset].Int64 + item.Amount.Int64,
				Valid: true,
			}
		} else {
			combinedMap[item.Asset] = item.Amount
		}
	}

	for _, item := range arr2 {
		if _, ok := combinedMap[item.Asset]; ok {
			combinedMap[item.Asset] = sql.NullInt64{
				Int64: combinedMap[item.Asset].Int64 + item.Amount.Int64,
				Valid: true,
			}
		} else {
			combinedMap[item.Asset] = item.Amount
		}
	}

	var combinedArray []LockedAmount
	for asset, amount := range combinedMap {
		combinedArray = append(combinedArray, LockedAmount{
			Asset:  asset,
			Amount: amount,
		})
	}

	return combinedArray
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
	if err := isWhitelisted(sendChain, strings.Replace(string(sendAsset), "secondary", "", 1)); err != nil {
		return "", "", "", "", err
	}
	if err := isWhitelisted(recieveChain, strings.Replace(string(recieveAsset), "secondary", "", 1)); err != nil {
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

func isWhitelisted(chain Chain, asset string) error {
	if chainMap, ok := config.ConfigMap[string(chain)]; ok {
		if _, ok := chainMap[asset]; ok {
			return nil
		}
		return fmt.Errorf("asset %v is not whitelisted for chain %v", asset, chain)
	} else {
		return fmt.Errorf("chain %v is not whitelisted", chain)
	}

}
