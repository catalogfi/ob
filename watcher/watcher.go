package watcher

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/catalogfi/wbtc-garden/blockchain"
	"github.com/catalogfi/wbtc-garden/model"
)

type Store interface {
	// update an order, used to update the status
	UpdateOrder(order *model.Order) error
	// get all active orders
	GetActiveOrders() ([]model.Order, error)
	// get locked value for a user on a chain
	GetValueLocked(user string, chain model.Chain) (*big.Int, error)
}

type watcher struct {
	store  Store
	config model.Config
}

func NewWatcher(store Store, config model.Config) *watcher {
	return &watcher{
		store:  store,
		config: config,
	}
}

func (w *watcher) Run() {
	for {
		orders, err := w.store.GetActiveOrders()
		if err != nil {
			fmt.Printf("error getting active orders: %v\n", err)
			continue
		}
		for _, order := range orders {
			if err := w.watch(order); err != nil {
				fmt.Printf("error updating order: %v\n", err)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func (w *watcher) watch(order model.Order) error {

	// to check isFinal when changing status from 2 -> 3
	iW, err := blockchain.LoadWatcher(*order.InitiatorAtomicSwap, order.SecretHash, w.config.RPC, order.InitiatorAtomicSwap.MinimumConfirmations)
	if err != nil {
		return err
	}

	// to check isFinal when changing status from 3 -> 4
	fW, err := blockchain.LoadWatcher(*order.FollowerAtomicSwap, order.SecretHash, w.config.RPC, order.FollowerAtomicSwap.MinimumConfirmations)
	if err != nil {
		return err
	}

	if order.Status == model.OrderFilled {
		expired := w.orderExpired(order)
		if expired {
			order.Status = model.OrderCancelled
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
		initiated, txHash, err := iW.IsInitiated()
		if err != nil {
			return err
		}
		if initiated {
			order.Status = model.InitiatorAtomicSwapInitiated
			order.InitiatorAtomicSwap.InitiateTxHash = strings.Join(txHash, ",")
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.InitiatorAtomicSwapInitiated {
		refunded, txHash, err := iW.IsRefunded()
		if err != nil {
			return err
		}
		if refunded {
			order.Status = model.InitiatorAtomicSwapRefunded
			order.InitiatorAtomicSwap.RefundTxHash = txHash
			return w.store.UpdateOrder(&order)
		}
		initiated, txHashs, err := fW.IsInitiated()
		if err != nil {
			return err
		}
		if initiated {
			order.Status = model.FollowerAtomicSwapInitiated
			order.FollowerAtomicSwap.InitiateTxHash = strings.Join(txHashs, ",")
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.FollowerAtomicSwapInitiated {
		refunded, txHash, err := fW.IsRefunded()
		if err != nil {
			return err
		}
		if refunded {
			order.Status = model.FollowerAtomicSwapRefunded
			order.FollowerAtomicSwap.RefundTxHash = txHash
			return w.store.UpdateOrder(&order)
		}

		refunded, txHash, err = iW.IsRefunded()
		if err != nil {
			return err
		}
		if refunded {
			order.Status = model.InitiatorAtomicSwapRefunded
			order.InitiatorAtomicSwap.RefundTxHash = txHash
			return w.store.UpdateOrder(&order)
		}

		redeemed, secret, txHash, err := fW.IsRedeemed()
		if err != nil {
			return err
		}
		if redeemed {
			order.Secret = hex.EncodeToString(secret)
			order.Status = model.FollowerAtomicSwapRedeemed
			order.FollowerAtomicSwap.RedeemTxHash = txHash
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.FollowerAtomicSwapRedeemed {
		refunded, txHash, err := iW.IsRefunded()
		if err != nil {
			return err
		}
		if refunded {
			order.Status = model.InitiatorAtomicSwapRefunded
			order.InitiatorAtomicSwap.RefundTxHash = txHash
			return w.store.UpdateOrder(&order)
		}
		redeemed, _, txHash, err := iW.IsRedeemed()
		if err != nil {
			return err
		}
		if redeemed {
			order.Status = model.InitiatorAtomicSwapRedeemed
			order.InitiatorAtomicSwap.RedeemTxHash = txHash
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.InitiatorAtomicSwapRedeemed {
		order.Status = model.OrderExecuted
		if err := w.store.UpdateOrder(&order); err != nil {
			return err
		}
	}

	if order.Status == model.FollowerAtomicSwapRefunded {
		refunded, txHash, err := iW.IsRefunded()
		if err != nil {
			return err
		}
		if refunded {
			order.InitiatorAtomicSwap.RefundTxHash = txHash
			w.store.UpdateOrder(&order)
		}
	}

	if (order.InitiatorAtomicSwap.RedeemTxHash != "" || order.FollowerAtomicSwap.RedeemTxHash != "") && (order.FollowerAtomicSwap.RefundTxHash != "" || order.InitiatorAtomicSwap.RefundTxHash != "") {
		order.Status = model.OrderFailedHard
		if err := w.store.UpdateOrder(&order); err != nil {
			return err
		}
	}

	if order.InitiatorAtomicSwap.RefundTxHash != "" && order.FollowerAtomicSwap.InitiateTxHash == "" && order.InitiatorAtomicSwap.RedeemTxHash == "" {

		order.Status = model.OrderFailedSoft
		if err := w.store.UpdateOrder(&order); err != nil {
			return err
		}
	}

	if order.InitiatorAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RedeemTxHash == "" && order.FollowerAtomicSwap.RefundTxHash != "" && order.InitiatorAtomicSwap.RefundTxHash != "" {
		order.Status = model.OrderFailedSoft
		if err := w.store.UpdateOrder(&order); err != nil {
			return err
		}
	}

	return nil
}

func (w *watcher) orderExpired(order model.Order) bool {
	return time.Since(order.CreatedAt) > time.Hour*12
}
