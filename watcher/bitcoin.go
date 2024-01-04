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
	"github.com/catalogfi/orderbook/model"
	"github.com/catalogfi/orderbook/screener"
	"github.com/catalogfi/orderbook/swapper"
	"github.com/catalogfi/orderbook/swapper/bitcoin"
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

type Confirmations struct {
	LatestTxConfirmations uint64
	LatestTxHeight        uint64
	FirstTxConfirmations  uint64
	FirstTxHeight         uint64
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
	w.logger.Info("started bitcoin watcher", zap.String("chain :", string(w.chain)))
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

	for _, swap := range swaps {
		fmt.Println(swap.ID)
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

		if err := UpdateSwapStatus(watcher, btcClient, w.screener, w.store, &swap, w.config.Network[swap.Chain].Expiry); err != nil {
			w.logger.Error("failed to update swap status", zap.Error(err))
			continue
		}
	}
	return nil
}

func UpdateSwapStatus(watcher swapper.Watcher, btcClient bitcoin.Client, screener screener.Screener, store Store, swap *model.AtomicSwap, expiry int64) error {

	var err error
	amount, err := strconv.ParseUint(swap.Amount, 10, 64)
	if err != nil {
		return err
	}
	if swap.Status == model.NotStarted {
		filledAmount, txHash, err := BTCInitiateStatus(btcClient, screener, swap.Chain, swap.OnChainIdentifier)
		if err != nil {
			return err
		}
		if filledAmount == 0 {
			return nil
		}

		swap.FilledAmount = strconv.FormatUint(filledAmount, 10)
		swap.InitiateTxHash = txHash
		if filledAmount >= amount {
			_, conf, isIw, err := watcher.Status(txHash)
			if err != nil {
				return err
			}
			if conf > 2 {
				order, err := store.GetOrderBySwapID(swap.ID)
				if err != nil {
					return fmt.Errorf("failed to get order of a non valid tx:%v", err)
				}
				order.Status = model.Cancelled
				if err = store.UpdateOrder(order); err != nil {
					return fmt.Errorf("failed to update a non valid order:%v", err)
				}
				return nil
				// TODO: blacklist this maker
			}
			if isIw {
				swap.Status = model.Initiated
				swap.IsInstantWallet = true
			} else {
				swap.Status = model.Detected
			}
		}

	} else if swap.InitiateTxHash != "" && swap.Status == model.Detected {
		filledAmount, txHash, err := BTCInitiateStatus(btcClient, screener, swap.Chain, swap.OnChainIdentifier)
		if err != nil {
			return err
		}
		swap.FilledAmount = strconv.FormatUint(filledAmount, 10)
		swap.InitiateTxHash = txHash
		confirmations, err := GetBTCConfirmations(btcClient, txHash)
		if err != nil {
			return err
		}
		if confirmations.LatestTxConfirmations > 0 && swap.InitiateBlockNumber == 0 && filledAmount >= amount {
			if confirmations.LatestTxConfirmations-confirmations.FirstTxConfirmations >= uint64(expiry/2) {
				return store.UpdateSwap(swap)
			}
			swap.InitiateBlockNumber = confirmations.LatestTxHeight
			swap.Status = model.Initiated
		}
		swap.CurrentConfirmations = confirmations.LatestTxConfirmations
		if swap.CurrentConfirmations >= swap.MinimumConfirmations {
			swap.CurrentConfirmations = swap.MinimumConfirmations
		}

	} else if swap.Status != model.Redeemed && swap.Status != model.Refunded {
		currentBlock, err := btcClient.GetTipBlockHeight()
		if err != nil {
			return err
		}

		timelock, err := strconv.ParseUint(swap.Timelock, 10, 64)
		if err != nil {
			return err
		}

		if swap.InitiateBlockNumber > 0 && currentBlock >= swap.InitiateBlockNumber+timelock {
			refunded, txHash, err := watcher.IsRefunded()
			if err != nil {
				return err
			}
			if !refunded {
				if swap.Status == model.Expired {
					return nil
				}
				swap.Status = model.Expired
			} else {
				swap.Status = model.Refunded
				swap.RefundTxHash = txHash
			}

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
	} else if swap.IsInstantWallet && swap.InitiateBlockNumber == 0 && swap.Status == model.Initiated {
		//when we detect iw, we set status to inited, but we would not
		//know the block number, so we do that here
		confs, err := GetBTCConfirmations(btcClient, swap.InitiateTxHash)
		if err != nil {
			return err
		}
		if confs.LatestTxConfirmations == 0 {
			return nil
		}
		// we have atleast 1 confirmation at this point
		if swap.InitiateBlockNumber == 0 {
			swap.InitiateBlockNumber = confs.LatestTxHeight
		}
		swap.CurrentConfirmations = confs.LatestTxConfirmations
		if confs.LatestTxConfirmations >= swap.MinimumConfirmations {
			swap.CurrentConfirmations = swap.MinimumConfirmations
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
	if err != nil || bal == 0 {
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

	if screener != nil && chain.IsMainnet() {
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

func GetBTCConfirmations(btcClient bitcoin.Client, txHash string) (Confirmations, error) {
	var conf Confirmations

	txHashes := strings.Split(txHash, ",")
	blockHeight, confirmations, err := btcClient.GetConfirmations(txHashes[0])
	if err != nil {
		return conf, err
	}

	conf.LatestTxConfirmations = confirmations
	conf.FirstTxConfirmations = confirmations
	conf.LatestTxHeight = blockHeight
	conf.FirstTxHeight = blockHeight

	for _, txHash := range txHashes[1:] {
		blockHeight, confirmations, err := btcClient.GetConfirmations(txHash)
		if err != nil {
			return conf, err
		}

		if blockHeight < uint64(conf.LatestTxHeight) {
			conf.LatestTxHeight = blockHeight
			conf.LatestTxConfirmations = confirmations
		}

		if blockHeight > uint64(conf.FirstTxHeight) {
			conf.FirstTxHeight = blockHeight
			conf.FirstTxConfirmations = confirmations
		}
	}
	return conf, nil
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

	return bitcoin.NewWatcher(scriptAddr, expiry.Int64(), swap.MinimumConfirmations, amt.Uint64(), config.IWRPC, client)
}
