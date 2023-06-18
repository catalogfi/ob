package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
)

type ClientConfig struct {
	PrivateKey      string `json:"privateKey"`
	API             string `json:"api"`
	BitcoinURL      string `json:"btcURL"`
	EthereumURL     string `json:"ethURL"`
	Network         string `json:"network"`
	DeployerAddress string `json:"deployerAddress"`

	BTCParams *chaincfg.Params
}

func main() {
	confFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("error reading config file (%s): %v", os.Args[1], err))
	}

	config := ClientConfig{}
	if err := json.Unmarshal(confFile, &config); err != nil {
		panic(fmt.Sprintf("error parsing config file (%s): %v", os.Args[1], err))
	}

	switch config.Network {
	case "regtest":
		config.BTCParams = &chaincfg.RegressionNetParams
	case "testnet":
		config.BTCParams = &chaincfg.TestNet3Params
	case "mainnet":
		config.BTCParams = &chaincfg.MainNetParams
	default:
		panic(fmt.Sprintf("invalid network: %s", config.Network))
	}

	// btcToWbtc(config)
	wbtcToBTC(config)
}

func btcToWbtc(config ClientConfig) {
	privKey, err := hex.DecodeString(config.PrivateKey)
	if err != nil {
		panic(err)
	}
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKey)

	bitcoinAddress, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(btcPrivKey.PubKey().SerializeCompressed()), config.BTCParams)
	if err != nil {
		panic(err)
	}
	fmt.Println(bitcoinAddress.EncodeAddress())

	ethPrivKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		panic(err)
	}
	ethereumAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	// 0.0001 BTC to 0.0001 WBTC

	secret := [32]byte{}
	rand.Read(secret[:])
	secretHash := sha256.Sum256(secret[:])

	account := getAccount(config.API)

	redeemerAddr, err := btcutil.DecodeAddress(account.BtcAddress, config.BTCParams)
	if err != nil {
		panic(err)
	}

	btcClient := bitcoin.NewClient(config.BitcoinURL, config.BTCParams)
	ethClient := ethereum.NewClient(config.EthereumURL)

	swap, err := bitcoin.NewInitiatorSwap(btcPrivKey, redeemerAddr, secretHash[:], 288, 10000, btcClient)
	if err != nil {
		panic(err)
	}
	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := currBlock + 5760

	rSwap, err := ethereum.NewRedeemerSwap(ethPrivKey, common.HexToAddress(account.WbtcAddress), common.HexToAddress(account.DeployerAddress), common.HexToAddress(account.WbtcTokenAddress), secretHash[:], big.NewInt(int64(expiry)), big.NewInt(9990), ethClient)
	if err != nil {
		panic(err)
	}

	txHash, err := swap.Initiate()
	if err != nil {
		panic(err)
	}

	fmt.Println("Atomic Swap Initiated", txHash)

	postTransaction(config.API, rest.PostTransactionReq{
		From:       bitcoinAddress.EncodeAddress(),
		To:         ethereumAddress.Hex(),
		SecretHash: hex.EncodeToString(secretHash[:]),
		WBTCExpiry: float64(expiry),
	})

	itxHash, err := rSwap.WaitForInitiate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Counter Party Initiated", itxHash)
	rtxHash, err := rSwap.Redeem(secret[:])
	if err != nil {
		panic(err)
	}
	fmt.Println("Atomic Swap Redeemed", rtxHash)
}

func wbtcToBTC(config ClientConfig) {
	// 0.1 WBTC to 0.1 BTC

	privKey, err := hex.DecodeString(config.PrivateKey)
	if err != nil {
		panic(err)
	}
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKey)

	ethPrivKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		panic(err)
	}
	ethAddr := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	fmt.Println(ethAddr.Hex())

	secret := [32]byte{}
	rand.Read(secret[:])
	secretHash := sha256.Sum256(secret[:])

	account := getAccount(config.API)

	redeemerAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(btcPrivKey.PubKey().SerializeCompressed()), config.BTCParams)
	if err != nil {
		panic(err)
	}

	initiatorAddr, err := btcutil.DecodeAddress(account.BtcAddress, config.BTCParams)
	if err != nil {
		panic(err)
	}

	btcClient := bitcoin.NewClient(config.BitcoinURL, config.BTCParams)
	ethClient := ethereum.NewClient(config.EthereumURL)

	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := currBlock + 2880
	swap, err := ethereum.NewInitiatorSwap(ethPrivKey, common.HexToAddress(account.WbtcAddress), common.HexToAddress(account.DeployerAddress), common.HexToAddress(account.WbtcTokenAddress), secretHash[:], big.NewInt(int64(expiry)), big.NewInt(10000), ethClient)
	if err != nil {
		panic(err)
	}

	rSwap, err := bitcoin.NewRedeemerSwap(btcPrivKey, initiatorAddr, secretHash[:], 144, 9990, btcClient)
	if err != nil {
		panic(err)
	}

	txHash, err := swap.Initiate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Atomic Swap Initiated", txHash)

	postTransaction(config.API, rest.PostTransactionReq{
		From:       ethAddr.Hex(),
		To:         redeemerAddr.EncodeAddress(),
		SecretHash: hex.EncodeToString(secretHash[:]),
		WBTCExpiry: float64(expiry),
	})

	itxHash, err := rSwap.WaitForInitiate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Counter Party Initiated", itxHash)
	rtxHash, err := rSwap.Redeem(secret[:])
	if err != nil {
		panic(err)
	}
	fmt.Println("Atomic Swap Redeemed", rtxHash)
}

func getAccount(api string) model.Account {
	resp, err := http.Get(api)
	if err != nil {
		panic(err)
	}

	account := model.Account{}

	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		panic(err)
	}

	return account
}

func postTransaction(api string, req rest.PostTransactionReq) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/transactions", api), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 201 {
		panic(fmt.Errorf("failed to create transaction: %s", resp.Status))
	}
}

func getTransactions(api, address string) []model.Transaction {
	resp, err := http.Get(fmt.Sprintf("%s/transactions/%s", api, address))
	if err != nil {
		panic(err)
	}

	txs := []model.Transaction{}
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		panic(err)
	}

	return txs
}
