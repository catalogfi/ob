package watcher

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/catalogfi/orderbook/model"
	"github.com/catalogfi/orderbook/screener"
	"github.com/catalogfi/orderbook/swapper/ethereum"
	"github.com/catalogfi/orderbook/swapper/ethereum/typings/AtomicSwap"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EthereumL2Watcher struct {
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

func NewEthereumL2Watchers(store Store, config model.Config, screener screener.Screener, logger *zap.Logger) ([]*EthereumL2Watcher, error) {
	var watchers []*EthereumL2Watcher
	for chain, netConfig := range config.Network {
		for asset, token := range netConfig.Assets {
			swapAddr := common.HexToAddress(asset.SecondaryID())
			watcher, err := NewEthereumL2Watcher(store, chain, netConfig, swapAddr, token.StartBlock, uint64(netConfig.EventWindow), screener, logger)
			if err != nil {
				return nil, err
			}
			watchers = append(watchers, watcher)
		}
	}
	return watchers, nil
}

func NewEthereumL2Watcher(store Store, chain model.Chain, config model.NetworkConfig, address common.Address, startBlock uint64, blockSpan uint64, screener screener.Screener, logger *zap.Logger) (*EthereumL2Watcher, error) {
	ethClient, err := ethereum.NewClient(logger, config.RPC["ethrpc"])
	if err != nil {
		return nil, fmt.Errorf("failed to load client: %v", err)
	}
	atomicSwap, _ := AtomicSwap.NewAtomicSwap(address, ethClient.GetProvider())
	atomicSwapAbi, _ := AtomicSwap.AtomicSwapMetaData.GetAbi()
	return &EthereumL2Watcher{
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

func (w *EthereumL2Watcher) Watch() {
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
		currentL1Block, err := w.client.GetL1CurrentBlock()
		if err != nil {
			w.logger.Error("failed to get current block number", zap.Error(err))
			continue
		}
		if w.startBlock == currentBlock {
			time.Sleep(1 * time.Second)
			continue
		} else if w.startBlock > currentBlock {
			//this might happen because of a reorg or rpc is giving incorrect block number
			//So when this happens, just process it again.
			w.logger.Error("start block is greater than current block", zap.Uint64("startBlock", w.startBlock), zap.Uint64("currentBlock", currentBlock))
			w.startBlock = currentBlock
		}
		fromBlock := w.startBlock
		toBlock := currentBlock
		logsSlice, err := w.client.GetLogs(w.atomicSwapAddr, fromBlock, toBlock, eventIds, w.blockSpan)
		if err != nil {
			w.logger.Error("failed to get logs", zap.Error(err))
			continue
		}
		werr := w.HandleEVML2Logs(eventIds, logsSlice, w.store, w.screener, w.AtomicSwap, w.logger)
		if werr != nil {
			var nonrecoverable *NonRecoverableError
			if errors.As(werr, &nonrecoverable) {
				w.logger.Error("an unrecoverable error occurred while handling EVML2 logs, shutting down", zap.Error(werr), zap.Any("chain", w.chain))
				return
			}
			w.logger.Error("failed to handle EVML2 logs", zap.Error(werr))
		}

		err = w.UpdateEVML2Confirmations(w.store, w.chain, currentL1Block)
		if err != nil {
			w.logger.Error("failed to update confirmations", zap.Error(err))
		}

		w.startBlock = currentBlock
		time.Sleep(w.interval)
	}
}

func (w *EthereumL2Watcher) HandleEVML2Logs(eventIds [][]common.Hash, logs []types.Log, store Store, screener screener.Screener, contract *AtomicSwap.AtomicSwap, logger *zap.Logger) error {
	for _, log := range logs {
		switch log.Topics[0] {
		case eventIds[0][0]:
			cSwap, err := RetryWithReturnValue(func() (Swap, error) {
				return contract.AtomicSwapOrders(nil, log.Topics[1])
			}, retryCount)
			if err != nil {
				return NewNonRecoverableError(fmt.Errorf("failed to get swap order: %s", err))
			}
			handler := func() error {
				return w.HandleEVML2Initiate(log, store, cSwap, screener)
			}
			wErr := RetryOnWatcherError(handler, retryCount)
			var nonrecoverable *NonRecoverableError
			if wErr != nil && errors.As(wErr, &nonrecoverable) {
				return wErr
			}
		case eventIds[0][1]:
			handler := func() error {
				return HandleEVML2Redeem(store, log)
			}
			err := RetryOnWatcherError(handler, retryCount)
			var nonrecoverable *NonRecoverableError
			if err != nil && errors.As(err, &nonrecoverable) {
				return err
			}
		case eventIds[0][2]:
			handler := func() error {
				return HandleEVML2Refund(store, log)
			}
			err := RetryOnWatcherError(handler, retryCount)
			var nonrecoverable *NonRecoverableError
			if err != nil && errors.As(err, &nonrecoverable) {
				return err
			}
		}
	}
	return nil
}

// update confirmation status of unconfirmed initiates
func (w *EthereumL2Watcher) UpdateEVML2Confirmations(store Store, chain model.Chain, currentBlock uint64) error {
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
		if ((swap.FilledAmount != "" && swap.Status == model.NotStarted) || swap.Status == model.Initiated) && currentBlock > swap.InitiateBlockNumber+uint64(timelock) {
			swap.Status = model.Expired
			if err := store.UpdateSwap(&swap); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *EthereumL2Watcher) HandleEVML2Initiate(log types.Log, store Store, cSwap Swap, screener screener.Screener) error {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return NewRecoverableError(fmt.Errorf("failed to get swap by ocid: %s", err))
	}
	if swap.InitiateTxHash != "" {
		return nil
	}

	if screener != nil {
		isBlacklisted, err := screener.IsBlacklisted(map[string]model.Chain{cSwap.Initiator.Hex(): swap.Chain})
		if err != nil {
			return NewRecoverableError(fmt.Errorf("failed to check if address is blacklisted, %s", cSwap.Initiator.Hex()))
		}
		if isBlacklisted {
			return NewIgnorableError(fmt.Errorf("address is blacklisted, %s", cSwap.Initiator.Hex()))
		}
	}

	amount, ok := new(big.Int).SetString(swap.Amount, 10)
	if !ok {
		return NewIgnorableError(fmt.Errorf("invalid amount: %s", swap.Amount))
	}
	if cSwap.Amount.Cmp(amount) < 0 {
		return NewIgnorableError(fmt.Errorf("insufficient amount: %s", swap.Amount))
	}
	expiry, ok := new(big.Int).SetString(swap.Timelock, 10)
	if !ok {
		return NewIgnorableError(fmt.Errorf("failed to decode timelock: %s", err))
	}

	if cSwap.Expiry.Cmp(expiry) != 0 {
		return NewIgnorableError(fmt.Errorf("incorrect expiry: %s", expiry))
	}

	if strings.ToLower(cSwap.Redeemer.String()) != strings.ToLower(swap.RedeemerAddress) {
		return NewIgnorableError(fmt.Errorf("incorrect redeemer: %s", swap.RedeemerAddress))
	}

	currentL1Block, err := w.client.GetL1BlockAt(log.BlockNumber)
	if err != nil {
		return NewRecoverableError(fmt.Errorf("failed to get current block number"))
	}

	w.logger.Info(fmt.Sprintf("Mapping %d on l2 to %d", log.BlockNumber, currentL1Block))

	swap.InitiateTxHash = log.TxHash.String()
	swap.InitiateBlockNumber = currentL1Block
	swap.Status = model.Detected

	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewRecoverableError(fmt.Errorf("failed to update swap: %s", err))
	}
	return nil
}

func HandleEVML2Redeem(store Store, log types.Log) error {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return NewRecoverableError(
			fmt.Errorf("failed to get swap by ocid: %s", err),
		)
	}

	if swap.RedeemTxHash != "" {
		return nil
	}
	if len(log.Data) < 64 {
		return NewIgnorableError(
			fmt.Errorf("invalid log data: %x", log.Data),
		)
	}
	swap.Secret = hex.EncodeToString(log.Data[64:])
	swap.RedeemTxHash = log.TxHash.Hex()
	swap.Status = model.Redeemed
	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewRecoverableError(
			fmt.Errorf("failed to update swap: %s", err),
		)
	}
	return nil
}

func HandleEVML2Refund(store Store, log types.Log) error {
	swap, err := store.SwapByOCID(log.Topics[1].Hex()[2:])
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return NewRecoverableError(
			fmt.Errorf("failed to get swap by ocid: %s", err),
		)
	}
	if swap.RefundTxHash != "" {
		return nil
	}
	swap.RefundTxHash = log.TxHash.String()
	swap.Status = model.Refunded
	err = store.UpdateSwap(&swap)
	if err != nil {
		return NewRecoverableError(
			fmt.Errorf("failed to update swap: %s", err),
		)
	}
	return nil
}