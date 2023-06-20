package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/susruth/wbtc-garden/swapper/ethereum/typings/ERC20"
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
	redeemerAddress  common.Address
	refunderAddress  common.Address
	deployerAddress  common.Address
	lastCheckedBlock *big.Int
	expiryBlock      *big.Int
	contractAddr     common.Address
	tokenAddr        common.Address
	amount           *big.Int
	secretHash       []byte
	client           Client
}

func GetAmount(client Client, deployerAddr, tokenAddr, redeemerAddr, initiatorAddr common.Address, secretHash []byte, expiryBlock *big.Int) (uint64, error) {
	contractAddr, err := GetAddress(client, deployerAddr, redeemerAddr, initiatorAddr, secretHash, expiryBlock)
	if err != nil {
		return 0, err
	}
	balance, err := client.GetERC20Balance(tokenAddr, contractAddr)
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

// var deployerAddr = common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2")

func NewInitiatorSwap(initiator *ecdsa.PrivateKey, redeemerAddr, deployerAddr, tokenAddr common.Address, secretHash []byte, expiryBlock *big.Int, amount *big.Int, client Client) (swapper.InitiatorSwap, error) {
	fmt.Println(expiryBlock.Text(10), hex.EncodeToString(secretHash))

	initiatorAddr := client.GetPublicAddress(initiator)

	fmt.Println(redeemerAddr, initiatorAddr, secretHash, expiryBlock)
	contractAddr, err := GetAddress(client, deployerAddr, redeemerAddr, initiatorAddr, secretHash, expiryBlock)
	if err != nil {
		return &initiatorSwap{}, err
	}
	fmt.Println("Contract Address : ", contractAddr.String())

	latestCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(6000))
	return &initiatorSwap{initiator: initiator, initiatorAddr: initiatorAddr, expiryBlock: expiryBlock, contractAddr: contractAddr, client: client, amount: amount, tokenAddr: tokenAddr, redeemerAddr: redeemerAddr, lastCheckedBlock: latestCheckedBlock}, nil
}

func (initiatorSwap *initiatorSwap) Initiate() (string, error) {
	defer fmt.Printf("Done Initiate on contract : %s : token : %s \n", initiatorSwap.contractAddr, initiatorSwap.tokenAddr)

	txHash, err := initiatorSwap.client.TransferERC20(initiatorSwap.initiator, initiatorSwap.amount, initiatorSwap.tokenAddr, initiatorSwap.contractAddr)
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

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, nil, "", err
	}

	redeemedEvent := atomicSwapAbi.Events["Redeemed"]
	query := ethereum.FilterQuery{
		FromBlock: initiatorSwap.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			initiatorSwap.contractAddr,
		},
		Topics: [][]common.Hash{{redeemedEvent.ID}},
	}

	logs, err := initiatorSwap.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, nil, "", err
	}

	if len(logs) == 0 {
		fmt.Println("No logs found")
		return false, nil, "", err
	}

	vLog := logs[0]

	val, err := redeemedEvent.Inputs.Unpack(vLog.Data)
	if err != nil {
		return false, nil, "", err
	}

	return true, []byte(val[0].(string)), vLog.TxHash.Hex(), nil
}

func (initiatorSwap *initiatorSwap) Refund() (string, error) {
	defer fmt.Println("Done refund")
	tx, err := initiatorSwap.client.RefundAtomicSwap(initiatorSwap.contractAddr, initiatorSwap.client.GetTransactOpts(initiatorSwap.initiator), initiatorSwap.tokenAddr)
	if err != nil {
		return "", err
	}
	return tx, nil
}

func NewRedeemerSwap(redeemer *ecdsa.PrivateKey, initiatorAddr, deployerAddr, tokenAddr common.Address, secretHash []byte, expiryBlock *big.Int, amount *big.Int, client Client) (swapper.RedeemerSwap, error) {
	fmt.Println(expiryBlock.Text(10), hex.EncodeToString(secretHash))

	redeemerAddress := crypto.PubkeyToAddress(redeemer.PublicKey)

	fmt.Println(redeemerAddress, initiatorAddr, secretHash, expiryBlock)
	contractAddr, err := GetAddress(client, deployerAddr, redeemerAddress, initiatorAddr, secretHash, expiryBlock)
	if err != nil {
		return &redeemerSwap{}, err
	}
	fmt.Println("Contract Address : ", contractAddr.String())

	lastCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(6000))
	return &redeemerSwap{
		redeemer:         redeemer,
		redeemerAddress:  redeemerAddress,
		refunderAddress:  initiatorAddr,
		lastCheckedBlock: lastCheckedBlock,
		expiryBlock:      expiryBlock,
		contractAddr:     contractAddr,
		tokenAddr:        tokenAddr,
		deployerAddress:  deployerAddr,
		amount:           amount,
		client:           client,
		secretHash:       secretHash,
	}, nil
}

