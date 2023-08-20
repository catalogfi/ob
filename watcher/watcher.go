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
	config  model.Config
	workers int
	orders  chan model.Order
}

// NewWatcher returns a new Watcher
func NewWatcher(logger *zap.Logger, store Store, config model.Config) Watcher {
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
				w.watch(order)
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

		// Wait at least 15 seconds and until the channel is empty
		time.Sleep(15 * time.Second)
		for {
			if len(w.orders) == 0 {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (w *watcher) watch(order model.Order) {
	log := w.logger.With(zap.Uint("order id", order.ID))

	// Special cases regardless of order status
	switch {
	// Redeem and refund happens at the same time which should not happen. Defensive check.
	case (order.InitiatorAtomicSwap.RedeemTxHash != "" || order.FollowerAtomicSwap.RedeemTxHash != "") && (order.FollowerAtomicSwap.RefundTxHash != "" || order.InitiatorAtomicSwap.RefundTxHash != ""):
		order.Status = model.OrderFailedHard
	// Follower swap never gets initiated and Initiator swap is initiated and refunded
	case order.InitiatorAtomicSwap.RefundTxHash != "" && order.FollowerAtomicSwap.InitiateTxHash == "" && order.InitiatorAtomicSwap.RedeemTxHash == "":
		order.Status = model.OrderFailedSoft
	// Follower swap not been redeemed and both swap been refunded
	case order.InitiatorAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RefundTxHash != "" && order.InitiatorAtomicSwap.RefundTxHash != "":
		order.Status = model.OrderFailedSoft
	}
	// Follower has not filled the swap before the order timeout
	if order.Status == model.OrderCreated && time.Since(order.CreatedAt) > OrderTimeout {
		order.Status = model.OrderCancelled
	}
	if order.Status == model.OrderFailedHard || order.Status == model.OrderFailedSoft || order.Status == model.OrderCancelled {
		if err := w.store.UpdateOrder(&order); err != nil {
			log.Error("failed to update status", zap.Error(err), zap.Uint("status", uint(order.Status)))
			return
		}
	}

	// Fetch swapper watchers for both parties
	initiatorWatcher, err := blockchain.LoadWatcher(*order.InitiatorAtomicSwap, order.SecretHash, w.config, order.InitiatorAtomicSwap.MinimumConfirmations)
	if err != nil {
		log.Error("failed to load initiator watcher", zap.Error(err))
		return
	}
	followerWatcher, err := blockchain.LoadWatcher(*order.FollowerAtomicSwap, order.SecretHash, w.config, order.FollowerAtomicSwap.MinimumConfirmations)
	if err != nil {
		log.Error("failed to load follower watcher", zap.Error(err))
		return
	}

	// Status update
	var orderUpdate bool
	for {
		var ctn bool
		order, ctn, err = w.statusCheck(order, initiatorWatcher, followerWatcher)
		if err != nil {
			log.Error("failed to check status", zap.Error(err))
			return
		}
		if !ctn {
			break
		}
		orderUpdate = true
	}
	if orderUpdate {
		if err := w.store.UpdateOrder(&order); err != nil {
			log.Error("failed to update order", zap.Any("order", order), zap.Error(err))
		}
	}
}

func (w *watcher) statusCheck(order model.Order, initiatorWatcher, followerWatcher swapper.Watcher) (model.Order, bool, error) {
	switch order.Status {
	case model.OrderFilled:
		if time.Since(order.CreatedAt) > SwapInitiationTimeout {
			order.Status = model.OrderCancelled
			return order, true, nil
		}

		fullyFilled, txHash, amount, err := initiatorWatcher.IsDetected()
		if err != nil {
			return order, false, fmt.Errorf("failed to detect initiator swap, %w", err)
		}
		if fullyFilled {
			// AtomicSwap initiation detected
			order.Status = model.FollowerAtomicSwapDetected
			order.FollowerAtomicSwap.FilledAmount = amount
			order.FollowerAtomicSwap.InitiateTxHash = txHash
			return order, true, nil
		} else if txHash != "" && order.FollowerAtomicSwap.FilledAmount != amount {
			// AtomicSwap partial initiation detected
			order.FollowerAtomicSwap.FilledAmount = amount
			order.FollowerAtomicSwap.InitiateTxHash = txHash
			return order, true, nil
		}
	case model.InitiatorAtomicSwapDetected:
		height, conf, err := initiatorWatcher.Status(order.InitiatorAtomicSwap.InitiateTxHash)
		if err != nil {
			return order, false, err
		}
		if order.InitiatorAtomicSwap.CurrentConfirmations != conf {
			order.InitiatorAtomicSwap.CurrentConfirmations = conf
			order.InitiatorAtomicSwap.InitiateBlockNumber = height
			if order.InitiatorAtomicSwap.CurrentConfirmations > order.InitiatorAtomicSwap.MinimumConfirmations {
				order.Status = model.InitiatorAtomicSwapInitiated
			}
			return order, true, nil
		}
	case model.InitiatorAtomicSwapInitiated:
		// Follower swap initiated
		fullyFilled, txHash, amount, err := followerWatcher.IsDetected()
		if err != nil {
			return order, false, fmt.Errorf("failed to detect initiator swap, %w", err)
		}
		if fullyFilled {
			// AtomicSwap initiation detected
			order.Status = model.FollowerAtomicSwapDetected
			order.FollowerAtomicSwap.FilledAmount = amount
			order.FollowerAtomicSwap.InitiateTxHash = txHash
			return order, true, nil
		} else if txHash != "" && order.FollowerAtomicSwap.FilledAmount != amount {
			// AtomicSwap partial initiation detected
			order.FollowerAtomicSwap.FilledAmount = amount
			order.FollowerAtomicSwap.InitiateTxHash = txHash
			return order, true, nil
		}

		// Initiateor swap expired
		expired, err := initiatorWatcher.Expired()
		if err == nil && expired {
			order.Status = model.InitiatorAtomicSwapExpired
			return order, true, nil
		}

	case model.InitiatorAtomicSwapExpired:
		if order.FollowerAtomicSwap.InitiateTxHash != "" && order.FollowerAtomicSwap.RefundTxHash == "" {
			// if initiator swap expires then follower swap should also have expired
			fRefunded, txHash, err := followerWatcher.IsRefunded()
			if err != nil {
				return order, false, fmt.Errorf("failed to check refund status %v", err)
			}
			if fRefunded {
				order.Status = model.FollowerAtomicSwapRefunded
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return order, true, nil
			}

			fRedeemed, secret, txHash, err := followerWatcher.IsRedeemed()
			if err != nil {
				return order, false, fmt.Errorf("failed to check refund status %v", err)
			}
			if fRedeemed {
				order.Secret = hex.EncodeToString(secret)
				order.Status = model.FollowerAtomicSwapRedeemed
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return order, true, nil
			}
		}

		refunded, txHash, err := initiatorWatcher.IsRefunded()
		if err != nil {
			return order, false, err
		}
		if refunded {
			order.Status = model.InitiatorAtomicSwapRefunded
			order.InitiatorAtomicSwap.RefundTxHash = txHash
			return order, true, nil
		}
		redeemed, _, txHash, err := initiatorWatcher.IsRedeemed()
		if err != nil {
			return order, false, err
		}
		if redeemed {
			order.Status = model.InitiatorAtomicSwapRedeemed
			order.InitiatorAtomicSwap.RedeemTxHash = txHash
			return order, true, nil
		}
	case model.FollowerAtomicSwapDetected:
		height, conf, err := followerWatcher.Status(order.FollowerAtomicSwap.InitiateTxHash)
		if err != nil {
			return order, false, err
		}
		if order.FollowerAtomicSwap.CurrentConfirmations != conf {
			order.FollowerAtomicSwap.CurrentConfirmations = conf
			order.FollowerAtomicSwap.InitiateBlockNumber = height
			if order.FollowerAtomicSwap.CurrentConfirmations > order.FollowerAtomicSwap.MinimumConfirmations {
				order.Status = model.FollowerAtomicSwapInitiated
			}
			return order, true, nil
		}
	case model.FollowerAtomicSwapInitiated:
		// Follower swap redeemed
		redeemed, secret, txHash, err := followerWatcher.IsRedeemed()
		if err != nil {
			return order, false, fmt.Errorf("follower swap redeeming, %w", err)
		}
		if redeemed {
			order.Secret = hex.EncodeToString(secret)
			order.Status = model.FollowerAtomicSwapRedeemed
			order.FollowerAtomicSwap.RedeemTxHash = txHash
			return order, true, nil
		}
		// Check for expiry on follower swap
		fExpired, err := followerWatcher.Expired()
		if err == nil && fExpired {
			order.Status = model.FollowerAtomicSwapExpired
			return order, true, nil
		}
		if err != nil {
			return order, false, err
		}

		// Check for expiry on initiator swap
		iExpired, err := initiatorWatcher.Expired()
		if err == nil && iExpired {
			order.Status = model.InitiatorAtomicSwapExpired
			return order, true, nil
		}
		if err != nil {
			return order, false, err
		}
	case model.FollowerAtomicSwapRedeemed:
		// Initiator swap redeemed
		redeemed, _, txHash, err := initiatorWatcher.IsRedeemed()
		if err != nil {
			return order, false, fmt.Errorf("initiator swap redeeming, %w", err)
		}
		if redeemed {
			order.Status = model.InitiatorAtomicSwapRedeemed
			order.InitiatorAtomicSwap.RedeemTxHash = txHash
			return order, true, nil
		}

		// Check for expiry on initiator swap
		iExpired, err := initiatorWatcher.Expired()
		if err != nil {
			return order, false, err
		}
		if iExpired {
			order.Status = model.InitiatorAtomicSwapExpired
			return order, true, nil
		}

	case model.FollowerAtomicSwapExpired:
		// if initiator swap expires then follower swap should also have expired
		fRefunded, txHash, err := followerWatcher.IsRefunded()
		if err != nil {
			return order, false, fmt.Errorf("failed to check refund status %v", err)
		}
		if fRefunded {
			order.Status = model.FollowerAtomicSwapRefunded
			order.FollowerAtomicSwap.RefundTxHash = txHash
			return order, true, nil
		}

		fRedeemed, secret, txHash, err := followerWatcher.IsRedeemed()
		if err != nil {
			return order, false, fmt.Errorf("failed to check refund status %v", err)
		}
		if fRedeemed {
			order.Secret = hex.EncodeToString(secret)
			order.Status = model.FollowerAtomicSwapRedeemed
			order.FollowerAtomicSwap.RefundTxHash = txHash
			return order, true, nil
		}
	case model.InitiatorAtomicSwapRedeemed:
		order.Status = model.OrderExecuted
		return order, true, nil
	case model.FollowerAtomicSwapRefunded:
		// Check for expiry on initiator swap
		iExpired, err := initiatorWatcher.Expired()
		if err != nil {
			return order, false, err
		}
		if iExpired {
			order.Status = model.InitiatorAtomicSwapExpired
			return order, true, nil
		}
	}

	return order, false, nil
}
