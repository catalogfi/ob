package watcher_test

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/wbtc-garden/mocks"
	"github.com/catalogfi/wbtc-garden/model"
	"github.com/catalogfi/wbtc-garden/swapper/bitcoin"

	. "github.com/catalogfi/wbtc-garden/watcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var _ = Describe("Bitcoin Watcher", func() {
	defer GinkgoRecover()

	logger := zap.NewNop()
	var (
		mockCtrl      *gomock.Controller
		mockWatcher   *mocks.MockWatcher
		mockStore     *mocks.MockStore
		mockBTCStore  *mocks.MockBTCStore
		mockBTCClient *mocks.MockBitcoinClient
		mockScreener  *mocks.MockScreener

		mockError  = errors.New("mock error")
		mockTxHash = "mock tx hash"
		// mockAmount = "mock amount"
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockWatcher = mocks.NewMockWatcher(mockCtrl)
		mockStore = mocks.NewMockStore(mockCtrl)
		mockBTCStore = mocks.NewMockBTCStore(mockCtrl)
		mockBTCClient = mocks.NewMockBitcoinClient(mockCtrl)
		mockScreener = mocks.NewMockScreener(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("can load btc client", func() {
		It("with blockstream url", func() {
			btcClient, err := LoadBTCClient(model.BitcoinTestnet, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, nil)
			Expect(err).Should(BeNil())
			Expect(btcClient).Should(Not(BeNil()))
		})

		It("should fail if no RPCs are provided", func() {
			_, err := LoadBTCClient(model.BitcoinTestnet, model.NetworkConfig{}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("with blockstream url should create an instant wallet", func() {
			btcClient, err := LoadBTCClient(model.BitcoinTestnet, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).Should(BeNil())
			Expect(btcClient).Should(Not(BeNil()))
		})
	})

	Describe("can load btc watcher", func() {
		It("should fail if no RPCs are provided", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain: model.BitcoinTestnet,
			}, model.NetworkConfig{}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid btc address is provided for the initiator", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "invalid address",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid btc address is provided for the redeemer", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "invalid address",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid amount is provided", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "invalid hash",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid secret hash is provided", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:           "ffee",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid timelock is provided", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:           "100000",
				Timelock:         "ffee",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})

		It("should fail if an invalid timelock is provided", func() {
			_, err := LoadBTCWatcher(model.AtomicSwap{
				Chain:            model.BitcoinTestnet,
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:           "100000",
				Timelock:         "ffee",
			}, model.NetworkConfig{
				RPC: map[string]string{
					"blockstream": "https://blockstream.info/testnet/api",
				},
			}, mockBTCStore)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Describe("can build and run the btc watcher", func() {
		It("should fail if ProcessSwaps fails", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{}, nil, time.Second, logger)
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return(nil, mockError).MaxTimes(2)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			btcWatcher.Watch(ctx)
		})

		It("should fail if get active orders fails", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{}, nil, time.Second, logger)
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return(nil, mockError)
			err := btcWatcher.ProcessBTCSwaps()
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the rpc url is invalid", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{
				Network: model.Network{
					model.BitcoinTestnet: model.NetworkConfig{
						RPC: map[string]string{
							"blockcypher": "fake url",
						},
					},
				},
			}, nil, time.Second, logger)
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return([]model.AtomicSwap{{
				Chain: model.BitcoinTestnet,
			}}, nil)
			err := btcWatcher.ProcessBTCSwaps()
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the chain is not btc", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{
				Network: model.Network{
					model.EthereumSepolia: model.NetworkConfig{
						RPC: map[string]string{
							"ethrpc": "https://gateway.tenderly.co/public/sepolia",
						},
					},
				},
			}, nil, time.Second, logger)
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return([]model.AtomicSwap{{
				Chain: model.EthereumSepolia,
			}}, nil)
			err := btcWatcher.ProcessBTCSwaps()
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if we fail to load the watcher", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{
				Network: model.Network{
					model.BitcoinTestnet: model.NetworkConfig{
						RPC: map[string]string{
							"mempool": "https://mempool.space/testnet/api",
						},
					},
				},
			}, nil, time.Second, logger)
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return([]model.AtomicSwap{{
				Chain: model.BitcoinTestnet,
			}}, nil)
			err := btcWatcher.ProcessBTCSwaps()
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the update swap fails", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{
				Network: model.Network{
					model.BitcoinTestnet: model.NetworkConfig{
						RPC: map[string]string{
							"mempool": "https://mempool.space/testnet/api",
						},
					},
				},
			}, nil, time.Second, logger)

			initialSwap := model.AtomicSwap{
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:           "100000",
				Timelock:         "144",
				Chain:            model.BitcoinTestnet,
			}

			updatedSwap := model.AtomicSwap{
				InitiatorAddress:  "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:   "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				OnChainIdentifier: "tb1qy88x434wze5keyakatgd9geqp74td3zqlxs7zcd0q45twrrdgjgqk8sr9e",
				SecretHash:        "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:            "100000",
				FilledAmount:      "0",
				Timelock:          "144",
				Chain:             model.BitcoinTestnet,
			}
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return([]model.AtomicSwap{initialSwap}, nil)
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(mockError)
			Expect(btcWatcher.ProcessBTCSwaps()).Should(Not(BeNil()))
		})

		It("should fail if the update swap fails", func() {
			btcWatcher := NewBTCWatcher(mockStore, mockBTCStore, model.BitcoinTestnet, model.Config{
				Network: model.Network{
					model.BitcoinTestnet: model.NetworkConfig{
						RPC: map[string]string{
							"mempool": "https://mempool.space/testnet/api",
						},
					},
				},
			}, nil, time.Second, logger)

			initialSwap := model.AtomicSwap{
				InitiatorAddress: "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:  "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				SecretHash:       "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:           "100000",
				Timelock:         "144",
				Chain:            model.BitcoinTestnet,
			}

			updatedSwap := model.AtomicSwap{
				InitiatorAddress:  "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd",
				RedeemerAddress:   "n1g3aBR4dhZnhwzT4PhfoaVB2doYJ5JvpX",
				OnChainIdentifier: "tb1qy88x434wze5keyakatgd9geqp74td3zqlxs7zcd0q45twrrdgjgqk8sr9e",
				SecretHash:        "0011223344556677889900112233445566778899001122334455667788990011",
				Amount:            "100000",
				FilledAmount:      "0",
				Timelock:          "144",
				Chain:             model.BitcoinTestnet,
			}
			mockStore.EXPECT().GetActiveSwaps(model.BitcoinTestnet).Return([]model.AtomicSwap{initialSwap}, nil)
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			Expect(btcWatcher.ProcessBTCSwaps()).Should(BeNil())
		})
	})

	Describe("can update swap status", func() {
		defer GinkgoRecover()

		It("should fail if BTCInitiateStatus check fails", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &model.AtomicSwap{OnChainIdentifier: ""})
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if BTCInitiateStatus check fails", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockWatcher.EXPECT().Identifier().Return("tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v")
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			err = UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "ffee"})
			Expect(err).Should(Not(BeNil()))
		})

		It("should update filled amount if a tx is detected", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockWatcher.EXPECT().Identifier().Return("tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v")
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "120000"}
			updatedSwap := model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "120000", FilledAmount: "100000", InitiateTxHash: "txHash1"}
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err = UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should update status to detected if the amount is >= swap amount", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockWatcher.EXPECT().Identifier().Return("tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v")
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "100000"}
			updatedSwap := model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "100000", FilledAmount: "100000", InitiateTxHash: "txHash1", Status: model.Detected}
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err = UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should update status to detected even if the deposit happened over multiple txs", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockWatcher.EXPECT().Identifier().Return("tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v").Times(2)
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}, {
				Amount: 100000,
				TxID:   "txHash2",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 102,
				},
			}}, uint64(200000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			mockBTCClient.EXPECT().GetTx("txHash2").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", Amount: "200000"}
			updatedSwap := model.AtomicSwap{OnChainIdentifier: mockAddress, Amount: "200000", FilledAmount: "200000", InitiateTxHash: "txHash1,txHash2", Status: model.Detected}
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err = UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should fail if get btc confirmations fails", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetConfirmations(mockTxHash).Return(uint64(0), uint64(0), mockError)

			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, Status: model.Detected})
			Expect(err).Should(Not(BeNil()))
		})

		It("should update the current confirmations if it is different from existing data", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetConfirmations(mockTxHash).Return(uint64(100), uint64(4), nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", MinimumConfirmations: 6, InitiateTxHash: mockTxHash, Status: model.Detected}
			updatedSwap := model.AtomicSwap{OnChainIdentifier: "", MinimumConfirmations: 6, InitiateTxHash: mockTxHash, Status: model.Detected, CurrentConfirmations: 4}
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should update the status to intiated after crossing the required number of confirmations", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetConfirmations("txHash1").Return(uint64(100), uint64(14), nil)
			mockBTCClient.EXPECT().GetConfirmations("txHash2").Return(uint64(102), uint64(12), nil)
			mockBTCClient.EXPECT().GetConfirmations("txHash3").Return(uint64(101), uint64(13), nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", MinimumConfirmations: 10, InitiateTxHash: "txHash1,txHash2,txHash3", Status: model.Detected}
			updatedSwap := model.AtomicSwap{OnChainIdentifier: "", MinimumConfirmations: 10, InitiateTxHash: "txHash1,txHash2,txHash3", Status: model.Initiated, CurrentConfirmations: 10, InitiateBlockNumber: 102}
			mockStore.EXPECT().UpdateSwap(&updatedSwap).Return(nil)
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should fail if unable to get tip block height", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(0), mockError)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, Status: model.Initiated}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if timelock is not a decimal number", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, Status: model.Initiated, Timelock: "ffee"}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if watcher's IsRedeemed returns an error", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRedeemed().Return(false, []byte{}, "", mockError)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, InitiateBlockNumber: 9999, Timelock: "1000", Status: model.Initiated}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(Not(BeNil()))
		})

		It("should return and not update swap if swap is not redeemed yet", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRedeemed().Return(false, []byte{}, "", nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, InitiateBlockNumber: 9999, Timelock: "1000", Status: model.Initiated}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should add the secret and the redeem tx hash if is redeemed succeeds", func() {
			secret := [32]byte{}
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRedeemed().Return(true, secret[:], mockTxHash, nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, InitiateBlockNumber: 9999, Timelock: "1000", Status: model.Initiated}
			finalSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateTxHash: mockTxHash, InitiateBlockNumber: 9999, Timelock: "1000", Secret: hex.EncodeToString(secret[:]), RedeemTxHash: mockTxHash, Status: model.Redeemed}
			mockStore.EXPECT().UpdateSwap(&finalSwap).Return(nil)
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should fail if watcher's IsRefunded returns an error", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRefunded().Return(false, "", mockError)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateBlockNumber: 5000, Timelock: "4999", InitiateTxHash: mockTxHash, Status: model.Initiated}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(Not(BeNil()))
		})

		It("should return and not update swap if swap is not redeemed yet", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRefunded().Return(false, "", nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateBlockNumber: 5000, Timelock: "4999", InitiateTxHash: mockTxHash, Status: model.Initiated}
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})

		It("should add the secret and the redeem tx hash if is redeemed succeeds", func() {
			mockWatcher.EXPECT().Identifier().Return("")
			mockBTCClient.EXPECT().GetTipBlockHeight().Return(uint64(10000), nil)
			mockWatcher.EXPECT().IsRefunded().Return(true, mockTxHash, nil)
			initialSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateBlockNumber: 5000, Timelock: "4999", InitiateTxHash: mockTxHash, Status: model.Initiated}
			finalSwap := model.AtomicSwap{OnChainIdentifier: "", InitiateBlockNumber: 5000, Timelock: "4999", InitiateTxHash: mockTxHash, RefundTxHash: mockTxHash, Status: model.Refunded}
			mockStore.EXPECT().UpdateSwap(&finalSwap).Return(nil)
			err := UpdateSwapStatus(mockWatcher, mockBTCClient, nil, mockStore, &initialSwap)
			Expect(err).Should(BeNil())
		})
	})

	Describe("get btc confirmations should work in all cases", func() {
		defer GinkgoRecover()

		It("should fail if the btc client fails to get confirmations", func() {
			mockBTCClient.EXPECT().GetConfirmations(mockTxHash).Return(uint64(0), uint64(0), mockError)
			height, conf, err := GetBTCConfirmations(mockBTCClient, mockTxHash)
			Expect(height).Should(Equal(uint64(0)))
			Expect(conf).Should(Equal(uint64(0)))
			Expect(err).Should(Not(BeNil()))
		})

		It("should successfully return block height and confirmations when there is one utxo", func() {
			mockBTCClient.EXPECT().GetConfirmations(mockTxHash).Return(uint64(100), uint64(4), nil)
			height, conf, err := GetBTCConfirmations(mockBTCClient, mockTxHash)
			Expect(height).Should(Equal(uint64(100)))
			Expect(conf).Should(Equal(uint64(4)))
			Expect(err).Should(BeNil())
		})

		It("should successfully return block height and confirmations when there are multiple utxos", func() {
			mockTxHashes := "txHash1,txHash2,txHash3"
			mockBTCClient.EXPECT().GetConfirmations("txHash1").Return(uint64(100), uint64(4), nil)
			mockBTCClient.EXPECT().GetConfirmations("txHash2").Return(uint64(102), uint64(2), nil)
			mockBTCClient.EXPECT().GetConfirmations("txHash3").Return(uint64(101), uint64(3), nil)
			height, conf, err := GetBTCConfirmations(mockBTCClient, mockTxHashes)
			Expect(height).Should(Equal(uint64(102)))
			Expect(conf).Should(Equal(uint64(2)))
			Expect(err).Should(BeNil())
		})
	})

	Describe("can check bitcoin initiate status", func() {
		defer GinkgoRecover()

		It("should fail if the script address is invalid", func() {
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			_, _, err := BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, "")
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the script address is invalid", func() {
			mockAddress := "mockAddress"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			_, _, err := BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the get utxos fails", func() {
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())
			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{}, uint64(0), mockError)
			_, _, err = BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the client fails to get the tx details", func() {
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{}, mockError)

			_, _, err = BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(Not(BeNil()))
		})

		It("should successfully return balance and txhash if the screener is nil", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)

			bal, txHash, err := BTCInitiateStatus(mockBTCClient, nil, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(BeNil())
			Expect(bal).Should(Equal(uint64(100000)))
			Expect(txHash).Should(Equal("txHash1"))
		})

		It("should successfully return balance and txhash if the address is not blacklisted", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{depositAddr: model.BitcoinTestnet}).Return(false, nil)
			bal, txHash, err := BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(BeNil())
			Expect(bal).Should(Equal(uint64(100000)))
			Expect(txHash).Should(Equal("txHash1"))
		})

		It("should fail if the depositor address is blacklisted", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{depositAddr: model.BitcoinTestnet}).Return(true, nil)
			_, _, err = BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(Not(BeNil()))
		})

		It("should fail if the screener fails to check if the account is blacklisted", func() {
			depositAddr := "n2psi3r4BpvzjPPXdaz3de1k1MgNi4Wyzd"
			mockAddress := "tb1qdcsqrldj6xhapxq55533j028dkyvsyc53w5gzkuys8dzng08y5jsthrz8v"
			mockBTCClient.EXPECT().Net().Return(&chaincfg.TestNet3Params)
			mockAddr, err := btcutil.DecodeAddress(mockAddress, &chaincfg.TestNet3Params)
			Expect(err).Should(BeNil())

			mockBTCClient.EXPECT().GetUTXOs(mockAddr, uint64(0)).Return(bitcoin.UTXOs{{
				Amount: 100000,
				TxID:   "txHash1",
				Vout:   0,
				Status: &bitcoin.Status{
					Confirmed:   false,
					BlockHeight: 100,
				},
			}}, uint64(100000), nil)
			mockBTCClient.EXPECT().GetTx("txHash1").Return(bitcoin.Transaction{VINs: []bitcoin.VIN{{Prevout: bitcoin.Prevout{ScriptPubKeyAddress: depositAddr}}}}, nil)
			mockScreener.EXPECT().IsBlacklisted(map[string]model.Chain{depositAddr: model.BitcoinTestnet}).Return(false, mockError)
			_, _, err = BTCInitiateStatus(mockBTCClient, mockScreener, model.BitcoinTestnet, mockAddress)
			Expect(err).Should(Not(BeNil()))
		})
	})
})
