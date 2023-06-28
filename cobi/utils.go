package cobi

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
	"github.com/tyler-smith/go-bip32"
)

func getKeys(entropy []byte, chain model.Chain, user uint32, selector []uint32) ([]interface{}, error) {
	masterKey, err := bip32.NewMasterKey(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %v", err)
	}

	var key *bip32.Key

	switch chain {
	case model.Bitcoin:
		key, err = masterKey.NewChildKey(0)
		if err != nil {
			return nil, fmt.Errorf("failed to create child key: %v", err)
		}
	case model.BitcoinTestnet, model.BitcoinRegtest:
		key, err = masterKey.NewChildKey(1)
		if err != nil {
			return nil, fmt.Errorf("failed to create child key: %v", err)
		}
	case model.Ethereum, model.EthereumLocalnet, model.EthereumSepolia:
		key, err = masterKey.NewChildKey(60)
		if err != nil {
			return nil, fmt.Errorf("failed to create child key: %v", err)
		}
	default:
		return nil, fmt.Errorf("invalid chain: %s", chain)
	}

	key, err = key.NewChildKey(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create child key: %v", err)
	}

	keys := make([]interface{}, len(selector))
	for i, sel := range selector {
		childKey, err := key.NewChildKey(sel)
		if err != nil {
			return nil, fmt.Errorf("failed to create %d child key: %v", sel, err)
		}

		if chain.IsBTC() {
			privKey, _ := btcec.PrivKeyFromBytes(childKey.PublicKey().Key)
			keys[i] = privKey
		} else if chain.IsEVM() {
			privKey, err := crypto.ToECDSA(childKey.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to create private key: %v", err)
			}
			keys[i] = privKey
		} else {
			return nil, fmt.Errorf("unsupported chain: %s", chain)
		}
	}
	return keys, nil
}

func getAddresses(entropy []byte, chain model.Chain, user uint32, selector []uint32) ([]interface{}, error) {
	keys, err := getKeys(entropy, chain, user, selector)
	if err != nil {
		return nil, err
	}

	addrs := make([]interface{}, len(keys))

	for i, key := range keys {
		if chain.IsBTC() {
			addrs[i], err = btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.(*btcec.PrivateKey).PubKey().SerializeCompressed()), getParams(chain))
			if err != nil {
				return nil, fmt.Errorf("failed to create address: %v", err)
			}
		} else if chain.IsEVM() {
			addrs[i] = crypto.PubkeyToAddress(key.(*ecdsa.PrivateKey).PublicKey)
		} else {
			return nil, fmt.Errorf("unsupported chain: %s", chain)
		}
	}
	return addrs, nil
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

func getBalances(entropy []byte, chain model.Chain, user uint32, selector []uint32, config model.Config, asset model.Asset) ([]interface{}, []uint64, error) {
	keys, err := getKeys(entropy, chain, user, selector)
	if err != nil {
		return nil, nil, err
	}

	addrs := make([]interface{}, len(keys))
	balances := make([]uint64, len(keys))

	client, err := blockchain.LoadClient(chain, config.RPC)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load client: %v", err)
	}

	for i, key := range keys {
		switch client := client.(type) {
		case bitcoin.Client:
			address, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.(*btcec.PrivateKey).PubKey().SerializeCompressed()), getParams(chain))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create address: %v", err)
			}

			_, balance, err := client.GetUTXOs(address, 0)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get UTXOs: %v", err)
			}
			addrs[i] = address
			balances[i] = balance
		case ethereum.Client:
			address := crypto.PubkeyToAddress(key.(*ecdsa.PrivateKey).PublicKey)
			if asset == model.Asset("primary") {
				balance, err := client.GetProvider().BalanceAt(context.Background(), address, nil)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to get ETH balance: %v", err)
				}
				addrs[i] = address
				balances[i] = balance.Uint64()
			}
		default:
			return nil, nil, fmt.Errorf("unsupported chain: %s", chain)
		}
	}
	return addrs, balances, nil
}

func getAddressString(entropy []byte, chain model.Chain, account, selector uint32) (string, error) {
	fromAddrs, err := getAddresses(entropy, chain, account, []uint32{selector})
	if err != nil {
		return "", fmt.Errorf("Error while getting addresses: %v", err)
	}
	switch addr := fromAddrs[0].(type) {
	case common.Address:
		return addr.Hex(), nil
	case btcutil.Address:
		return addr.EncodeAddress(), nil
	default:
		return "", fmt.Errorf("Error while getting addresses: %v", err)
	}
}