func (redeemerSwap *redeemerSwap) Redeem(secret []byte) (string, error) {
	defer fmt.Println("Done redeem")
	data, err := redeemerSwap.client.GetProvider().CodeAt(context.Background(), redeemerSwap.contractAddr, nil)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		txHash, err := Deploy(redeemerSwap.deployerAddress, redeemerSwap.client.GetTransactOpts(redeemerSwap.redeemer), redeemerSwap.client.GetProvider(), redeemerSwap.redeemerAddress, redeemerSwap.refunderAddress, redeemerSwap.secretHash, redeemerSwap.expiryBlock)
		if err != nil {
			return "", err
		}
		fmt.Println("Deployed contract at ", txHash)
	}
	return redeemerSwap.client.RedeemAtomicSwap(redeemerSwap.contractAddr, redeemerSwap.client.GetTransactOpts(redeemerSwap.redeemer), redeemerSwap.tokenAddr, secret)
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

	erc20Abi, err := ERC20.ERC20MetaData.GetAbi()
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

	amount := big.NewInt(0)
	for _, vLog := range logs {
		inputs, err := transferEvent.Inputs.Unpack(vLog.Data)
		if err != nil {
			return false, "", err
		}
		amount.Add(amount, inputs[0].(*big.Int))
		if amount.Cmp(redeemerSwap.amount) >= 0 {
			final, err := redeemerSwap.client.IsFinal(vLog.TxHash.Hex())
			if err != nil {
				return false, "", err
			}
			if !final {
				return false, "", nil
			}
			return true, vLog.TxHash.Hex(), nil
		}
	}
	return false, "", err
}

type watcher struct {
	client           Client
	tokenAddr        common.Address
	contractAddr     common.Address
	lastCheckedBlock *big.Int
	amount           *big.Int
}

func NewWatcher(initiator, redeemerAddr, deployerAddr, tokenAddr common.Address, secretHash []byte, expiryBlock *big.Int, amount *big.Int, client Client) (swapper.Watcher, error) {
	contractAddr, err := GetAddress(client, deployerAddr, redeemerAddr, initiator, secretHash, expiryBlock)
	if err != nil {
		return &watcher{}, err
	}
	latestCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(6000))
	return &watcher{
		client:           client,
		tokenAddr:        tokenAddr,
		contractAddr:     contractAddr,
		lastCheckedBlock: latestCheckedBlock,
		amount:           amount,
	}, nil
}

func (watcher *watcher) IsInitiated() (bool, string, error) {
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))

	erc20Abi, err := ERC20.ERC20MetaData.GetAbi()
	if err != nil {
		return false, "", err
	}
	transferEvent := erc20Abi.Events["Transfer"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.tokenAddr,
		},
		Topics: [][]common.Hash{{transferEvent.ID}, {}, {watcher.contractAddr.Hash()}},
	}
	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", err
	}

	if len(logs) == 0 {
		watcher.lastCheckedBlock = currentBlock
		return false, "", err
	}

	amount := big.NewInt(0)
	for _, vLog := range logs {
		inputs, err := transferEvent.Inputs.Unpack(vLog.Data)
		if err != nil {
			return false, "", err
		}
		amount.Add(amount, inputs[0].(*big.Int))
		if amount.Cmp(watcher.amount) >= 0 {
			return true, vLog.TxHash.Hex(), nil
		}
	}
	return false, "", err
}

func (watcher *watcher) IsRedeemed() (bool, []byte, string, error) {
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, nil, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, nil, "", err
	}

	redeemedEvent := atomicSwapAbi.Events["Redeemed"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.contractAddr,
		},
		Topics: [][]common.Hash{{redeemedEvent.ID}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, nil, "", err
	}

	if len(logs) == 0 {
		fmt.Println("No logs found")
		return false, nil, "", err
	}

	vLog := logs[0]

	val, err := redeemedEvent.Inputs.Unpack(vLog.Data)
	if err != nil {
		return false, nil, "", err
	}

	return true, []byte(val[0].(string)), vLog.TxHash.Hex(), nil
}
