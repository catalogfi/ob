package watcher_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"

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
		mockScreener *mocks.MockScreener

		mockError  = errors.New("mock error")
		mockTxHash = "mock tx hash"
		// mockAmount = "mock amount"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		// mockWatcher = mocks.NewMockWatcher(mockCtrl)
		mockStore = mocks.NewMockStore(mockCtrl)
		// mockBTCClient = mocks.NewMockBitcoinClient(mockCtrl)
		mockScreener = mocks.NewMockScreener(mockCtrl)
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
	})

	Describe("can handle initiate on EVM chains", func() {
		It("should fail if SwapByOCID fails", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{InitiateTxHash: mockTxHash}, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{}, mockScreener)
			Expect(err).Should(BeNil())
		})

		It("should fail if the blacklisted check fails", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, mockError)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the blacklisted check returns true", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(true, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the swap amount is invalid", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "ffee"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the initiate amount is less than swap amount", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(99999)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the expiry is invalid", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "ffee"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the initiate timelock is not the swap timelock", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000), Expiry: big.NewInt(12)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should update the swap if initiate is valid", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])

			txhash := [32]byte{}
			rand.Read(txhash[:])
			txhashHash := common.BytesToHash(txhash[:])

			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144", InitiateTxHash: txhashHash.Hex(), InitiateBlockNumber: 100, OnChainIdentifier: ocidHash.Hex(), Status: model.Detected})
			err := HandleEVMInitiate(types.Log{TxHash: txhashHash, BlockNumber: 100, Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000), Expiry: big.NewInt(144)}, mockScreener)
			Expect(err).Should(BeNil())
		})
	})

	Describe("can update EVM confirmations", func() {
		It("should fail if get active swaps fails", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return(nil, mockError)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).ShouldNot(BeNil())
		})

		It("should return nil if no swaps are found", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return(nil, nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(BeNil())
		})

		It("should return nil if no swaps are found with status detected", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Initiated}}, nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(BeNil())
		})

		It("should return nil if no swaps are found with status detected", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Detected, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 6, MinimumConfirmations: 6}).Return(nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(BeNil())
		})

		It("should fail if update order fails", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Detected, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 6, MinimumConfirmations: 6}).Return(mockError)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).ShouldNot(BeNil())
		})
	})
})
