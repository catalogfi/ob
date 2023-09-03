package watcher_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/catalogfi/wbtc-garden/mocks"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	. "github.com/catalogfi/wbtc-garden/watcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/mock/gomock"
)

var _ = Describe("Ethereum Watcher", func() {
	defer GinkgoRecover()

	// logger := zap.NewNop()
	var (
		mockCtrl *gomock.Controller
		// mockWatcher   *mocks.MockWatcher
		mockStore *mocks.MockStore
		// mockBTCClient *mocks.MockBitcoinClient
		// mockScreener  *mocks.MockScreener

		mockError  = errors.New("mock error")
		mockTxHash = "mock tx hash"
		// mockAmount = "mock amount"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		// mockWatcher = mocks.NewMockWatcher(mockCtrl)
		mockStore = mocks.NewMockStore(mockCtrl)
		// mockBTCClient = mocks.NewMockBitcoinClient(mockCtrl)
		// mockScreener = mocks.NewMockScreener(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("can handle refund on EVM chains", func() {
		It("should fail if SwapByOCID fails", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMRefund(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{RefundTxHash: mockTxHash}, nil)
			err := HandleEVMRefund(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).Should(BeNil())
		})

		It("should update store if it is a new status", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			txhash := [32]byte{}
			rand.Read(txhash[:])
			txhashHash := common.BytesToHash(txhash[:])
			initialSwap := model.AtomicSwap{}
			updatedSwap := model.AtomicSwap{RefundTxHash: txhashHash.Hex(), Status: model.Refunded}
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(initialSwap, nil)
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err := HandleEVMRefund(mockStore, types.Log{TxHash: txhashHash, Topics: []common.Hash{{}, ocidHash}})
			Expect(err).Should(BeNil())
		})
	})

	Describe("can handle redeem on EVM chains", func() {
		It("should fail if SwapByOCID fails", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMRedeem(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{RedeemTxHash: mockTxHash}, nil)
			err := HandleEVMRedeem(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).Should(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{}, nil)
			err := HandleEVMRedeem(mockStore, types.Log{Data: []byte("hello"), Topics: []common.Hash{{}, ocidHash}})
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			txhash := [32]byte{}
			rand.Read(txhash[:])
			txhashHash := common.BytesToHash(txhash[:])

			secret := [32]byte{}
			rand.Read(secret[:])
			len := [32]byte{}
			len[31] = 0x20
			offset := [32]byte{}
			offset[31] = 0x20
			abiSecret := append(append(len[:], offset[:]...), secret[:]...)

			updatedSwap := model.AtomicSwap{RedeemTxHash: txhashHash.Hex(), Secret: hex.EncodeToString(secret[:]), Status: model.Redeemed}

			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{}, nil)
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)

			err := HandleEVMRedeem(mockStore, types.Log{TxHash: txhashHash, Data: abiSecret, Topics: []common.Hash{{}, ocidHash}})
			Expect(err).Should(BeNil())
		})

		// It("should update store if it is a new status", func() {
		// 	ocid := [32]byte{}
		// 	rand.Read(ocid[:])
		// 	ocidHash := common.BytesToHash(ocid[:])
		// 	txhash := [32]byte{}
		// 	rand.Read(txhash[:])
		// 	txhashHash := common.BytesToHash(txhash[:])
		// 	initialSwap := model.AtomicSwap{}
		// 	updatedSwap := model.AtomicSwap{RedeemTxHash: txhashHash.Hex(), Status: model.Redeemed}
		// 	mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(initialSwap, nil)
		// 	mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
		// 	err := HandleEVMRedeem(mockStore,  types.Log{TxHash: txhashHash, Topics: []common.Hash{{}, ocidHash}})
		// 	Expect(err).Should(BeNil())
		// })
	})
})
