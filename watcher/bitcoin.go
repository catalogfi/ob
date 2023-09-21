package watcher

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/screener"
	"github.com/catalogfi/wbtc-garden/swapper"
	"github.com/catalogfi/wbtc-garden/swapper/bitcoin"
	"github.com/ethereum/go-ethereum/common/math"
	"go.uber.org/zap"
)

type BTCWatcher struct {
	store    Store
	config   model.Config
	screener screener.Screener
	interval time.Duration
	logger   *zap.Logger
	chain    model.Chain
}

func NewBTCWatcher(store Store, chain model.Chain, config model.Config, screener screener.Screener, interval time.Duration, logger *zap.Logger) *BTCWatcher {
	return &BTCWatcher{
		chain:    chain,
		store:    store,
		config:   config,
		logger:   logger,
		screener: screener,
		interval: interval,
	}
}

func (w *BTCWatcher) Watch(ctx context.Context) {
	fmt.Println("started bitcoin watcher")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := w.ProcessBTCSwaps(); err != nil {
				w.logger.Error("failed to process swaps", zap.Error(err))
			}
			time.Sleep(w.interval)
		}
	}
}

func (w *BTCWatcher) ProcessBTCSwaps() error {
	swaps, err := w.store.GetActiveSwaps(w.chain)
	if err != nil {
		return fmt.Errorf("failed to fetch active orders %v", err)
	}
	fmt.Println(len(swaps), "btc watcher")
	for _, swap := range swaps {
		btcClient, err := LoadBTCClient(swap.Chain, w.config.Network[swap.Chain], nil)
		if err != nil {
			w.logger.Error("failed to load client", zap.Error(err))
			continue
		}
		watcher, err := LoadBTCWatcher(btcClient, swap, w.config.Network[swap.Chain])
		if err != nil {
			w.logger.Error("failed to load watcher", zap.Error(err))
			continue
		}

		if err := UpdateSwapStatus(watcher, btcClient, w.screener, w.store, &swap); err != nil {
			w.logger.Error("failed to update swap status", zap.Error(err))
			continue
		}
	}
	return nil
}

func UpdateSwapStatus(watcher swapper.Watcher, btcClient bitcoin.Client, screener screener.Screener, store Store, swap *model.AtomicSwap) error {
	if swap.InitiateTxHash == "" || (swap.InitiateTxHash != "" && strings.Compare(swap.FilledAmount, swap.Amount) < 0 && swap.Chain.IsBTC()) {
		filledAmount, txHash, err := BTCInitiateStatus(btcClient, screener, swap.Chain, swap.OnChainIdentifier)
		if err != nil {
			return err
		}

		amount, err := strconv.ParseUint(swap.Amount, 10, 64)
		if err != nil {
			return err
		}

		swap.FilledAmount = strconv.FormatUint(filledAmount, 10)
		swap.InitiateTxHash = txHash
		if filledAmount >= amount {
			_, _, isIw, err := watcher.Status(swap.InitiateTxHash)
			if err != nil {
				return err
			}

			if isIw {
				currentBlock, err := btcClient.GetTipBlockHeight()
				if err != nil {
					return err
				}
				swap.InitiateBlockNumber = currentBlock
				swap.CurrentConfirmations = swap.MinimumConfirmations
				swap.Status = model.Initiated
				swap.IsInstantWallet = true
			} else {
				swap.Status = model.Detected
			}
		}

	} else if swap.InitiateTxHash != "" && swap.Status == model.Detected {
		height, confirmations, _, err := watcher.Status(swap.InitiateTxHash)
		if err != nil {
			return err
		}

		if confirmations != swap.CurrentConfirmations {
			swap.CurrentConfirmations = confirmations
		}
		if swap.CurrentConfirmations >= swap.MinimumConfirmations {
			swap.InitiateBlockNumber = height
			swap.CurrentConfirmations = swap.MinimumConfirmations
			swap.Status = model.Initiated
		}

	} else {
		currentBlock, err := btcClient.GetTipBlockHeight()
		if err != nil {
			return err
		}

		expiry, err := strconv.ParseUint(swap.Timelock, 10, 64)
		if err != nil {
			return err
		}

		if currentBlock > swap.InitiateBlockNumber+expiry {
			refunded, txHash, err := watcher.IsRefunded()
			if err != nil {
				return err
			}
			if !refunded {
				return nil
			}
			swap.Status = model.Refunded
			swap.RefundTxHash = txHash
		} else {
			redeemed, secret, txHash, err := watcher.IsRedeemed()
			if err != nil {
				return err
			}
			if !redeemed {
				return nil
			}
			swap.Status = model.Redeemed
			swap.RedeemTxHash = txHash
			swap.Secret = hex.EncodeToString(secret)
		}
	}
	return store.UpdateSwap(swap)
}

