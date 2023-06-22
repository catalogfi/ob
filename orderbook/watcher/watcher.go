package watcher

import (
	"encoding/hex"
	"fmt"

	"github.com/susruth/wbtc-garden/orderbook/model"
	"github.com/susruth/wbtc-garden/swapper"
)

type Store interface {
	// update an order, used to update the status
	UpdateOrder(order *model.Order) error
	// get all active orders
	GetActiveOrders() ([]model.Order, error)
}

type watcher struct {
	store Store
}

func NewWatcher(store Store) *watcher {
	return &watcher{
		store: store,
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
	}
}

// TODO: does not handle the case where the order is refunded
func (w *watcher) watch(order model.Order) error {
	iW, fW, err := w.getWatchers(order)
	if err != nil {
		return err
	}

	if order.Status == model.OrderFilled {
		initiated, txHash, err := iW.IsInitiated()
		if err != nil {
			return err
		}
		if initiated {
			order.Status = model.InitiatorAtomicSwapInitiated
			order.InitiatorAtomicSwap.InitiateTxHash = txHash
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.InitiatorAtomicSwapInitiated {
		initiated, txHash, err := fW.IsInitiated()
		if err != nil {
			return err
		}
		if initiated {
			order.Status = model.FollowerAtomicSwapInitiated
			order.FollowerAtomicSwap.InitiateTxHash = txHash
			if err := w.store.UpdateOrder(&order); err != nil {
				return err
			}
		}
	}

	if order.Status == model.FollowerAtomicSwapInitiated {
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

	return nil
}

func (w *watcher) getInitiateWatcher(chain model.Chain) (swapper.Watcher, error) {
	return nil, nil
}

func (w *watcher) getWatchers(order model.Order) (swapper.Watcher, swapper.Watcher, error) {
	// fromChain, toChain, fromAsset, toAsset, err := model.ParseOrderPair(order.OrderPair)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// if fromChain == model.Ethereum && toChain == model.Bitcoin {
	// 	initiatorAddr := common.HexToAddress(order.InitiatorAtomicSwap.InitiatorAddress)
	// 	redeemerAddr := common.HexToAddress(order.InitiatorAtomicSwap.RedeemerAddress)
	// 	var tokenAddr, deployerAddr common.Address
	// 	if fromAsset[:9] != "secondary" {
	// 		tokenAddr = common.HexToAddress(string(fromAsset[9:]))
	// 		deployerAddr = common.HexToAddress("")
	// 	}

	// 	watcher, err := ethereum.NewWatcher(initiatorAddr, redeemerAddr, deployerAddr, tokenAddr, )
	// 	return nil, nil, nil
	// }
	return nil, nil, nil
}
