package watcher

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/screener"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum"
	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/AtomicSwap"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type EthereumWatcher struct {
	chain          model.Chain
	netConfig      model.NetworkConfig
	startBlock     uint64
	interval       time.Duration
	store          Store
	atomicSwapAddr common.Address
	client         ethereum.Client
	ABI            *abi.ABI
	AtomicSwap     *AtomicSwap.AtomicSwap
	screener       screener.Screener
	logger         *zap.Logger
	ignoreOrders   map[string]bool
	blockSpan      uint64
}

type Swap struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}

var retryCount = 5

func NewEthereumWatchers(store Store, config model.Config, screener screener.Screener, logger *zap.Logger) ([]*EthereumWatcher, error) {
	var watchers []*EthereumWatcher
	for chain, netConfig := range config.Network {
		for asset, token := range netConfig.Assets {
			swapAddr := common.HexToAddress(asset.SecondaryID())
			watcher, err := NewEthereumWatcher(store, chain, netConfig, swapAddr, token.StartBlock, uint64(netConfig.EventWindow), screener, logger)
			if err != nil {
				return nil, err
			}
			watchers = append(watchers, watcher)
		}
	}
	return watchers, nil
}

func NewEthereumWatcher(store Store, chain model.Chain, config model.NetworkConfig, address common.Address, startBlock uint64, blockSpan uint64, screener screener.Screener, logger *zap.Logger) (*EthereumWatcher, error) {
	ethClient, err := ethereum.NewClient(logger, config.RPC["ethrpc"])
	if err != nil {
		return nil, fmt.Errorf("failed to load client: %v", err)
	}
	atomicSwap, _ := AtomicSwap.NewAtomicSwap(address, ethClient.GetProvider())
	atomicSwapAbi, _ := AtomicSwap.AtomicSwapMetaData.GetAbi()
	return &EthereumWatcher{
		chain:          chain,
		netConfig:      config,
		interval:       5 * time.Second,
		store:          store,
		atomicSwapAddr: address,
		client:         ethClient,
		startBlock:     startBlock,
		AtomicSwap:     atomicSwap,
		screener:       screener,
		ABI:            atomicSwapAbi,
		logger:         logger,
		ignoreOrders:   make(map[string]bool),
		blockSpan:      blockSpan,
	}, nil
}

func (w *EthereumWatcher) Watch() {
	eventIds := [][]common.Hash{{
		w.ABI.Events["Initiated"].ID,
		w.ABI.Events["Redeemed"].ID,
		w.ABI.Events["Refunded"].ID,
	}}
	for {
		currentBlock, err := w.client.GetCurrentBlock()
		if err != nil {
			w.logger.Error("failed to get current block number", zap.Error(err))
			continue
		}

		var logs []types.Log
		toBlock := w.startBlock
		fetchedAll := false
		for !fetchedAll {
			fromBlock := toBlock
			toBlock += w.blockSpan
			if toBlock >= currentBlock {
				toBlock = currentBlock
				fetchedAll = true
			}
			logsSlice, err := w.client.GetLogs(w.atomicSwapAddr, fromBlock, toBlock, eventIds)
			if err != nil {
				w.logger.Error("failed to get logs", zap.Error(err), zap.Any("chain", w.chain), zap.Any("toBlock", toBlock), zap.Any("fromBlock", fromBlock))
				fetchedAll = false
				toBlock = fromBlock
				continue
			}
			logs = append(logs, logsSlice...)
		}

		werr := HandleEVMLogs(eventIds, logs, w.store, w.screener, w.AtomicSwap, w.logger)
		if werr != nil {
			if werr.isUnrecoverable() {
				w.logger.Error("an unrecoverable error occurred while handling evm logs, shutting down", zap.Error(werr.Err), zap.Any("chain", w.chain))
				return
			}
			w.logger.Error("failed to handle evm logs", zap.Error(werr.Err))
		}
		err = UpdateEVMConfirmations(w.store, w.chain, currentBlock)
		if err != nil {
			w.logger.Error("failed to update confirmations", zap.Error(err))
		}

		w.startBlock = currentBlock
		time.Sleep(w.interval)
	}
}

func HandleEVMLogs(eventIds [][]common.Hash, logs []types.Log, store Store, screener screener.Screener, contract *AtomicSwap.AtomicSwap, logger *zap.Logger) *WatcherError {
	for _, log := range logs {
		switch log.Topics[0] {
		case eventIds[0][0]:
			cSwap, err := RetryWithReturnValue(func() (Swap, error) {
				return contract.AtomicSwapOrders(nil, log.Topics[1])
			}, retryCount)
			if err != nil {
				return NewWatcherError(fmt.Errorf("failed to get swap order: %s", err))
			}
			handler := func() *WatcherError {
				return HandleEVMInitiate(log, store, cSwap, screener)
			}
			wErr := RetryOnWatcherError(handler, retryCount)
			if wErr != nil && wErr.isUnrecoverable() {
				return wErr
			}
		case eventIds[0][1]:
			handler := func() *WatcherError {
				return HandleEVMRedeem(store, log)
			}
			err := RetryOnWatcherError(handler, retryCount)
			if err != nil && err.isUnrecoverable() {
				return err
			}
		case eventIds[0][2]:
			handler := func() *WatcherError {
				return HandleEVMRefund(store, log)
			}
			err := RetryOnWatcherError(handler, retryCount)
			if err != nil && err.isUnrecoverable() {
				return err
			}
		}
	}
	return nil
}

