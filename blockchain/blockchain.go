package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
)

func LoadClient(chain model.Chain, urls map[model.Chain]string) (interface{}, error) {
	if chain.IsBTC() {
		return bitcoin.NewClient(urls[chain], getParams(chain)), nil
	}
	if chain.IsEVM() {
		return ethereum.NewClient(urls[chain])
	}
	return nil, fmt.Errorf("invalid chain: %s", chain)
}

func LoadInitiatorSwap(atomicSwap model.AtomicSwap, initiatorPrivateKey interface{}, secretHash string, urls, deployers map[model.Chain]string) (swapper.InitiatorSwap, error) {
	client, err := LoadClient(atomicSwap.Chain, urls)
	if err != nil {
		fmt.Println(err)
	}

	redeemerAddress, err := ParseAddress(client, atomicSwap.RedeemerAddress)
	if err != nil {
		fmt.Println(err)
	}

	secHash, err := hex.DecodeString(secretHash)
	if err != nil {
		fmt.Println(err)
	}

	amt, ok := new(big.Int).SetString(atomicSwap.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", atomicSwap.Amount)
	}

	expiry, ok := new(big.Int).SetString(atomicSwap.Timelock, 10)
	if !ok {
		return nil, fmt.Errorf("invalid timelock: %s", atomicSwap.Timelock)
	}

	switch client := client.(type) {
	case bitcoin.Client:
		return bitcoin.NewInitiatorSwap(initiatorPrivateKey.(*btcec.PrivateKey), redeemerAddress.(btcutil.Address), secHash, expiry.Int64(), amt.Uint64(), client)
	case ethereum.Client:
		deployerAddress := common.HexToAddress(deployers[atomicSwap.Chain])
		tokenAddress := common.HexToAddress(atomicSwap.Asset.SecondaryID())
		return ethereum.NewInitiatorSwap(initiatorPrivateKey.(*ecdsa.PrivateKey), redeemerAddress.(common.Address), deployerAddress, tokenAddress, secHash, expiry, amt, client)
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func LoadWatcher(atomicSwap model.AtomicSwap, secretHash string, urls, deployers map[model.Chain]string) (swapper.Watcher, error) {
	client, err := LoadClient(atomicSwap.Chain, urls)
	if err != nil {
		return nil, fmt.Errorf("failed to load client: %v", err)
	}

	initiatorAddress, err := ParseAddress(client, atomicSwap.InitiatorAddress)
	if err != nil {
		return nil, err
	}

	redeemerAddress, err := ParseAddress(client, atomicSwap.RedeemerAddress)
	if err != nil {
		return nil, err
	}

	secHash, err := hex.DecodeString(secretHash)
	if err != nil {
		return nil, err
	}

	amt, ok := new(big.Int).SetString(atomicSwap.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", atomicSwap.Amount)
	}

	expiry, ok := new(big.Int).SetString(atomicSwap.Timelock, 10)
	if !ok {
		return nil, fmt.Errorf("invalid timelock: %s", atomicSwap.Timelock)
	}

	switch client := client.(type) {
	case bitcoin.Client:
		return bitcoin.NewWatcher(initiatorAddress.(btcutil.Address), redeemerAddress.(btcutil.Address), secHash, expiry.Int64(), amt.Uint64(), client)
	case ethereum.Client:
		deployerAddress := common.HexToAddress(deployers[atomicSwap.Chain])
		tokenAddress := common.HexToAddress(atomicSwap.Asset.SecondaryID())
		return ethereum.NewWatcher(initiatorAddress.(common.Address), redeemerAddress.(common.Address), deployerAddress, tokenAddress, secHash, expiry, amt, client)
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func CalculateExpiry(chain model.Chain, goingFirst bool, urls map[model.Chain]string) (string, error) {
	if chain.IsBTC() {
		expiry := bitcoin.GetExpiry(goingFirst)
		return strconv.FormatInt(expiry, 10), nil
	}
	client, err := LoadClient(chain, urls)
	if err != nil {
		return "", err
	}
	expiry, err := ethereum.GetExpiry(client.(ethereum.Client), goingFirst)
	if err != nil {
		return "", err
	}
	return expiry.String(), nil
}

func LoadRedeemerSwap(atomicSwap model.AtomicSwap, redeemerPrivateKey interface{}, secretHash string, urls, deployers map[model.Chain]string) (swapper.RedeemerSwap, error) {
	client, err := LoadClient(atomicSwap.Chain, urls)
	if err != nil {
		fmt.Println(err)
	}

	initiatorAddress, err := ParseAddress(client, atomicSwap.InitiatorAddress)
	if err != nil {
		fmt.Println(err)
	}

	secHash, err := hex.DecodeString(secretHash)
	if err != nil {
		fmt.Println(err)
	}

	amt, ok := new(big.Int).SetString(atomicSwap.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", atomicSwap.Amount)
	}

	expiry, ok := new(big.Int).SetString(atomicSwap.Timelock, 10)
	if !ok {
		return nil, fmt.Errorf("invalid timelock: %s", atomicSwap.Timelock)
	}

	switch client := client.(type) {
	case bitcoin.Client:
		return bitcoin.NewRedeemerSwap(redeemerPrivateKey.(*btcec.PrivateKey), initiatorAddress.(btcutil.Address), secHash, expiry.Int64(), amt.Uint64(), client)
	case ethereum.Client:
		deployerAddress := common.HexToAddress(deployers[atomicSwap.Chain])
		tokenAddress := common.HexToAddress(atomicSwap.Asset.SecondaryID())
		return ethereum.NewRedeemerSwap(redeemerPrivateKey.(*ecdsa.PrivateKey), initiatorAddress.(common.Address), deployerAddress, tokenAddress, secHash, expiry, amt, client)
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func ParseKey(chain model.Chain, key string) (interface{}, error) {
	switch chain {
	case model.Bitcoin:
		privKeyBytes, err := hex.DecodeString(key)
		if err != nil {
			return nil, err
		}
		// ignoring public key as we do not need it
		privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
		return privKey, nil
	case model.Ethereum:
		return crypto.HexToECDSA(key)
	default:
		return nil, fmt.Errorf("unknown chain: %s", chain)
	}
}

func ParseAddress(client interface{}, address string) (interface{}, error) {
	switch client := client.(type) {
	case bitcoin.Client:
		return btcutil.DecodeAddress(address, client.Net())
	case ethereum.Client:
		return common.HexToAddress(address), nil
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func getParams(chain model.Chain) *chaincfg.Params {
	switch chain {
	case model.Bitcoin:
		return &chaincfg.MainNetParams
	case model.BitcoinTestnet:
		return &chaincfg.TestNet3Params
	case model.BitcoinRegtest:
		return &chaincfg.RegressionNetParams
	default:
		panic("constraint violation: unknown chain")
	}
}
