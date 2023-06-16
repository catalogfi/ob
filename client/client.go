package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
)

func main() {
	// btcToWbtc("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
	wbtcToBTC("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d")
}

func btcToWbtc(privKeyStr string) {
	privKey, err := hex.DecodeString(privKeyStr)
	if err != nil {
		panic(err)
	}
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKey)

	// addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(btcPrivKey.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
	// if err != nil {
	// 	panic(err)
	// }
	// panic(addr.EncodeAddress())

	ethPrivKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		panic(err)
	}
	ethAddr := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	// 0.0001 BTC to 0.0001 WBTC

	secret := [32]byte{}
	rand.Read(secret[:])
	secretHash := sha256.Sum256(secret[:])

	account := getAccount()

	redeemerPubKey, err := hex.DecodeString(account.BtcPubKey)
	if err != nil {
		panic(err)
	}

	initiatorPubKey := hex.EncodeToString(btcPrivKey.PubKey().SerializeCompressed())

	btcClient := bitcoin.NewClient("http://localhost:30000", &chaincfg.RegressionNetParams)
	ethClient := ethereum.NewClient("http://localhost:8545/")

	swap, err := bitcoin.NewInitiatorSwap(btcPrivKey, redeemerPubKey, secretHash[:], 288, 10000, btcClient)
	if err != nil {
		panic(err)
	}
	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := currBlock + 5760

	rSwap, err := ethereum.NewRedeemerSwap(ethPrivKey, common.HexToAddress(account.WbtcAddress), common.HexToAddress(account.WbtcTokenAddress), secretHash[:], big.NewInt(int64(expiry)), big.NewInt(9990), ethClient)
	if err != nil {
		panic(err)
	}

	txHash, err := swap.Initiate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Atomic Swap Initiated", txHash)

	postTransaction(rest.PostTransactionReq{
		From:       initiatorPubKey,
		To:         ethAddr.Hex(),
		SecretHash: hex.EncodeToString(secretHash[:]),
		Amount:     0.0001,
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

func wbtcToBTC(privKeyStr string) {
	// 0.1 WBTC to 0.1 BTC

	privKey, err := hex.DecodeString(privKeyStr)
	if err != nil {
		panic(err)
	}
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKey)

	ethPrivKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		panic(err)
	}
	ethAddr := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	fmt.Println(ethAddr.Hex())

	secret := [32]byte{}
	rand.Read(secret[:])
	secretHash := sha256.Sum256(secret[:])

	account := getAccount()

	redeemerPubKey, err := hex.DecodeString(account.BtcPubKey)
	if err != nil {
		panic(err)
	}

	initiatorPubKey := hex.EncodeToString(btcPrivKey.PubKey().SerializeCompressed())

	btcClient := bitcoin.NewClient("http://localhost:30000", &chaincfg.RegressionNetParams)
	ethClient := ethereum.NewClient("http://localhost:8545/")

	currBlock, err := ethClient.GetCurrentBlock()
	if err != nil {
		panic(err)
	}
	expiry := currBlock + 2880
	swap, err := ethereum.NewInitiatorSwap(ethPrivKey, common.HexToAddress(account.WbtcAddress), common.HexToAddress(account.WbtcTokenAddress), secretHash[:], big.NewInt(int64(expiry)), big.NewInt(10000), ethClient)
	if err != nil {
		panic(err)
	}

	rSwap, err := bitcoin.NewRedeemerSwap(btcPrivKey, redeemerPubKey, secretHash[:], 144, 9990, btcClient)
	if err != nil {
		panic(err)
	}

	txHash, err := swap.Initiate()
	if err != nil {
		panic(err)
	}
	fmt.Println("Atomic Swap Initiated", txHash)

	postTransaction(rest.PostTransactionReq{
		From:       ethAddr.Hex(),
		To:         initiatorPubKey,
		SecretHash: hex.EncodeToString(secretHash[:]),
		Amount:     0.0001,
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

func getAccount() model.Account {
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		panic(err)
	}

	account := model.Account{}

	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		panic(err)
	}

	return account
}

func postTransaction(req rest.PostTransactionReq) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:8080/transactions", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 201 {
		panic(fmt.Errorf("failed to create transaction: %s", resp.Status))
	}
}

func getTransactions(address string) []model.Transaction {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/transactions/%s", address))
	if err != nil {
		panic(err)
	}

	txs := []model.Transaction{}
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		panic(err)
	}

	return txs
}
