package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

	client := rest.NewChainClient(config.API)

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

	// btcToWbtc(config, client)
	wbtcToBTC(config, client)
}

func btcToWbtc(config ClientConfig, chainClient rest.ChainClient) {
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

	account, err := chainClient.GetAccount()
	if err != nil {
		panic(err)
	}

	redeemerAddr, err := btcutil.DecodeAddress(account.BtcAddress, config.BTCParams)
	if err != nil {
		panic(err)
	}

	btcClient := bitcoin.NewClient(config.BitcoinURL, config.BTCParams)
	ethClient, err := ethereum.NewClient(config.EthereumURL)
	if err != nil {
		panic(err)
	}

	swap, err := bitcoin.NewInitiatorSwap(btcPrivKey, redeemerAddr, secretHash[:], 288, 10000, btcClient)
	if err != nil {
		panic(err)
	}
	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := int64(currBlock + 5760)

	rSwap, err := ethereum.NewRedeemerSwap(ethPrivKey, common.HexToAddress(account.WbtcAddress), common.HexToAddress(account.DeployerAddress), common.HexToAddress(account.WbtcTokenAddress), secretHash[:], big.NewInt(int64(expiry)), big.NewInt(9990), ethClient)
	if err != nil {
		panic(err)
	}

	txHash, err := swap.Initiate()
	if err != nil {
		panic(err)
	}

	fmt.Println("Atomic Swap Initiated", txHash)

	if err := chainClient.PostTransaction(bitcoinAddress.EncodeAddress(), ethereumAddress.Hex(), hex.EncodeToString(secretHash[:]), expiry); err != nil {
		panic(err)
	}

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

func wbtcToBTC(config ClientConfig, chainClient rest.ChainClient) {
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

	account, err := chainClient.GetAccount()
	if err != nil {
		panic(err)
	}

	redeemerAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(btcPrivKey.PubKey().SerializeCompressed()), config.BTCParams)
	if err != nil {
		panic(err)
	}

	initiatorAddr, err := btcutil.DecodeAddress(account.BtcAddress, config.BTCParams)
	if err != nil {
		panic(err)
	}

	btcClient := bitcoin.NewClient(config.BitcoinURL, config.BTCParams)
	ethClient, err := ethereum.NewClient(config.EthereumURL)
	if err != nil {
		panic(err)
	}

	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := int64(currBlock + 2880)
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

	if err := chainClient.PostTransaction(ethAddr.Hex(), redeemerAddr.EncodeAddress(), hex.EncodeToString(secretHash[:]), expiry); err != nil {
		panic(err)
	}

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