func BTCInitiateStatus(btcClient bitcoin.Client, screener screener.Screener, chain model.Chain, scriptAddress string) (uint64, string, error) {
	addr, err := btcutil.DecodeAddress(scriptAddress, btcClient.Net())
	if err != nil {
		return 0, "", err
	}

	utxos, bal, err := btcClient.GetUTXOs(addr, 0)
	if err != nil {
		return 0, "", err
	}

	txs := make([]string, len(utxos))
	txSenders := map[string]model.Chain{}
	for i, utxo := range utxos {
		txs[i] = utxo.TxID
		tx, err := btcClient.GetTx(utxo.TxID)
		if err != nil {
			return 0, "", err
		}
		for _, vin := range tx.VINs {
			txSenders[vin.Prevout.ScriptPubKeyAddress] = chain
		}
	}

	if screener != nil {
		isBlacklisted, err := screener.IsBlacklisted(txSenders)
		if err != nil {
			return 0, "", err
		}

		if isBlacklisted {
			return 0, "", fmt.Errorf("blacklisted deposits detected")
		}
	}

	return bal, strings.Join(txs, ","), nil
}

func GetBTCConfirmations(btcClient bitcoin.Client, txHash string) (uint64, uint64, error) {
	latestDepositBlockHeight := uint64(0)
	latestDepositConfirmations := uint64(math.MaxUint64)
	txHashes := strings.Split(txHash, ",")
	for _, txHash := range txHashes {
		blockHeight, confirmations, err := btcClient.GetConfirmations(txHash)
		if err != nil {
			return 0, 0, err
		}

		if blockHeight > uint64(latestDepositBlockHeight) {
			latestDepositBlockHeight = blockHeight
		}

		if confirmations < latestDepositConfirmations {
			latestDepositConfirmations = confirmations
		}
	}
	return latestDepositBlockHeight, latestDepositConfirmations, nil
}

func LoadBTCClient(chain model.Chain, config model.NetworkConfig, btcStore bitcoin.Store) (bitcoin.Client, error) {
	indexers := []bitcoin.Indexer{}
	for iType, url := range config.RPC {
		switch iType {
		case "blockstream":
			indexers = append(indexers, bitcoin.NewBlockstream(url))
		case "mempool":
			indexers = append(indexers, bitcoin.NewMempool(url))
		default:
			return nil, fmt.Errorf("unknown indexer: %s", iType)
		}
	}
	indexer, err := bitcoin.NewMultiIndexer(indexers...)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexer: %v", err)
	}
	client := bitcoin.NewClient(indexer, chain.Params())
	if btcStore != nil {
		return bitcoin.InstantWalletWrapper(config.IWRPC, btcStore, client), nil
	}
	return client, nil
}

func LoadBTCWatcher(client bitcoin.Client, swap model.AtomicSwap, config model.NetworkConfig) (swapper.Watcher, error) {

	amt, ok := new(big.Int).SetString(swap.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount: %s", swap.Amount)
	}

	expiry, ok := new(big.Int).SetString(swap.Timelock, 10)
	if !ok {
		return nil, fmt.Errorf("invalid timelock: %s", swap.Timelock)
	}

	scriptAddr, err := btcutil.DecodeAddress(swap.OnChainIdentifier, client.Net())
	if err != nil {
		return nil, err
	}
	fmt.Println("watching :", scriptAddr, " on chain :", swap.Chain)
	return bitcoin.NewWatcher(scriptAddr, expiry.Int64(), swap.MinimumConfirmations, amt.Uint64(), config.IWRPC, client)
}
