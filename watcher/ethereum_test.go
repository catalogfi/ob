package watcher_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/catalogfi/orderbook/mocks"
	"github.com/catalogfi/orderbook/model"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	. "github.com/catalogfi/orderbook/watcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var _ = Describe("Ethereum Watcher", func() {
	defer GinkgoRecover()

	logger := zap.NewNop()
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
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMRefund(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{RefundTxHash: mockTxHash}, nil)
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
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(initialSwap, nil)
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
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMRedeem(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).ShouldNot(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{RedeemTxHash: mockTxHash}, nil)
			err := HandleEVMRedeem(mockStore, types.Log{Topics: []common.Hash{{}, ocidHash}})
			Expect(err).Should(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{}, nil)
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

			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{}, nil)
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
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{}, mockError)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if reedemer is incorrect", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{RedeemerAddress: "0xA1a547358A9Ca8E7b320d7742729e3334Ad96546", Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Redeemer: common.HexToAddress("0xA1a547368A9Ca8E7b320d7742729e3334Ad96546"), Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000), Expiry: big.NewInt(144)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should properly handle evm intiate", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Detected, RedeemerAddress: "0xA1a547358A9Ca8E7b320d7742729e3334Ad96546", Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144", InitiateTxHash: "0x0000000000000000000000000000000000000000000000000000000000000000"}).Return(nil)
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{RedeemerAddress: "0xA1a547358A9Ca8E7b320d7742729e3334Ad96546", Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Redeemer: common.HexToAddress("0xA1a547358A9Ca8E7b320d7742729e3334Ad96546"), Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000), Expiry: big.NewInt(144)}, mockScreener)
			Expect(err).Should(BeNil())
		})

		It("should not update store if it is already updated", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{InitiateTxHash: mockTxHash}, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{}, mockScreener)
			Expect(err).Should(BeNil())
		})

		It("should fail if the blacklisted check fails", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, mockError)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the blacklisted check returns true", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(true, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the swap amount is invalid", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "ffee"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890")}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the initiate amount is less than swap amount", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(99999)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the expiry is invalid", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "ffee"}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			err := HandleEVMInitiate(types.Log{Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000)}, mockScreener)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if the initiate timelock is not the swap timelock", func() {
			ocid := [32]byte{}
			rand.Read(ocid[:])
			ocidHash := common.BytesToHash(ocid[:])
			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144"}, nil)
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

			mockStore.EXPECT().SwapByOCID(ocidHash.Hex()[2:]).Return(model.AtomicSwap{RedeemerAddress: "0xA1a547358A9Ca8E7b320d7742729e3334Ad96546", Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144", InitiateTxHash: txhashHash.Hex(), InitiateBlockNumber: 100, OnChainIdentifier: ocidHash.Hex(), Status: model.Detected}, nil)
			// mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{"0x1234567890123456789012345678901234567890": model.EthereumSepolia}).Return(false, nil)
			// mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{RedeemerAddress: "0xA1a547358A9Ca8E7b320d7742729e3334Ad96546", Chain: model.EthereumSepolia, Amount: "100000", Timelock: "144", InitiateTxHash: txhashHash.Hex(), InitiateBlockNumber: 100, OnChainIdentifier: ocidHash.Hex(), Status: model.Detected})
			err := HandleEVMInitiate(types.Log{TxHash: txhashHash, BlockNumber: 100, Topics: []common.Hash{{}, ocidHash}}, mockStore, Swap{Initiator: common.HexToAddress("0x1234567890123456789012345678901234567890"), Amount: big.NewInt(100000), Expiry: big.NewInt(144), Redeemer: common.HexToAddress("0xA1a547358A9Ca8E7b320d7742729e3334Ad96546")}, mockScreener)
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
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Initiated, Timelock: "5000"}}, nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(BeNil())
		})

		It("should return nil if no swaps are found with status detected", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Detected, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 6, MinimumConfirmations: 6, Timelock: "5000"}).Return(nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(BeNil())
		})

		It("should return nil if no swaps are found with status detected and update swap is an err", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Detected, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 6, MinimumConfirmations: 6, Timelock: "5000"}).Return(mockError)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).Should(Not(BeNil()))
		})

		It("should return nil if order is expired and status is intitiated", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Expired, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}).Return(nil)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 6000)
			Expect(err).Should(BeNil())
		})
		It("should return err if order is expired and status is intitiated and update swap fails", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}}, nil)
			mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Expired, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6, Timelock: "5000"}).Return(mockError)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 6000)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if update order fails", func() {
			mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Detected, InitiateBlockNumber: 90, CurrentConfirmations: 5, MinimumConfirmations: 6}}, nil)
			// mockStore.EXPECT().UpdateSwap(&model.AtomicSwap{Status: model.Initiated, InitiateBlockNumber: 90, CurrentConfirmations: 6, MinimumConfirmations: 6}).Return(mockError)
			err := UpdateEVMConfirmations(mockStore, model.EthereumSepolia, 100)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Describe("creating new ethereum watcher", func() {
		It("Should succesfully new a ethereum watcher", func() {
			_, err := NewEthereumWatchers(mockStore, model.Config{
				Network: model.Network{
					model.EthereumSepolia: model.NetworkConfig{
						EventWindow: 1000,
						RPC: map[string]string{
							"ethrpc": "https://sepolia.infura.io/v3/68fc281b537345f2b9af8dfe4a72c75b",
						},
						Assets: map[model.Asset]model.Token{
							"0xA5E38d098b54C00F10e32E51647086232a9A0afD": {
								Oracle:       "https://api.coincap.io/v2/assets/bitcoin",
								TokenAddress: "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
								StartBlock:   18139000,
								Decimals:     8,
							},
						},
						Expiry: 7200,
					},
				},
			}, nil, logger)
			Expect(err).Should(BeNil())

		})

		It("Should fail if Asset is wrong", func() {
			_, err := NewEthereumWatchers(mockStore, model.Config{
				Network: model.Network{
					model.EthereumSepolia: model.NetworkConfig{
						EventWindow: 1000,
						RPC: map[string]string{
							"ethrp": "https://sepolia.infura.io/v3/68fc281b537345f2b9af8dfe4a72c75b",
						},
						Assets: map[model.Asset]model.Token{
							"0xA5E38d098b54C00F10e32E51647086232a9A0afD": {
								Oracle:       "https://api.coincap.io/v2/assets/bitcoin",
								TokenAddress: "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
								StartBlock:   18139000,
								Decimals:     8,
							},
						},
						Expiry: 7200,
					},
				},
			}, nil, logger)
			Expect(err).ShouldNot(BeNil())

		})
	})

	// Describe("Ethereum watch function", func(){
	// 	It("yi", func(){
	// 		ethWatcher, err := NewEthereumWatchers(mockStore, model.Config{
	// 			Network: model.Network{
	// 				model.EthereumSepolia: model.NetworkConfig{
	// 					EventWindow: 1000,
	// 					RPC: map[string]string{
	// 						"ethrpc": "https://sepolia.infura.io/v3/68fc281b537345f2b9af8dfe4a72c75b",
	// 					},
	// 					Assets: map[model.Asset]model.Token{
	// 						"0xA5E38d098b54C00F10e32E51647086232a9A0afD": {
	// 							Oracle:       "https://api.coincap.io/v2/assets/bitcoin",
	// 							TokenAddress: "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599",
	// 							StartBlock:   18139000,
	// 							Decimals:     8,
	// 						},
	// 					},
	// 					Expiry: 7200,
	// 				},
	// 			},
	// 		}, nil, logger)
	// 		Expect(err).Should(BeNil())
	// 		mockStore.EXPECT().GetActiveSwaps(model.EthereumSepolia).Return([]model.AtomicSwap{{Status: model.Initiated}}, nil).AnyTimes()
	// 		updatedSwap := model.AtomicSwap{RefundTxHash: "", Status: model.Initiated}
	// 		mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)

	// 		ethWatcher[0].Watch()
	// 	})
	// })
})
