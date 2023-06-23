package blockchain

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

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

func LoadClient(chain model.Chain) (interface{}, error) {
	if chain == model.Bitcoin {
		vals := strings.Split(os.Getenv("BITCOIN_RPC"), "-")
		var params *chaincfg.Params
		switch vals[0] {
		case "mainnet":
			params = &chaincfg.MainNetParams
		case "testnet":
			params = &chaincfg.TestNet3Params
		case "regtest":
			params = &chaincfg.RegressionNetParams
		default:
			return nil, fmt.Errorf("invalid bitcoin network: %s", vals[0])
		}
		return bitcoin.NewClient(vals[1], params), nil
	}
	if chain.IsEVM() {
		fmt.Println( "ERORR" , fmt.Sprintf("%s_RPC", strings.ToUpper(string(chain))))
		return ethereum.NewClient(os.Getenv(fmt.Sprintf("%s_RPC", strings.ToUpper(string(chain)))))
	}
	return nil, fmt.Errorf("invalid chain: %s", chain)
}

func LoadInitiatorSwap(atomicSwap model.AtomicSwap, initiatorKey, secretHash string) (swapper.InitiatorSwap, error) {
	client, err := LoadClient(atomicSwap.Chain)
	if err != nil {
		fmt.Println(err)
	}

	initiatorPrivateKey, err := ParseKey(atomicSwap.Chain, initiatorKey)
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
		deployerAddress := common.HexToAddress(os.Getenv(fmt.Sprintf("%s_DEPLOYER_ADDRESS", strings.ToUpper(string(atomicSwap.Chain)))))
		tokenAddress := common.HexToAddress(atomicSwap.Asset.SecondaryID())
		return ethereum.NewInitiatorSwap(initiatorPrivateKey.(*ecdsa.PrivateKey), redeemerAddress.(common.Address), deployerAddress, tokenAddress, secHash, expiry, amt, client)
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func LoadWatcher(atomicSwap model.AtomicSwap, secretHash string) (swapper.Watcher, error) {
	client, err := LoadClient(atomicSwap.Chain)
	if err != nil {
		fmt.Println(err)
	}

	initiatorAddress, err := ParseAddress(client, atomicSwap.InitiatorAddress)
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
		return bitcoin.NewWatcher(initiatorAddress.(btcutil.Address), redeemerAddress.(btcutil.Address), secHash, expiry.Int64(), amt.Uint64(), client)
	case ethereum.Client:
		deployerAddress := common.HexToAddress(os.Getenv(fmt.Sprintf("%s_DEPLOYER_ADDRESS", strings.ToUpper(string(atomicSwap.Chain)))))
		tokenAddress := common.HexToAddress(atomicSwap.Asset.SecondaryID())
		return ethereum.NewWatcher(initiatorAddress.(common.Address), redeemerAddress.(common.Address), deployerAddress, tokenAddress, secHash, expiry, amt, client)
	default:
		return nil, fmt.Errorf("unknown chain: %T", client)
	}
}

func CalculateExpiry(chain model.Chain, goingFirst bool) (string, error) {
	if chain == model.Bitcoin {
		expiry := bitcoin.GetExpiry(goingFirst)
		return strconv.FormatInt(expiry, 10), nil
	}
	client, err := LoadClient(chain)
	if err != nil {
		return "", err
	}
	expiry, err := ethereum.GetExpiry(client.(ethereum.Client), goingFirst)
	if err != nil {
		return "", err
	}
	return expiry.String(), nil
}

func LoadRedeemerSwap(atomicSwap model.AtomicSwap, redeemerKey, secretHash string) (swapper.RedeemerSwap, error) {
	client, err := LoadClient(atomicSwap.Chain)
	if err != nil {
		fmt.Println(err)
	}

	redeemerPrivateKey, err := ParseKey(atomicSwap.Chain, redeemerKey)
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
		deployerAddress := common.HexToAddress(os.Getenv(fmt.Sprintf("%s_DEPLOYER_ADDRESS", strings.ToUpper(string(atomicSwap.Chain)))))
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
