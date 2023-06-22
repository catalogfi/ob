package watcher

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/susruth/wbtc-garden/orderbook/model"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
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

func (w *watcher) watch(order model.Order) error {
	fromChain, toChain, fromAsset, toAsset, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return err
	}

	iW, err := w.getWatcher(fromChain, fromAsset, order.SecretHash, *order.InitiatorAtomicSwap)
	if err != nil {
		return err
	}

	fW, err := w.getWatcher(toChain, toAsset, order.SecretHash, *order.FollowerAtomicSwap)
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

	if order.InitiatorAtomicSwap.RefundTxHash != "" && order.FollowerAtomicSwap.InitiateTxHash == "" {
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

	if (order.InitiatorAtomicSwap.RedeemTxHash != "" || order.FollowerAtomicSwap.RedeemTxHash != "") && (order.FollowerAtomicSwap.RefundTxHash != "" || order.InitiatorAtomicSwap.RefundTxHash != "") {
		order.Status = model.OrderFailedHard
		if err := w.store.UpdateOrder(&order); err != nil {
			return err
		}
	}

	return nil
}

func (w *watcher) getWatcher(chain model.Chain, asset model.Asset, secretHashStr string, atomicSwap model.AtomicSwap) (swapper.Watcher, error) {
	secretHash, err := hex.DecodeString(secretHashStr)
	if err != nil {
		return nil, err
	}

	if chain == model.Bitcoin {
		vals := strings.Split(os.Getenv("BITCOIN_RPC"), "-")
		var params *chaincfg.Params
		switch vals[0] {
		case "mainnet":
			params = &chaincfg.MainNetParams
		case "testnet":
			params = &chaincfg.TestNet3Params
		case "regtest":
			params = &chaincfg.RegressionNetParams
		default:
			return nil, fmt.Errorf("invalid bitcoin network: %s", vals[0])
		}
		btcClient := bitcoin.NewClient(vals[1], params)

		initiatorAddr, err := btcutil.DecodeAddress(atomicSwap.InitiatorAddress, params)
		if err != nil {
			return nil, err
		}

		redeemerAddr, err := btcutil.DecodeAddress(atomicSwap.RedeemerAddress, params)
		if err != nil {
			return nil, err
		}

		expiry, err := strconv.ParseInt(atomicSwap.Timelock, 10, 32)
		if err != nil {
			return nil, err
		}

		amount, err := strconv.ParseUint(atomicSwap.Amount, 10, 64)
		if err != nil {
			return nil, err
		}

		return bitcoin.NewWatcher(initiatorAddr, redeemerAddr, secretHash, expiry, amount, btcClient)
	}

	client, err := ethereum.ClientFromChain(chain)
	if err != nil {
		return nil, err
	}
	deployerAddr := common.HexToAddress(os.Getenv(fmt.Sprintf("%s_DEPLOYER", strings.ToUpper(string(chain)))))
	initiatorAddr := common.HexToAddress(atomicSwap.InitiatorAddress)
	redeemerAddr := common.HexToAddress(atomicSwap.RedeemerAddress)
	if asset[:9] == "secondary" {
		return nil, fmt.Errorf("invalid asset: %s", asset)
	}
	tokenAddr := common.HexToAddress(string(asset[9:]))
	expiry, ok := new(big.Int).SetString(atomicSwap.Timelock, 10)
	if !ok {
		return nil, fmt.Errorf("invalid timelock: %s", atomicSwap.Timelock)
	}
	amount, ok := new(big.Int).SetString(atomicSwap.Amount, 10)
	if !ok {
		return nil, err
	}
	return ethereum.NewWatcher(initiatorAddr, redeemerAddr, deployerAddr, tokenAddr, secretHash, expiry, amount, client)
}
