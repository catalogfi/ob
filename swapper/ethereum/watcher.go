package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/catalogfi/wbtc-garden/swapper"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type watcher struct {
	client           Client
	atomicSwapAddr   common.Address
	amount           *big.Int
	expiry           *big.Int
	minConfirmations *big.Int
	secretHash       []byte
	orderId          []byte
	lastCheckedBlock *big.Int
}

func NewWatcher(atomicSwapAddr common.Address, secretHash, orderId []byte, expiry, minConfirmations, amount *big.Int, client Client) (swapper.Watcher, error) {
	currentBlock, err := client.GetCurrentBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get the current block: %v", err)
	}

	// TODO: we only look at last 100 expiries from the current block, could potentially optimised
	lastCheckedBlock := new(big.Int).Sub(new(big.Int).SetUint64(currentBlock), new(big.Int).Mul(big.NewInt(100), expiry))
	if lastCheckedBlock.Cmp(big.NewInt(0)) < 0 {
		lastCheckedBlock = big.NewInt(0)
	}

	return &watcher{
		client:           client,
		atomicSwapAddr:   atomicSwapAddr,
		expiry:           expiry,
		amount:           amount,
		secretHash:       secretHash,
		minConfirmations: minConfirmations,
		lastCheckedBlock: lastCheckedBlock,
		orderId:          orderId,
	}, nil
}

func (watcher *watcher) Expired() (bool, error) {
	initiated, txHash, _, _ := watcher.IsInitiated()
	if !initiated {
		return false, nil
	}
	currentBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, err
	}
	height, _, err := watcher.Status(txHash)
	if err != nil {
		return false, err
	}
	if currentBlock > height+watcher.expiry.Uint64() {
		return true, nil
	} else {
		return false, nil
	}
}

func (watcher *watcher) Status(txHash string) (uint64, uint64, error) {
	return watcher.client.GetConfirmations(txHash)
}

func (watcher *watcher) IsDetected() (bool, string, string, error) {
	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, "", "", err
	}

	initiatedEvent := atomicSwapAbi.Events["Initiated"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		Addresses: []common.Address{
			watcher.atomicSwapAddr,
		},
		Topics: [][]common.Hash{{initiatedEvent.ID}, {common.BytesToHash(watcher.orderId)}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", "", err
	}

	if len(logs) == 0 {
		return false, "", "", fmt.Errorf("no logs found")
	}

	vLog := logs[0]
	values, err := atomicSwapAbi.Unpack("Initiated", vLog.Data)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to unpack Initiated event data: %v", err)
	}

	val, ok := values[1].(*big.Int)
	if !ok {
		return false, "", "", fmt.Errorf("unable to decode amount from Initiated event data")
	}

	if val.Cmp(watcher.amount) < 0 {
		return false, "", "", fmt.Errorf("initiated with lower than expected amount")
	}

	return true, vLog.TxHash.Hex(), val.String(), nil
}

func (watcher *watcher) IsInitiated() (bool, string, uint64, error) {
	fmt.Println("Checking if initiated")
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, "", 0, err
	}
	currentBlock := big.NewInt(int64(currBlock))
	// if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
	// 	currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	// }

	atomicSwapAbi, err := AtomicSwap.AtomicSwapMetaData.GetAbi()
	if err != nil {
		return false, "", 0, err
	}

	initiatedEvent := atomicSwapAbi.Events["Initiated"]
	query := ethereum.FilterQuery{
		FromBlock: watcher.lastCheckedBlock,
		ToBlock:   currentBlock,
		Addresses: []common.Address{
			watcher.atomicSwapAddr,
		},
		Topics: [][]common.Hash{{initiatedEvent.ID}, {common.BytesToHash(watcher.orderId)}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", 0, err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		// newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		// if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
		// 	watcher.lastCheckedBlock = currentBlock
		// }
		fmt.Println("No logs found")
		return false, "", 0, err
	}

	vLog := logs[0]

	isFinal, progress, err := watcher.client.IsFinal(vLog.TxHash.Hex(), watcher.minConfirmations.Uint64())
	if err != nil {
		return false, "", 0, err
	}

	if !isFinal {
		return false, "", progress, nil
	}

	return true, vLog.TxHash.Hex(), watcher.minConfirmations.Uint64(), nil
}

func (watcher *watcher) IsRedeemed() (bool, []byte, string, error) {
	currBlock, err := watcher.client.GetCurrentBlock()
	if err != nil {
		return false, nil, "", err
	}
	currentBlock := big.NewInt(int64(currBlock))
	// if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
	// 	currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	// }

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
		Topics: [][]common.Hash{{redeemedEvent.ID}, {common.BytesToHash(watcher.orderId)}, {common.BytesToHash(watcher.secretHash)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, nil, "", err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		// newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		// if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
		// 	watcher.lastCheckedBlock = currentBlock
		// }
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
	// if currentBlock.Int64() > watcher.lastCheckedBlock.Int64()+MaxQueryBlockRange {
	// 	currentBlock = big.NewInt(0).Add(watcher.lastCheckedBlock, big.NewInt(MaxQueryBlockRange))
	// }

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
		Topics: [][]common.Hash{{refundedEvent.ID}, {common.BytesToHash(watcher.orderId)}},
	}

	logs, err := watcher.client.GetProvider().FilterLogs(context.Background(), query)
	if err != nil {
		return false, "", err
	}

	if len(logs) == 0 {
		// Update the last checked block height
		// newLastCheckedBlock := big.NewInt(0).Sub(currentBlock, watcher.minConfirmations)
		// if newLastCheckedBlock.Cmp(watcher.lastCheckedBlock) == 1 {
		// 	watcher.lastCheckedBlock = currentBlock
		// }
		fmt.Println("No logs found")
		return false, "", err
	}
	return true, logs[0].TxHash.Hex(), nil
}
