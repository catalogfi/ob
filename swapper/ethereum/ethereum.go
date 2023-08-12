package ethereum

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/catalogfi/wbtc-garden/swapper"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

const MaxQueryBlockRange = 500

type initiatorSwap struct {
	initiator        *ecdsa.PrivateKey
	initiatorAddr    common.Address
	redeemerAddr     common.Address
	lastCheckedBlock *big.Int
	expiryBlock      *big.Int
	atomicSwapAddr   common.Address
	secretHash       []byte
	client           Client
	amount           *big.Int
	tokenAddr        common.Address
	watcher          swapper.Watcher
}
type redeemerSwap struct {
	redeemer         *ecdsa.PrivateKey
	lastCheckedBlock *big.Int
	expiryBlock      *big.Int
	atomicSwapAddr   common.Address
	tokenAddr        common.Address
	amount           *big.Int
	secretHash       []byte
	client           Client
	watcher          swapper.Watcher
}

func GetExpiry(client Client, goingFirst bool) (*big.Int, error) {
	blockNumber, err := client.GetProvider().BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	if goingFirst {
		return new(big.Int).Add(new(big.Int).SetUint64(blockNumber), big.NewInt(11520)), nil
	}
	return new(big.Int).Add(new(big.Int).SetUint64(blockNumber), big.NewInt(5760)), nil
}

func NewInitiatorSwap(initiator *ecdsa.PrivateKey, redeemerAddr, atomicSwapAddr common.Address, secretHash []byte, expiryBlock, minConfirmations, amount *big.Int, client Client) (swapper.InitiatorSwap, error) {

	initiatorAddr := client.GetPublicAddress(initiator)

	latestCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(12000))
	if latestCheckedBlock.Cmp(big.NewInt(0)) == -1 {
		latestCheckedBlock = big.NewInt(0)
	}

	watcher, err := NewWatcher(atomicSwapAddr, secretHash, expiryBlock, minConfirmations, amount, client)
	if err != nil {
		return &initiatorSwap{}, err
	}
	tokenAddr, err := client.GetTokenAddress(atomicSwapAddr)
	if err != nil {
		return &initiatorSwap{}, err
	}
	return &initiatorSwap{
		initiator: initiator, watcher: watcher,
		initiatorAddr:  initiatorAddr,
		expiryBlock:    expiryBlock,
		atomicSwapAddr: atomicSwapAddr,
		client:         client, amount: amount,
		tokenAddr:        tokenAddr,
		redeemerAddr:     redeemerAddr,
		lastCheckedBlock: latestCheckedBlock,
		secretHash:       secretHash}, nil
}

func (initiatorSwap *initiatorSwap) Initiate() (txHash string, err error) {
	defer func() {
		fmt.Printf("Done Initiate on contract : %s : token : %s : err : %v \n", initiatorSwap.atomicSwapAddr, initiatorSwap.tokenAddr, err)
	}()
	txHash, err = initiatorSwap.client.InitiateAtomicSwap(initiatorSwap.atomicSwapAddr, initiatorSwap.initiator, initiatorSwap.redeemerAddr, initiatorSwap.tokenAddr, initiatorSwap.expiryBlock, initiatorSwap.amount, initiatorSwap.secretHash)
	return
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
		if err != nil {
			fmt.Println("failed to check redeemed status", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if redeemed {
			return secret, txHash, err
		}
		time.Sleep(5 * time.Second)
	}
}

func (initiatorSwap *initiatorSwap) IsRedeemed() (bool, []byte, string, error) {
	return initiatorSwap.watcher.IsRedeemed()
}

func (initiatorSwap *initiatorSwap) Refund() (string, error) {
	defer fmt.Println("Done refund")
	tx, err := initiatorSwap.client.RefundAtomicSwap(initiatorSwap.atomicSwapAddr, initiatorSwap.client.GetTransactOpts(initiatorSwap.initiator), initiatorSwap.tokenAddr, initiatorSwap.secretHash)
	if err != nil {
		return "", err
	}
	return tx, nil
}

func NewRedeemerSwap(redeemer *ecdsa.PrivateKey, initiatorAddr, atomicSwapAddr common.Address, secretHash []byte, expiryBlock, amount, minConfirmations *big.Int, client Client) (swapper.RedeemerSwap, error) {
	watcher, err := NewWatcher(atomicSwapAddr, secretHash, expiryBlock, minConfirmations, amount, client)
	if err != nil {
		return &redeemerSwap{}, err
	}

	tokenAddr, err := client.GetTokenAddress(atomicSwapAddr)
	if err != nil {
		return &redeemerSwap{}, err
	}

	lastCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(12000))
	return &redeemerSwap{
		redeemer:         redeemer,
		watcher:          watcher,
		lastCheckedBlock: lastCheckedBlock,
		expiryBlock:      expiryBlock,
		atomicSwapAddr:   atomicSwapAddr,
		tokenAddr:        tokenAddr,
		amount:           amount,
		client:           client,
		secretHash:       secretHash,
	}, nil
}

