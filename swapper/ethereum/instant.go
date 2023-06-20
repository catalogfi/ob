package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type instantClient struct {
	url           string
	indexerClient Client
}

func InstantWalletWrapper(url string, client Client) Client {
	return &instantClient{url: url, indexerClient: client}
}

func (client *instantClient) GetTransactOpts(privKey *ecdsa.PrivateKey) *bind.TransactOpts {
	return client.indexerClient.GetTransactOpts(privKey)
}

func (client *instantClient) GetCallOpts() *bind.CallOpts {
	return client.indexerClient.GetCallOpts()
}

func (client *instantClient) RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error) {
	return client.indexerClient.RedeemAtomicSwap(contract, auth, token, secret)
}

func (client *instantClient) RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address) (string, error) {
	return client.indexerClient.RefundAtomicSwap(contract, auth, token)
}

func (client *instantClient) GetPublicAddress(privKey *ecdsa.PrivateKey) common.Address {
	return client.indexerClient.GetPublicAddress(privKey)
}

func (client *instantClient) GetProvider() *ethclient.Client {
	return client.indexerClient.GetProvider()
}

func (client *instantClient) TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
	panic("not implemented")
}

func (client *instantClient) IsFinal(txHash string) (bool, error) {
	// TODO: check whether it is an instant wallet transaction, if it is return true, nil
	panic("not implemented")

	return client.indexerClient.IsFinal(txHash)
}

func (client *instantClient) GetCurrentBlock() (uint64, error) {
	return client.indexerClient.GetCurrentBlock()
}

func (client *instantClient) GetERC20Balance(tokenAddr common.Address, address common.Address) (*big.Int, error) {
	return client.indexerClient.GetERC20Balance(tokenAddr, address)
}
