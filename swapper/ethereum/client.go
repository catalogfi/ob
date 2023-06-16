package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	// "crypto/rand"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/susruth/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/susruth/wbtc-garden/swapper/ethereum/typings/ERC20"
)

type Client interface {
	GetTransactOpts(privKey *ecdsa.PrivateKey) *bind.TransactOpts
	GetCallOpts() *bind.CallOpts
	RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error)
	RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address) (string, error)
	GetPublicAddress(privKey *ecdsa.PrivateKey) common.Address
	GetProvider() *ethclient.Client
	TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address, auth *bind.TransactOpts) (string, error)
	GetCurrentBlock() (uint64, error)
	GetERC20Balance(tokenAddr common.Address, address common.Address, auth *bind.CallOpts) (*big.Int, error)
}
type client struct {
	url      string
	provider *ethclient.Client
}

func NewClient(url string) Client {
	provider, _ := ethclient.Dial(url)
	return &client{url: url, provider: provider}
}
func (client *client) GetTransactOpts(privKey *ecdsa.PrivateKey) *bind.TransactOpts {
	provider := client.provider
	chainId, err := provider.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	fromAddress := client.GetPublicAddress(privKey)
	nonce, err := provider.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	// gasPrice, err := provider.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, chainId)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(20305454254)

	return auth
}

func (client *client) GetCallOpts() *bind.CallOpts {
	var auth *bind.CallOpts = &bind.CallOpts{}
	auth.Pending = true
	return auth
}
func (client *client) RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Redeem(auth, token, secret)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func (client *client) RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Refund(auth, token)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func (client *client) GetPublicAddress(privKey *ecdsa.PrivateKey) common.Address {
	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress
}
func (client *client) GetProvider() *ethclient.Client {
	return client.provider
}
func (client *client) TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address, auth *bind.TransactOpts) (string, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Transfer(auth, toAddr, amount)
	if err != nil {
		return "", err
	}
	fmt.Printf("Transfering %v %s to %s txhash : %s\n", amount, tokenAddr, toAddr, tx.Hash().Hex())
	return tx.Hash().Hex(), err
}
func (client *client) GetCurrentBlock() (uint64, error) {
	bn, err := client.provider.BlockNumber(context.Background())
	return bn, err
}

func (client *client) GetERC20Balance(tokenAddr common.Address, ofAddr common.Address, auth *bind.CallOpts) (*big.Int, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return big.NewInt(0), err
	}
	balance, err := instance.BalanceOf(auth, ofAddr)
	return balance, err

}
