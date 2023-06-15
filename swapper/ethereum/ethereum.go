package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden-server/swapper"
)

type initiatorSwap struct {
	initiator        *ecdsa.PrivateKey
	initiatorAddr    common.Address
	redeemerAddr     common.Address
	lastCheckedBlock *big.Int
	expiryBlock      *big.Int
	contractAddr     common.Address
	client           Client
	amount           *big.Int
	tokenAddr        common.Address
}
type redeemerSwap struct {
	redeemer         *ecdsa.PrivateKey
	lastCheckedBlock *big.Int
	expiryBlock      *big.Int
	contractAddr     common.Address
	tokenAddr        common.Address
	amount           *big.Int
	client           Client
}

func NewInitiatorSwap(initiator *ecdsa.PrivateKey, redeemerAddr, tokenAddr common.Address, secretHash []byte, expiryBlock *big.Int, amount *big.Int, client Client) (swapper.InitiatorSwap, error) {
	deployerAddr := common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2")
	initiatorAddr := client.GetPublicAddress(initiator)
	contractAddr, err := GetAddress(deployerAddr, client.GetCallOpts(), client.GetProvider(), redeemerAddr, initiatorAddr, secretHash, expiryBlock)
	if err != nil {
		return &initiatorSwap{}, err
	}
	txHash, deploymentErr := Deploy(deployerAddr, client.GetTransactOpts(initiator), client.GetProvider(), redeemerAddr, initiatorAddr, secretHash, expiryBlock)
	if deploymentErr != nil {
		panic(deploymentErr)
	}
	txReceipt, err := client.GetProvider().TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return &initiatorSwap{}, err
	}
	txBlock := txReceipt.BlockNumber
	return &initiatorSwap{initiator: initiator, initiatorAddr: initiatorAddr, expiryBlock: expiryBlock, contractAddr: contractAddr, client: client, amount: amount, tokenAddr: tokenAddr, redeemerAddr: redeemerAddr, lastCheckedBlock: txBlock}, nil
}

func (initiatorSwap *initiatorSwap) Initiate() (string, error) {
	defer fmt.Printf("Done Initiate on contract : %s : token : %s \n", initiatorSwap.contractAddr, initiatorSwap.tokenAddr)
	txHash, err := initiatorSwap.client.TransferERC20(initiatorSwap.initiator, initiatorSwap.amount, initiatorSwap.tokenAddr, initiatorSwap.contractAddr, initiatorSwap.client.GetTransactOpts(initiatorSwap.initiator))
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (initiatorSwap *initiatorSwap) Expired() (bool, error) {
	currentBlock, err := initiatorSwap.client.GetCurrentBlock()
	if err != nil {
		return false, err
	}

	if currentBlock > initiatorSwap.expiryBlock.Uint64() {
		return true, nil
	} else {
		return false, nil
	}
}

func (initiatorSwap *initiatorSwap) WaitForRedeem() ([]byte, string, error) {
	for {
		redeemed, secret, txHash, err := initiatorSwap.IsRedeemed()
		if redeemed {
			return secret, txHash, err
		}
		if err != nil {
			fmt.Println("failed to check redeemed status", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (initiatorSwap *initiatorSwap) IsRedeemed() (bool, []byte, string, error) {
	currBlock, err := initiatorSwap.client.GetCurrentBlock()
	if err != nil {
		return false, nil, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))

	query := ethereum.FilterQuery{
		FromBlock: initiatorSwap.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			initiatorSwap.contractAddr,
		},
	}

	logs, err := initiatorSwap.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, nil, "", err
	}

	if len(logs) == 0 {
		initiatorSwap.lastCheckedBlock = currentBlock
		return false, nil, "", nil
	}

	AtomicSwapAbi, _ := abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"_secret","type":"bytes"}],"name":"redeemed","type":"event"}]`))
	vLog := logs[0]
	event, err := AtomicSwapAbi.Unpack("redeemed", vLog.Data)
	if err != nil {
		return false, nil, "", err
	}
	return true, event[0].([]byte), vLog.TxHash.Hex(), nil
}

func (initiatorSwap *initiatorSwap) Refund() (string, error) {
	defer fmt.Println("Done refund")
	tx, err := initiatorSwap.client.ExecuteAtomicSwap(initiatorSwap.contractAddr, initiatorSwap.client.GetTransactOpts(initiatorSwap.initiator), initiatorSwap.tokenAddr, []byte{})
	if err != nil {
		return "", err
	}
	return tx, nil
}

func NewRedeemerSwap(redeemer *ecdsa.PrivateKey, initiatorAddr, tokenAddr common.Address, secretHash []byte, expiryBlock *big.Int, amount *big.Int, client Client) (swapper.RedeemerSwap, error) {
	deployerAddress := common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2")
	redeemerAddress := crypto.PubkeyToAddress(redeemer.PublicKey)
	contractAddr, err := GetAddress(deployerAddress, nil, client.GetProvider(), redeemerAddress, initiatorAddr, secretHash, expiryBlock)
	if err != nil {
		return &redeemerSwap{}, err
	}
	lastCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(150000))
	return &redeemerSwap{contractAddr: contractAddr, tokenAddr: tokenAddr, client: client, redeemer: redeemer, amount: amount, lastCheckedBlock: lastCheckedBlock, expiryBlock: expiryBlock}, nil
}

func (redeemerSwap *redeemerSwap) Redeem(secret []byte) (string, error) {
	defer fmt.Println("Done redeem")
	return redeemerSwap.client.ExecuteAtomicSwap(redeemerSwap.contractAddr, redeemerSwap.client.GetTransactOpts(redeemerSwap.redeemer), redeemerSwap.tokenAddr, secret)
}

func (redeemerSwap *redeemerSwap) WaitForInitiate() (string, error) {
	defer fmt.Println("Done WaitForInitiate")
	for {
		initiated, txHash, err := redeemerSwap.IsInitiated()
		if initiated {
			return txHash, nil
		}
		if err != nil {
			fmt.Println("failed to check initiated status", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (redeemerSwap *redeemerSwap) IsInitiated() (bool, string, error) {
	currBlock, err := redeemerSwap.client.GetCurrentBlock()
	if err != nil {
		return false, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))

	erc20Abi, err := abi.JSON(strings.NewReader("{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}"))
	if err != nil {
		return false, "", err
	}

	transferEvent := erc20Abi.Events["Transfer"]

	query := ethereum.FilterQuery{
		FromBlock: redeemerSwap.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			redeemerSwap.tokenAddr,
		},
		Topics: [][]common.Hash{{transferEvent.ID}, {}, {redeemerSwap.contractAddr.Hash()}},
	}
	logs, err := redeemerSwap.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", err
	}

	if len(logs) == 0 {
		redeemerSwap.lastCheckedBlock = currentBlock
		return false, "", err
	}

	vLog := logs[0]
	return true, vLog.TxHash.Hex(), nil
}
