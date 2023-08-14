package watcher

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/catalogfi/wbtc-garden/blockchain"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/swapper"
)

const SwapInitiationTimeout = 30 * time.Minute

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
	store   Store
	config  model.Config
	workers int
	orders  chan model.Order
}

// NewWatcher returns a new Watcher
func NewWatcher(store Store, config model.Config) Watcher {
	workers := runtime.NumCPU() * 4
	if workers == 0 {
		workers = 4
	}
	return &watcher{
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
				if err := w.watch(order); err != nil {
					fmt.Printf("error updating order: %v\n", err)
				}
			}
		}()
	}

	defer close(w.orders)
	for {
		orders, err := w.store.GetActiveOrders()
		if err != nil {
			fmt.Printf("error getting active orders: %v\n", err)
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

func (w *watcher) watch(order model.Order) error {

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
	if order.Status == model.OrderFailedHard || order.Status == model.OrderFailedSoft {
		return w.store.UpdateOrder(&order)
	}

	// Fetch swapper watchers for both parties
	initiatorWatcher, err := blockchain.LoadWatcher(*order.InitiatorAtomicSwap, order.SecretHash, w.config.RPC, order.InitiatorAtomicSwap.MinimumConfirmations)
	if err != nil {
		return err
	}
	followerWatcher, err := blockchain.LoadWatcher(*order.FollowerAtomicSwap, order.SecretHash, w.config.RPC, order.FollowerAtomicSwap.MinimumConfirmations)
	if err != nil {
		return err
	}

	// Status update
	var orderUpdate bool
	for {
		var ctn bool
		order, ctn, err = w.statusCheck(order, initiatorWatcher, followerWatcher)
		if err != nil {
			return err
		}
		if !ctn {
			break
		}
		orderUpdate = true
	}
	if orderUpdate {
		return w.store.UpdateOrder(&order)
	}
	return nil
}

func (w *watcher) statusCheck(order model.Order, initiatorWatcher, followerWatcher swapper.Watcher) (model.Order, bool, error) {
	initiatorMinRefundTime := order.CreatedAt.Add(24 * time.Hour)
	followerMinRefundTime := order.CreatedAt.Add(12 * time.Hour)

	switch order.Status {
	case model.OrderFilled:
		// Initiator swap initiated
		initiated, txHashes, err := initiatorWatcher.IsInitiated()
		if err != nil {
			return order, false, err
		}
		if initiated {
			order.Status = model.InitiatorAtomicSwapInitiated
			order.InitiatorAtomicSwap.InitiateTxHash = strings.Join(txHashes, ",")
			return order, true, nil
		}

		// Order expired
		if time.Since(order.CreatedAt) > SwapInitiationTimeout {
			order.Status = model.OrderCancelled
			return order, true, nil
		}
	case model.InitiatorAtomicSwapInitiated:
		// Initiator swap refunded
		if time.Now().After(initiatorMinRefundTime) {
			refunded, txHash, err := initiatorWatcher.IsRefunded()
			if err != nil {
				return order, false, err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.InitiatorAtomicSwap.RefundTxHash = txHash
				return order, true, nil
			}
		}

		// Follower swap initiated
		initiated, txHashes, err := followerWatcher.IsInitiated()
		if err != nil {
			return order, false, err
		}
		if initiated {
			order.Status = model.FollowerAtomicSwapInitiated
			order.FollowerAtomicSwap.InitiateTxHash = strings.Join(txHashes, ",")
			return order, true, nil
		}
	case model.FollowerAtomicSwapInitiated:
		//  Check if follower swap is refunded
		if time.Now().After(followerMinRefundTime) {
			refunded, txHash, err := followerWatcher.IsRefunded()
			if err != nil {
				return order, false, err
			}
			if refunded {
				order.Status = model.FollowerAtomicSwapRefunded
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return order, true, nil
			}
		}

		//  Check if initiator swap is refunded
		if time.Now().After(initiatorMinRefundTime) {
			refunded, txHash, err := initiatorWatcher.IsRefunded()
			if err != nil {
				return order, false, err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.InitiatorAtomicSwap.RefundTxHash = txHash
				return order, true, nil
			}
		}

		// Follower swap redeemed
		redeemed, secret, txHash, err := followerWatcher.IsRedeemed()
		if err != nil {
			return order, false, err
		}
		if redeemed {
			order.Secret = hex.EncodeToString(secret)
			order.Status = model.FollowerAtomicSwapRedeemed
			order.FollowerAtomicSwap.RedeemTxHash = txHash
			return order, true, nil
		}
	case model.FollowerAtomicSwapRedeemed:
		// Initiator swap refunded
		if time.Now().After(initiatorMinRefundTime) {
			refunded, txHash, err := initiatorWatcher.IsRefunded()
			if err != nil {
				return order, false, err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.InitiatorAtomicSwap.RefundTxHash = txHash
				return order, true, err
			}
		}

		// Initiator swap redeemed
		redeemed, _, txHash, err := initiatorWatcher.IsRedeemed()
		if err != nil {
			return order, false, err
		}
		if redeemed {
			order.Status = model.InitiatorAtomicSwapRedeemed
			order.InitiatorAtomicSwap.RedeemTxHash = txHash
			return order, true, nil
		}
	case model.InitiatorAtomicSwapRedeemed:
		order.Status = model.OrderExecuted
		return order, true, nil
	case model.FollowerAtomicSwapRefunded:
		if time.Now().After(initiatorMinRefundTime) {
			refunded, txHash, err := initiatorWatcher.IsRefunded()
			if err != nil {
				return order, false, err
			}
			if refunded {
				order.InitiatorAtomicSwap.RefundTxHash = txHash
				return order, false, err
			}
		}
	}

	return order, false, nil
}
