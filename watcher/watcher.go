package watcher

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"time"

	"github.com/catalogfi/wbtc-garden/blockchain"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/swapper"
	"go.uber.org/zap"
)

const SwapInitiationTimeout = 12 * time.Hour
const OrderTimeout = 3 * time.Minute

type Store interface {
	// UpdateOrder updates and order status in the db
	UpdateOrder(order *model.Order) error
	// GetActiveOrders fetches all orders which are active.
	GetActiveOrders() ([]model.Order, error)
}

// Watcher watches the blockchain and update order status accordingly.
type Watcher interface {
	Run()
}

type watcher struct {
	logger  *zap.Logger
	store   Store
	config  model.Network
	workers int
	orders  chan model.Order
}

// NewWatcher returns a new Watcher
func NewWatcher(logger *zap.Logger, store Store, config model.Network) Watcher {
	workers := runtime.NumCPU() * 4
	if workers == 0 {
		workers = 4
	}
	childLogger := logger.With(zap.String("service", "watcher"))
	return &watcher{
		logger:  childLogger,
		store:   store,
		config:  config,
		workers: workers,
		orders:  make(chan model.Order, 32),
	}
}

func (w *watcher) Run() {
	for i := 0; i < w.workers; i++ {
		go func() {
			for order := range w.orders {
				ProcessOrder(order, w.store, w.config, w.logger.With(zap.Uint("order id", order.ID)))
			}
		}()
	}

	defer close(w.orders)
	for {
		orders, err := w.store.GetActiveOrders()
		if err != nil {
			w.logger.Error("get active order", zap.Error(err))
			continue
		}
		for _, order := range orders {
			w.orders <- order
		}

		// Wait at least 5 seconds and until the channel is empty
		// reduced to test increse in performance
		time.Sleep(5 * time.Second)
		for {
			if len(w.orders) == 0 {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func ProcessOrder(order model.Order, store Store, config model.Network, logger *zap.Logger) {
	// Special cases regardless of order status
	switch {
	case order.Status == model.Created && time.Since(order.CreatedAt) > OrderTimeout:
		order.Status = model.Cancelled

	// Redeem and refund happens at the same time which should not happen. Defensive check.
	case (order.InitiatorAtomicSwap.RedeemTxHash != "" || order.FollowerAtomicSwap.RedeemTxHash != "") && (order.FollowerAtomicSwap.RefundTxHash != "" || order.InitiatorAtomicSwap.RefundTxHash != ""):
		logger.Error("atomic swap hard failed as someone both redeemed and refunded")
		order.Status = model.FailedHard
	// Follower swap never gets initiated and Initiator swap is initiated and refunded
	case order.InitiatorAtomicSwap.RefundTxHash != "" && order.FollowerAtomicSwap.InitiateTxHash == "" && order.InitiatorAtomicSwap.RedeemTxHash == "":
		logger.Error("atomic swap soft failed due to initiator refunding")
		order.Status = model.FailedSoft
	// Follower swap not been redeemed and both swap been refunded
	case order.InitiatorAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RefundTxHash != "" && order.InitiatorAtomicSwap.RefundTxHash != "":
		logger.Error("atomic swap soft failed due to both parties refunding")
		order.Status = model.FailedSoft
	}
	// Follower has not filled the swap before the order timeout
	if order.Status == model.Created && time.Since(order.CreatedAt) > OrderTimeout {
		logger.Info("atomic swap cancelled due to fill timeout")
		order.Status = model.Cancelled
	}
	if order.Status == model.Filled && time.Since(order.CreatedAt) > SwapInitiationTimeout && order.InitiatorAtomicSwap.Status == model.NotStarted {
		logger.Error("atomic swap cancelled due to initiator failing to initiate before timeout")
		order.Status = model.Cancelled
	}

	if order.Status == model.FailedHard || order.Status == model.FailedSoft || order.Status == model.Cancelled {
		if err := store.UpdateOrder(&order); err != nil {
			logger.Error("failed to update status", zap.Error(err), zap.Uint("status", uint(order.Status)))
			return
		}
	}

	if order.Status != model.Filled {
		return
	}

	if order.Status == model.Created {
		return
	}

	// Fetch swapper watchers for both parties
	initiatorWatcher, err := blockchain.LoadWatcher(*order.InitiatorAtomicSwap, order.SecretHash, config, order.InitiatorAtomicSwap.MinimumConfirmations)
	if err != nil {
		logger.Error("failed to load initiator watcher", zap.Error(err))
		return
	}
	followerWatcher, err := blockchain.LoadWatcher(*order.FollowerAtomicSwap, order.SecretHash, config, order.FollowerAtomicSwap.MinimumConfirmations)
	if err != nil {
		logger.Error("failed to load follower watcher", zap.Error(err))
		return
	}

	// Status update
	var orderUpdate bool
	for {
		var ctn bool
		order, ctn, err = updateStatus(logger, order, initiatorWatcher, followerWatcher)
		if err != nil {
			logger.Error("failed to update status", zap.Error(err))
			// return
		}
		if !ctn {
			break
		}
		orderUpdate = true
	}

	if orderUpdate {
		if err := store.UpdateOrder(&order); err != nil {
			logger.Error("failed to update order", zap.Any("order", order), zap.Error(err))
		}
	}
}

func updateSwapStatus(log *zap.Logger, swap model.AtomicSwap, watcher swapper.Watcher) (model.AtomicSwap, bool, error) {
	switch swap.Status {
	case model.NotStarted:
		fullyFilled, txHash, amount, err := watcher.IsDetected()
		if err != nil {
			return swap, false, fmt.Errorf("failed to detect swap, %w", err)
		}
		if fullyFilled {
			// AtomicSwap initiation detected
			swap.Status = model.Detected
			swap.FilledAmount = amount
			swap.InitiateTxHash = txHash
			log.Info("atomic swap detected", zap.Uint("swap id", swap.ID), zap.String("txhash", txHash), zap.String("amount", amount))
			return swap, true, nil
		} else if txHash != "" && swap.FilledAmount != amount {
			// AtomicSwap partial initiation detected
			swap.FilledAmount = amount
			swap.InitiateTxHash = txHash
			return swap, true, nil
		}
	case model.Detected:
		height, conf, err := watcher.Status(swap.InitiateTxHash)
		if err != nil {
			return swap, false, err
		}
		if swap.CurrentConfirmations != conf {
			if conf > swap.MinimumConfirmations {
				conf = swap.MinimumConfirmations
			}
			swap.CurrentConfirmations = conf
			swap.InitiateBlockNumber = height
			if swap.CurrentConfirmations == swap.MinimumConfirmations {
				log.Info("atomic swap initiated", zap.Uint("swap id", swap.ID), zap.Uint64("confirmations", conf), zap.Uint64("height", height))
				swap.Status = model.Initiated
			}
			return swap, true, nil
		}
	case model.Initiated:
		// Swap redeemed
		redeemed, secret, txHash, err := watcher.IsRedeemed()
		if err != nil {
			return swap, false, fmt.Errorf("follower swap redeeming, %w", err)
		}
		if redeemed {
			log.Info("atomic swap redeemed", zap.Uint("swap id", swap.ID), zap.String("tx hash", txHash))
			swap.Secret = hex.EncodeToString(secret)
			swap.Status = model.Redeemed
			swap.RedeemTxHash = txHash
			return swap, true, nil
		}

		// Swap Expired
		watcher.Expired()
		expired, err := watcher.Expired()
		if err != nil {
			return swap, false, err
		}
		if expired {
			log.Info("atomic swap expired", zap.Uint("swap id", swap.ID))
			swap.Status = model.Expired
			return swap, true, nil
		}

	case model.Expired:
		// Swap expired
		refunded, txHash, err := watcher.IsRefunded()
		if err != nil {
			return swap, false, fmt.Errorf("failed to check refund status %v", err)
		}
		if refunded {
			log.Info("atomic swap expired", zap.Uint("swap id", swap.ID), zap.String("tx hash", txHash))
			swap.Status = model.Refunded
			swap.RefundTxHash = txHash
			return swap, true, nil
		}
	}
	return swap, false, nil
}

func updateStatus(log *zap.Logger, order model.Order, initiatorWatcher, followerWatcher swapper.Watcher) (model.Order, bool, error) {
	initiatorSwap, ictn, ierr := updateSwapStatus(log, *order.InitiatorAtomicSwap, initiatorWatcher)
	if ierr != nil {
		return order, false, ierr
	}
	followerSwap, fctn, ferr := updateSwapStatus(log, *order.FollowerAtomicSwap, followerWatcher)
	if ferr != nil {
		return order, false, ferr
	}
	if ictn || fctn {
		if order.Secret != followerSwap.Secret {
			order.Secret = followerSwap.Secret
		}
		if initiatorSwap.Status == model.Redeemed && followerSwap.Status == model.Redeemed {
			log.Info("atomic swap executed", zap.Uint("order id", order.ID))
			order.Status = model.Executed
		}
		order.InitiatorAtomicSwap = &initiatorSwap
		order.FollowerAtomicSwap = &followerSwap
		return order, true, nil
	}
	return order, false, nil
}
