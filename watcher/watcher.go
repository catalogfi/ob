package watcher

import (
	"context"
	"time"

	"github.com/catalogfi/orderbook/model"
	"go.uber.org/zap"
)

const SwapInitiationTimeout = 1 * time.Hour
const OrderTimeout = 3 * time.Minute

type Store interface {
	// UpdateOrder updates and order status in the db
	UpdateOrder(order *model.Order) error
	// GetActiveOrders fetches all orders which are active.
	GetActiveOrders() ([]model.Order, error)

	// get order by atomic swap id
	GetOrderBySwapID(swapID uint) (*model.Order, error)

	UpdateSwap(swap *model.AtomicSwap) error
	GetActiveSwaps(chain model.Chain) ([]model.AtomicSwap, error)
	SwapByOCID(ocID string) (model.AtomicSwap, error)
}

// Watcher watches the blockchain and update order status accordingly.
type Watcher interface {
	Run(ctx context.Context)
	RunWorker(ctx context.Context)
}

type watcher struct {
	logger  *zap.Logger
	store   Store
	workers int
	orders  chan model.Order
}

// NewWatcher returns a new Watcher
func NewWatcher(logger *zap.Logger, store Store, workers int) Watcher {
	return &watcher{
		logger:  logger.With(zap.String("service", "watcher")),
		store:   store,
		workers: workers,
		orders:  make(chan model.Order, 32),
	}
}

func (w *watcher) Run(ctx context.Context) {
	for i := 0; i < w.workers; i++ {
		childCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		go w.RunWorker(childCtx)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer close(w.orders)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(w.orders) == 0 {
				orders, err := w.store.GetActiveOrders()
				if err != nil {
					w.logger.Error("get active order", zap.Error(err))
					continue
				}
				for _, order := range orders {
					w.orders <- order
				}
			}
		}
	}
}

// allows for concurrent processing of orders
func (w *watcher) RunWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case order := <-w.orders:
			logger := w.logger.With(zap.Uint("order id", order.ID))
			order, hasUpdated := ProcessOrder(order, w.store, logger)
			if hasUpdated {
				if err := w.store.UpdateOrder(&order); err != nil {
					logger.Error("update order failed with", zap.Error(err))
				}
			}
		}
	}
}

func ProcessOrder(order model.Order, store Store, logger *zap.Logger) (model.Order, bool) {
	// copy secret from follower atomic swap
	secretUpdated := false
	if order.Secret != order.FollowerAtomicSwap.Secret {
		order.Secret = order.FollowerAtomicSwap.Secret
		order.SecretUpdatedAt = time.Now().UTC()
		secretUpdated = true
	}
	// Special cases regardless of order status
	switch {
	case order.Status == model.Created && time.Since(order.CreatedAt) > OrderTimeout:
		order.Status = model.Cancelled
	// Redeem and refund happens at the same time which should not happen. Defensive check.
	case (order.InitiatorAtomicSwap.Status == model.Redeemed && order.FollowerAtomicSwap.Status == model.Refunded) || (order.FollowerAtomicSwap.Status == model.Redeemed && order.InitiatorAtomicSwap.Status == model.Refunded):
		logger.Error("atomic swap hard failed as someone both redeemed and refunded")
		order.Status = model.FailedHard
	// Follower swap never gets initiated and initiator swap is initiated and refunded or both swaps are refunded
	case (order.InitiatorAtomicSwap.Status == model.Refunded && order.FollowerAtomicSwap.Status == model.NotStarted) || (order.InitiatorAtomicSwap.Status == model.Refunded && order.FollowerAtomicSwap.Status == model.Refunded):
		logger.Error("atomic swap soft failed due to initiator refunding")
		order.Status = model.FailedSoft
	// Follower has not filled the swap before the order timeout
	case (order.Status == model.Created && time.Since(order.CreatedAt) > OrderTimeout) || order.Status == model.Filled && time.Since(order.CreatedAt) > SwapInitiationTimeout && order.InitiatorAtomicSwap.Status == model.NotStarted:
		//check if the user partially filled the order
		if order.InitiatorAtomicSwap.FilledAmount != "" {
			break
		}
		logger.Info("atomic swap cancelled due to fill timeout")
		order.Status = model.Cancelled
	case (order.InitiatorAtomicSwap.Status == model.Redeemed && order.FollowerAtomicSwap.Status == model.Redeemed):
		logger.Info("atomic swap executed")
		order.Status = model.Executed
	}
	return order, (order.Status != model.Created && order.Status != model.Filled) || secretUpdated
}
