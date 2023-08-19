package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/ERC20"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

var (
	maxApproval = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
)

type Client interface {
	GetTransactOpts(privKey *ecdsa.PrivateKey) (*bind.TransactOpts, error)
	GetCurrentBlock() (uint64, error)
	GetProvider() *ethclient.Client
	GetTokenAddress(contractAddr common.Address) (common.Address, error)
	GetERC20Balance(tokenAddr common.Address, address common.Address) (*big.Int, error)
	GetDecimals(tokenAddr common.Address) (uint8, error)
	ApproveERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error)
	InitiateAtomicSwap(contract common.Address, initiator *ecdsa.PrivateKey, redeemerAddr, token common.Address, expiry *big.Int, amount *big.Int, secretHash []byte) (string, error)
	RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, orderID [32]byte, secret []byte) (string, error)
	RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secretHash []byte) (string, error)
	IsFinal(txHash string, waitBlocks uint64) (bool, uint64, error)
}

type client struct {
	logger   *zap.Logger
	url      string
	provider *ethclient.Client
	chainID  *big.Int
}

func NewClient(logger *zap.Logger, url string) (Client, error) {
	provider, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	childLogger := logger.With(zap.String("service", "ethClient"))
	chainID, err := provider.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &client{
		logger:   childLogger,
		url:      url,
		provider: provider,
		chainID:  chainID,
	}, nil
}

func (client *client) GetTransactOpts(privKey *ecdsa.PrivateKey) (*bind.TransactOpts, error) {
	transactor, err := bind.NewKeyedTransactorWithChainID(privKey, client.chainID)
	if err != nil {
		return nil, err
	}
	return transactor, nil
}

func (client *client) GetCurrentBlock() (uint64, error) {
	return client.provider.BlockNumber(context.Background())
}

func (client *client) GetProvider() *ethclient.Client {
	return client.provider
}

func (client *client) GetTokenAddress(contractAddr common.Address) (common.Address, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contractAddr, client.provider)
	if err != nil {
		return common.Address{}, err
	}
	return instance.Token(client.callOpts())
}

func (client *client) GetERC20Balance(tokenAddr common.Address, userAddr common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return big.NewInt(0), err
	}
	return instance.BalanceOf(client.callOpts(), userAddr)
}
func (client *client) GetDecimals(tokenAddr common.Address) (uint8, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return 0, err
	}
	return instance.Decimals(client.callOpts())
}

func (client *client) ApproveERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return "", err
	}
	transactor, err := client.GetTransactOpts(privKey)
	if err != nil {
		return "", err
	}
	tx, err := instance.Approve(transactor, toAddr, amount)
	if err != nil {
		return "", err
	}
	client.logger.Debug("approve erc20",
		zap.String("amount", amount.String()),
		zap.String("token address", tokenAddr.Hex()),
		zap.String("to address", toAddr.Hex()),
		zap.String("txHash", tx.Hash().Hex()))
	receipt, err := bind.WaitMined(context.Background(), client.provider, tx)
	if err != nil {
		return "", err
	}
	return receipt.TxHash.Hex(), err
}

func (client *client) InitiateAtomicSwap(contract common.Address, initiator *ecdsa.PrivateKey, redeemerAddr, token common.Address, expiry *big.Int, amount *big.Int, secretHash []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	var hash [32]byte
	copy(hash[:], secretHash)

	initiatorAddr := crypto.PubkeyToAddress(initiator.PublicKey)
	val, err := client.allowance(token, contract, initiatorAddr)
	if err != nil {
		return "", err
	}
	if val.Cmp(amount) < 0 {
		_, err := client.ApproveERC20(initiator, maxApproval, token, contract)
		if err != nil {
			return "", err
		}
	}

	transactor, err := client.GetTransactOpts(initiator)
	if err != nil {
		return "", err
	}
	initTx, err := instance.Initiate(transactor, redeemerAddr, expiry, amount, hash)
	if err != nil {
		return "", err
	}
	client.logger.Info("initiate swap", zap.String("txHash", initTx.Hash().Hex()))
	return initTx.Hash().Hex(), nil
}

func (client *client) RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, orderID [32]byte, secret []byte) (string, error) {
	instance, err := AtomicSwap.NewAtomicSwap(contract, client.provider)
	if err != nil {
		return "", err
	}
	tx, err := instance.Redeem(auth, orderID, secret)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
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

func (client *client) IsFinal(txHash string, minConf uint64) (bool, uint64, error) {
	tx, err := client.provider.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return false, 0, err
	}
	if tx.Status == 0 {
		return false, 0, nil
	}
	currentBlock, err := client.GetCurrentBlock()
	if err != nil {
		return false, 0, fmt.Errorf("error getting current block %v", err)
	}
	progress := int64(currentBlock) - tx.BlockNumber.Int64()
	if progress > 0 {
		return progress >= int64(minConf)-1, uint64(progress), nil
	}
	return progress+1 >= int64(minConf), 0, nil
}

func (client *client) allowance(tokenAddr common.Address, spender common.Address, owner common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20(tokenAddr, client.provider)
	if err != nil {
		return nil, err
	}
	return instance.Allowance(client.callOpts(), owner, spender)
}

func (client *client) callOpts() *bind.CallOpts {
	return &bind.CallOpts{
		Pending: true,
	}
}
