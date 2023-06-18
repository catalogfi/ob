package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
)

func main() {
	initiatorAddr, err := btcutil.DecodeAddress("mg54DDo5jfNkx5tF4d7Ag6G6VrJaSjr7ES", &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	redeemerAddr, err := btcutil.DecodeAddress("mpifLw6JYHJoAfWTtqieaync4ZdRpBeDSc", &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}

	secretHash, err := hex.DecodeString("cacc4b29040889acee9404c7c93cf4bd2852d8e16fd8781a33b10de2b8e09e49")
	if err != nil {
		panic(err)
	}

	addr1, err := bitcoin.GetAddress(bitcoin.NewClient("https://blockstream.info/testnet/api", &chaincfg.TestNet3Params), redeemerAddr, initiatorAddr, secretHash, 144)
	if err != nil {
		panic(err)
	}

	addr2, err := bitcoin.GetAddress(bitcoin.NewClient("https://blockstream.info/testnet/api", &chaincfg.TestNet3Params), initiatorAddr, redeemerAddr, secretHash, 144)
	if err != nil {
		panic(err)
	}

	addr3, err := bitcoin.GetAddress(bitcoin.NewClient("https://blockstream.info/testnet/api", &chaincfg.TestNet3Params), redeemerAddr, initiatorAddr, secretHash, 288)
	if err != nil {
		panic(err)
	}

	addr4, err := bitcoin.GetAddress(bitcoin.NewClient("https://blockstream.info/testnet/api", &chaincfg.TestNet3Params), initiatorAddr, redeemerAddr, secretHash, 288)
	if err != nil {
		panic(err)
	}

	fmt.Println(addr1.String())
	fmt.Println(addr2.String())
	fmt.Println(addr3.String())
	fmt.Println(addr4.String())
}
