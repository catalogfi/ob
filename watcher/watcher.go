package watcher

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
)

type Store interface {
	// update an order, used to update the status
	UpdateOrder(order *model.Order) error
	// get all active orders
	GetActiveOrders() ([]model.Order, error)
	// get locked value for a user on a chain
	GetValueLocked(user string, chain model.Chain) (int64, error)
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

	initiatorLockValue, err := w.store.GetValueLocked(order.Maker, order.InitiatorAtomicSwap.Chain)
	if err != nil {
		return nil
	}

	initiatorMinConfirmations := GetMinConfirmations(initiatorLockValue, order.InitiatorAtomicSwap.Chain)

	//to check isFinal when changing status from 2 -> 3
	iW, err := blockchain.LoadWatcher(*order.InitiatorAtomicSwap, order.SecretHash, w.config.RPC, initiatorMinConfirmations)
	if err != nil {
		return err
	}

	followerLockValue, err := w.store.GetValueLocked(order.Taker, order.InitiatorAtomicSwap.Chain)
	if err != nil {
		return nil
	}
	followerMinConfirmations := GetMinConfirmations(followerLockValue, order.InitiatorAtomicSwap.Chain)

	//to check isFinal when changing status from 3 -> 4
	fW, err := blockchain.LoadWatcher(*order.FollowerAtomicSwap, order.SecretHash, w.config.RPC, followerMinConfirmations)
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
		expired, err := iW.Expired()
		if err != nil {
			return err
		}
		if expired {
			refunded, txHash, err := iW.IsRefunded()
			if err != nil {
				return err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return w.store.UpdateOrder(&order)
			}
		}
		initiated, txHash, err := fW.IsInitiated()
		if err != nil {
			return err
		}
		if initiated {
			order.Status = model.FollowerAtomicSwapInitiated
			order.FollowerAtomicSwap.InitiateTxHash = strings.Join(txHash, ",")
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.FollowerAtomicSwapInitiated {
		fmt.Println("ckeckpoint one")
		expired, err := fW.Expired()
		if err != nil {
			return err
		}
		if expired {
			refunded, txHash, err := fW.IsRefunded()
			if err != nil {
				return err
			}
			if refunded {
				order.Status = model.FollowerAtomicSwapRefunded
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return w.store.UpdateOrder(&order)
			}
		}
		expired, err = iW.Expired()
		if err != nil {
			return err
		}
		if expired {
			refunded, txHash, err := iW.IsRefunded()
			if err != nil {
				return err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.FollowerAtomicSwap.RefundTxHash = txHash
				return w.store.UpdateOrder(&order)
			}
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

		expired, err := iW.Expired()
		if err != nil {
			return err
		}
		if expired {
			refunded, txHash, err := iW.IsRefunded()
			if err != nil {
				return err
			}
			if refunded {
				order.Status = model.InitiatorAtomicSwapRefunded
				order.InitiatorAtomicSwap.RefundTxHash = txHash
				return w.store.UpdateOrder(&order)
			}
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

func GetMinConfirmations(value int64, chain model.Chain) uint64 {
	if chain.IsBTC() {
		switch {
		case value < 10000:
			return 1

		case value < 100000:
			return 2

		case value < 1000000:
			return 4

		case value < 10000000:
			return 6

		case value < 100000000:
			return 8

		default:
			return 12
		}
	} else if chain.IsEVM() {
		switch {
		case value < 10000:
			return 6

		case value < 100000:
			return 12

		case value < 1000000:
			return 18

		case value < 10000000:
			return 24

		case value < 100000000:
			return 30

		default:
			return 100
		}
	}
	return 0
}