func (redeemerSwap *redeemerSwap) Redeem(secret []byte) (string, error) {
	defer fmt.Println("Done redeem")
	fmt.Println("redeeming...")

	return redeemerSwap.client.RedeemAtomicSwap(redeemerSwap.atomicSwapAddr, redeemerSwap.client.GetTransactOpts(redeemerSwap.redeemer), redeemerSwap.tokenAddr, secret)
}

func (redeemerSwap *redeemerSwap) IsInitiated() (bool, []string, error) {
	return redeemerSwap.watcher.IsInitiated()
}

func (redeemerSwap *redeemerSwap) WaitForInitiate() ([]string, error) {
	defer fmt.Println("Done WaitForInitiate")
	for {
		initiated, txHash, err := redeemerSwap.IsInitiated()
		if initiated {
			fmt.Printf("Initiation Found on contract : %s : token : %s \n", redeemerSwap.atomicSwapAddr, redeemerSwap.tokenAddr)
			return txHash, nil
		}
		if err != nil {
			fmt.Println("failed to check initiated status", err)
		}
		time.Sleep(5 * time.Second)
	}
}

type watcher struct {
	client           Client
	atomicSwapAddr   common.Address
	lastCheckedBlock *big.Int
	amount           *big.Int
	expiryBlock      *big.Int
	secretHash       []byte
	minConfirmations *big.Int
}

func NewWatcher(atomicSwapAddr common.Address, secretHash []byte, expiryBlock, minConfirmations, amount *big.Int, client Client) (swapper.Watcher, error) {
	latestCheckedBlock := new(big.Int).Sub(expiryBlock, big.NewInt(12000))
	if latestCheckedBlock.Cmp(big.NewInt(0)) == -1 {
		latestCheckedBlock = big.NewInt(0)
	}
	return &watcher{
		client:           client,
		atomicSwapAddr:   atomicSwapAddr,
		lastCheckedBlock: latestCheckedBlock,
		expiryBlock:      expiryBlock,
		amount:           amount,
		secretHash:       secretHash,
		minConfirmations: minConfirmations,
	}, nil
}

func (watcher *watcher) Expired() (bool, error) {
	currentBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, err
	}
	if currentBlock > watcher.expiryBlock.Uint64() {
		return true, nil
	} else {
		return false, nil
	}
}

func (watcher *watcher) IsInitiated() (bool, []string, error) {
	fmt.Println("Checking if initiated")
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, []string{}, err
	}
	currentBlock := big.NewInt(int64(currBlock))
	if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
		currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	}

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, []string{}, err
	}

	initiatedEvent := atomicSwapAbi.Events["Initiated"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.atomicSwapAddr,
		},
		Topics: [][]common.Hash{{initiatedEvent.ID}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, []string{}, err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
			watcher.lastCheckedBlock = currentBlock
		}
		fmt.Println("No logs found")
		return false, []string{}, err
	}

	vLog := logs[0]

	isFinal, err := watcher.client.IsFinal(vLog.TxHash.Hex(), watcher.minConfirmations.Uint64())
	if err != nil {
		return false, []string{}, err
	}

	if !isFinal {
		return false, []string{}, fmt.Errorf("transaction not finalized yet")
	}

	return true, []string{vLog.TxHash.Hex()}, nil
}

func (watcher *watcher) IsRedeemed() (bool, []byte, string, error) {
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, nil, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))
	if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
		currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	}

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, nil, "", err
	}

	redeemedEvent := atomicSwapAbi.Events["Redeemed"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.atomicSwapAddr,
		},
		Topics: [][]common.Hash{{redeemedEvent.ID}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, nil, "", err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
			watcher.lastCheckedBlock = currentBlock
		}
		fmt.Println("No logs found")
		return false, nil, "", err
	}

	vLog := logs[0]

	val, err := redeemedEvent.Inputs.Unpack(vLog.Data)
	if err != nil {
		return false, nil, "", err
	}

	return true, []byte(val[0].([]uint8)), vLog.TxHash.Hex(), nil
}

func (watcher *watcher) IsRefunded() (bool, string, error) {
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))
	if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
		currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	}

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, "", err
	}

	refundedEvent := atomicSwapAbi.Events["Refunded"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.atomicSwapAddr,
		},
		Topics: [][]common.Hash{{refundedEvent.ID}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
			watcher.lastCheckedBlock = currentBlock
		}
		fmt.Println("No logs found")
		return false, "", err
	}
	return true, logs[0].TxHash.Hex(), nil
}
