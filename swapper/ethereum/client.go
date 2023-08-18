package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	// "crypto/rand"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	maxApproval = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
)

type Client interface {
	GetTransactOpts(privKey *ecdsa.PrivateKey) *bind.TransactOpts
	GetCallOpts() *bind.CallOpts
	InitiateAtomicSwap(contract common.Address, initiator *ecdsa.PrivateKey, redeemerAddr, token common.Address, expiry *big.Int, amount *big.Int, secretHash []byte) (string, error)
	RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error)
	RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secretHash []byte) (string, error)
	GetPublicAddress(privKey *ecdsa.PrivateKey) common.Address
	GetProvider() *ethclient.Client
	ApproveERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error)
	TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error)
	GetCurrentBlock() (uint64, error)
	GetERC20Balance(tokenAddr common.Address, address common.Address) (*big.Int, error)
	GetTokenAddress(contractAddr common.Address) (common.Address, error)
	IsFinal(txHash string, waitBlocks uint64) (bool, error)
}
type client struct {
	url      string
	provider *ethclient.Client
}

func NewClient(url string) (Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	provider, err := ethclient.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return &client{url: url, provider: provider}, nil
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

	return auth
}

func (client *client) GetCallOpts() *bind.CallOpts {
	var auth *bind.CallOpts = &bind.CallOpts{}
	auth.Pending = true
	return auth
}
func (client *client) GetTokenAddress(contractAddr common.Address) (common.Address, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contractAddr, client.provider)
	if err != nil {
		return common.Address{}, err
	}
	tokenAddr, err := instance.Token(client.GetCallOpts())
	if err != nil {
		return common.Address{}, err
	}
	return tokenAddr, nil
}
func (client *client) RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Redeem(auth, secret)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}
func (client *client) InitiateAtomicSwap(contract common.Address, initiator *ecdsa.PrivateKey, redeemerAddr, token common.Address, expiry *big.Int, amount *big.Int, secretHash []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	var hash [32]byte
	copy(hash[:], secretHash)

	val, err := client.Allowance(token, contract, client.GetPublicAddress(initiator))
	if err != nil {
		return "", err
	}
	if val.Cmp(amount) <= 0 {
		_, err := client.ApproveERC20(initiator, maxApproval, token, contract)
		if err != nil {
			return "", err
		}
	}

	auth := client.GetTransactOpts(initiator)
	initTx, err := instance.Initiate(auth, redeemerAddr, expiry, amount, hash)

	if err != nil {
		return "", err
	}
	fmt.Println("hash", initTx.Hash().Hex())
	return initTx.Hash().Hex(), nil
}

func (client *client) RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secretHash []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	var hash [32]byte
	copy(hash[:], secretHash)
	tx, err := instance.Refund(auth, hash)
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
func (client *client) TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Transfer(client.GetTransactOpts(privKey), toAddr, amount)
	if err != nil {
		return "", err
	}
	fmt.Printf("Transfering %v %s to %s txhash : %s\n", amount, tokenAddr, toAddr, tx.Hash().Hex())
	return tx.Hash().Hex(), err
}
func (client *client) ApproveERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Approve(client.GetTransactOpts(privKey), toAddr, amount)
	if err != nil {
		return "", err
	}
	fmt.Printf("Approving %v %s to %s txhash : %s\n", amount, tokenAddr, toAddr, tx.Hash().Hex())
	bind.WaitMined(context.Background(), client.provider, tx)
	return tx.Hash().Hex(), err
}
func (client *client) Allowance(tokenAddr common.Address, spender common.Address, owner common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return nil, err
	}
	allowance, err := instance.Allowance(client.GetCallOpts(), owner, spender)
	if err != nil {
		return nil, err
	}
	return allowance, nil

}
func (client *client) GetCurrentBlock() (uint64, error) {
	bn, err := client.provider.BlockNumber(context.Background())
	return bn, err
}

func (client *client) GetERC20Balance(tokenAddr common.Address, ofAddr common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return big.NewInt(0), err
	}
	balance, err := instance.BalanceOf(client.GetCallOpts(), ofAddr)
	return balance, err
}

func (client *client) IsFinal(txHash string, waitBlocks uint64) (bool, error) {
	tx, err := client.provider.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return false, err
	}
	if tx.Status == 0 {
		return false, nil
	}
	currentBlock, err := client.GetCurrentBlock()
	if err != nil {
		return false, fmt.Errorf("error getting current block %v", err)
	}
	if int64(currentBlock-tx.BlockNumber.Uint64()) >= int64(waitBlocks)-1 || waitBlocks == 0 {
		return true, nil
	}
	return false, nil
}

func (client *client) FetchOrder(secretHash []byte) {
	// TODO:
}