// update confirmation status of unconfirmed initiates
func UpdateEVMConfirmations(store Store, chain model.Chain, currentBlock uint64) error {
	swaps, err := store.GetActiveSwaps(chain)
	if err != nil {
		return err
	}
	for _, swap := range swaps {
		timelock, err := strconv.ParseUint(swap.Timelock, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to ParseUint timelock: %s", err)
		}
		if swap.Status == model.Detected && currentBlock > swap.InitiateBlockNumber {
			confirmations := currentBlock - swap.InitiateBlockNumber + 1
			if confirmations != swap.CurrentConfirmations {
				swap.CurrentConfirmations = confirmations
				if confirmations >= swap.MinimumConfirmations {
					swap.CurrentConfirmations = swap.MinimumConfirmations
					swap.Status = model.Initiated
				}
				if err := store.UpdateSwap(&swap); err != nil {
					return err
				}
			}
		}
		if swap.Status == model.Initiated && currentBlock > swap.InitiateBlockNumber+uint64(timelock) {
			swap.Status = model.Expired
			if err := store.UpdateSwap(&swap); err != nil {
				return err
			}
		}
	}
	return nil
}

func HandleEVMInitiate(log types.Log, store Store, cSwap Swap, screener screener.Screener) *WatcherError {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])
	if err != nil {
		if err.Error() != "record not found" {
			// we just ignore the error and move on to the next log
			return nil
		}
		return NewWatcherError(
			fmt.Errorf("failed to get swap by ocid: %s", err),
			WithRecoverable,
		)
	}
	if swap.InitiateTxHash != "" {
		return nil
	}

	isBlacklisted, err := screener.IsBlacklisted(map[string]model.Chain{cSwap.Initiator.Hex(): swap.Chain})
	if err != nil {
		return NewWatcherError(fmt.Errorf("failed to check if address is blacklisted, %s", cSwap.Initiator.Hex()), WithRecoverable)
	}
	if isBlacklisted {
		return NewWatcherError(fmt.Errorf("address is blacklisted, %s", cSwap.Initiator.Hex()), WithIgnore)
	}

	amount, ok := new(big.Int).SetString(swap.Amount, 10)
	if !ok {
		return NewWatcherError(fmt.Errorf("invalid amount: %s", swap.Amount), WithIgnore)
	}
	if cSwap.Amount.Cmp(amount) < 0 {
		return NewWatcherError(fmt.Errorf("insufficient amount: %s", swap.Amount), WithIgnore)
	}
	expiry, ok := new(big.Int).SetString(swap.Timelock, 10)
	if !ok {
		return NewWatcherError(fmt.Errorf("failed to decode timelock: %s", err), WithIgnore)
	}

	if cSwap.Expiry.Cmp(expiry) != 0 {
		return NewWatcherError(fmt.Errorf("incorrect expiry: %s", swap.Amount), WithIgnore)
	}

	if strings.ToLower(cSwap.Redeemer.String()) != strings.ToLower(swap.RedeemerAddress) {
		return NewWatcherError(fmt.Errorf("incorrect redeemer: %s", swap.RedeemerAddress), WithIgnore)
	}

	swap.InitiateTxHash = log.TxHash.String()
	swap.InitiateBlockNumber = log.BlockNumber
	swap.Status = model.Detected

	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewWatcherError(fmt.Errorf("failed to update swap: %s", err), WithRecoverable)
	}
	return nil
}

func HandleEVMRedeem(store Store, log types.Log) *WatcherError {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])

	if err != nil {
		if err.Error() != "record not found" {
			// we just ignore the error and move on to the next log
			return nil
		}
		return NewWatcherError(
			fmt.Errorf("failed to get swap by ocid: %s", err),
			WithRecoverable,
		)
	}

	if swap.RedeemTxHash != "" {
		return nil
	}
	if len(log.Data) < 64 {
		return NewWatcherError(
			fmt.Errorf("invalid log data: %x", log.Data),
			WithIgnore,
		)
	}
	swap.Secret = hex.EncodeToString(log.Data[64:])
	swap.RedeemTxHash = log.TxHash.Hex()
	swap.Status = model.Redeemed
	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewWatcherError(
			fmt.Errorf("failed to update swap: %s", err),
			WithRecoverable,
		)
	}
	return nil
}

func HandleEVMRefund(store Store, log types.Log) *WatcherError {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])
	if err != nil {
		if err.Error() != "record not found" {
			// we just ignore the error and move on to the next log
			return nil
		}
		return NewWatcherError(
			fmt.Errorf("failed to get swap by ocid: %s", err),
			WithRecoverable,
		)
	}
	if swap.RefundTxHash != "" {
		return nil
	}
	swap.RefundTxHash = log.TxHash.String()
	swap.Status = model.Refunded
	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewWatcherError(
			fmt.Errorf("failed to update swap: %s", err),
			WithRecoverable,
		)
	}
	return nil
}

func RetryOnWatcherError(f func() *WatcherError, retries int) *WatcherError {
	var err *WatcherError
	for i := 0; i < retries; i++ {
		err = f()
		if err == nil {
			return nil
		}
		if err.ignore || !err.Recoverable {
			return err
		}
		time.Sleep(time.Duration(retries) * 1000 * time.Millisecond)
	}
	err.Recoverable = false
	return err
}

func RetryWithReturnValue[T any](f func() (T, error), retries int) (T, error) {
	var err error
	var nilData T
	for i := 0; i < retries; i++ {
		data, err := f()
		if err == nil {
			return data, nil
		}
		time.Sleep(time.Duration(retries) * 1000 * time.Millisecond)
	}
	return nilData, err
}
