package watcher_test

import (
	"context"
	"errors"
	"time"

	"github.com/catalogfi/orderbook/mocks"
	"github.com/catalogfi/orderbook/model"
	"gorm.io/gorm"

	. "github.com/catalogfi/orderbook/watcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var _ = Describe("Watcher", func() {
	defer GinkgoRecover()

	logger := zap.NewNop()
	var (
		mockCtrl  *gomock.Controller
		mockStore *mocks.MockStore

		mockError = errors.New("mock error")
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockStore = mocks.NewMockStore(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("should be able to process an order", func() {
		It("should hard fail when initiator redeems and refunds", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedHard,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.FailedHard))
		})

		It("should hard fail when follower redeems and refunds", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedHard,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.FailedHard))
		})

		It("should soft fail when both initiator and follower refunds", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedSoft,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.FailedSoft))
		})

		It("should soft fail when both initiator refunds and follower does not start", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedSoft,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.FailedSoft))
		})

		It("should update order status to successful if both orders are redeemed", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder := model.Order{
				Status: model.Executed,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.Executed))
		})

		It("should update secret to be same as the follower swap's secret", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
					Secret: "secret",
				},
			}
			updatedOrder := model.Order{
				Status: model.Executed,
				Secret: "secret",
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
					Secret: "secret",
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.Executed))
			Expect(updatedOrder.Secret).To(Equal(updatedOrder.FollowerAtomicSwap.Secret))
		})

		It("should cancel an order if OrderTimeout passes", func() {
			createdAt := time.Now().Add(-OrderTimeout - time.Minute)
			order := model.Order{
				Status: model.Created,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder := model.Order{
				Status: model.Cancelled,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.Cancelled))
		})

		It("should not panic if the order status fails", func() {
			createdAt := time.Now().Add(-OrderTimeout - time.Minute)
			order := model.Order{
				Status: model.Created,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder := model.Order{
				Status: model.Cancelled,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.Cancelled))
		})

		It("should cancel an order if OrderTimeout passes", func() {
			createdAt := time.Now().Add(-SwapInitiationTimeout - time.Minute)
			order := model.Order{
				Status: model.Filled,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder := model.Order{
				Status: model.Cancelled,
				Model: gorm.Model{
					CreatedAt: createdAt,
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeTrue())
			Expect(updatedOrder.Status).To(Equal(model.Cancelled))
		})
		It("should not update order if the status is created and order does not expire", func() {
			order := model.Order{
				Status: model.Created,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}
			updatedOrder, ctn := ProcessOrder(order, mockStore, logger)
			Expect(ctn).To(BeFalse())
			Expect(updatedOrder.Status).To(Equal(model.Created))
		})
	})

	Describe("should be able to run the watcher", func() {
		defer GinkgoRecover()

		var minWorkers = 4

		It("should build a watcher", func() {
			watcher := NewWatcher(logger, mockStore, 1)
			Expect(watcher).ToNot(BeNil())
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedHard,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}

			mockStore.EXPECT().GetActiveOrders().Return([]model.Order{order}, nil).MaxTimes(2)
			mockStore.EXPECT().UpdateOrder(&updatedOrder).Return(nil)

			ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
			defer cancel()
			watcher.Run(ctx)
		})

		It("should build a watcher", func() {
			watcher := NewWatcher(logger, mockStore, minWorkers)
			Expect(watcher).ToNot(BeNil())
			order := model.Order{
				Status: model.Filled,
				Model: gorm.Model{
					CreatedAt: time.Now(),
				},
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.NotStarted,
				},
			}

			mockStore.EXPECT().GetActiveOrders().Return([]model.Order{order}, nil).MaxTimes(2)
			ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
			defer cancel()
			watcher.Run(ctx)
		})

		It("should build a watcher", func() {
			watcher := NewWatcher(logger, mockStore, minWorkers)
			Expect(watcher).ToNot(BeNil())
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			updatedOrder := model.Order{
				Status: model.FailedHard,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}

			mockStore.EXPECT().GetActiveOrders().Return([]model.Order{order}, nil).AnyTimes()
			mockStore.EXPECT().UpdateOrder(&updatedOrder).Return(mockError)
			ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
			watcher.Run(ctx)
		})

		It("should build a watcher", func() {
			watcher := NewWatcher(logger, mockStore, minWorkers)
			Expect(watcher).ToNot(BeNil())
			mockStore.EXPECT().GetActiveOrders().Return(nil, mockError).AnyTimes()
			ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
			watcher.Run(ctx)
		})

		It("should build a watcher", func() {
			order := model.Order{
				Status: model.Filled,
				InitiatorAtomicSwap: &model.AtomicSwap{
					Status: model.Refunded,
				},
				FollowerAtomicSwap: &model.AtomicSwap{
					Status: model.Redeemed,
				},
			}
			watcher := NewWatcher(logger, mockStore, 0)
			Expect(watcher).ToNot(BeNil())
			mockStore.EXPECT().GetActiveOrders().Return([]model.Order{order}, nil)
			ctx, _ := context.WithTimeout(context.Background(), 8*time.Second)
			watcher.Run(ctx)
		})
	})
})
