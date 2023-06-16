package executor

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
)

type executor struct {
	bitcoinPrivateKey  *btcec.PrivateKey
	ethereumPrivateKey *ecdsa.PrivateKey
	client             bitcoin.Client
	ethereumClient     ethereum.Client
	wbtcAddress        common.Address
	store              Store
}

type Store interface {
	PendingTransactions() ([]model.Transaction, error)
	PutTransaction(tx model.Transaction) error
	UpdateTransaction(tx model.Transaction) error
}

type Config struct {
	Network *chaincfg.Params

	BitcoinURL  string
	EthereumURL string
	WBTCAddress string
}

type Executor interface {
	Run()
	rest.Swapper
}

func New(privateKey string, config Config, store Store) (Executor, error) {
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	btcPrivKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	ethPrivKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &executor{
		bitcoinPrivateKey:  btcPrivKey,
		ethereumPrivateKey: ethPrivKey,
		client:             bitcoin.NewClient(config.BitcoinURL, config.Network),
		ethereumClient:     ethereum.NewClient(config.EthereumURL),
		wbtcAddress:        common.HexToAddress(config.WBTCAddress),
		store:              store,
	}, nil
}

func (s *executor) Run() {
	for {
		txs, err := s.store.PendingTransactions()
		if err != nil {
			fmt.Println("Failed to retrieve transactions: ", err)
			time.Sleep(5 * time.Minute)
			continue
		}
		for _, tx := range txs {
			fmt.Println("Executing transaction: ", tx)
			s.execute(tx)
		}
		time.Sleep(10 * time.Second)
	}
}

func (s *executor) execute(tx model.Transaction) {
	initiatorSwap, err := s.getInitiatorSwap(tx.ToAddress, tx.SecretHash, tx.WBTCExpiry, deductFee(tx.Amount))
	if err != nil {
		panic(fmt.Errorf("constraint violation, this check should have been done before storage into DB: %v", err))
	}

	redeemerSwap, err := s.getRedeemerSwap(tx.FromAddress, tx.SecretHash, tx.WBTCExpiry, tx.Amount)
	if err != nil {
		panic(fmt.Errorf("constraint violation, this check should have been done before storage into DB: %v", err))
	}

	if tx.Status == 1 {
		txHash, err := initiatorSwap.Initiate()
		if err != nil {
			fmt.Println("Failed to initiate swap: ", err)
			return
		}
		tx.Status = 2
		tx.FollowerInitiateTxHash = txHash
		if err := s.store.UpdateTransaction(tx); err != nil {
			fmt.Println("Failed to update transaction: ", err)
			return
		}
	}

	if tx.Status == 2 {
		expired, err := initiatorSwap.Expired()
		if expired {
			tx.Status = 4
			if err := s.store.UpdateTransaction(tx); err != nil {
				fmt.Println("Failed to update transaction: ", err)
				return
			}
		} else if err != nil {
			fmt.Println("Failed to check if swap is expired: ", err)
			return
		} else {
			redeemed, secret, txHash, err := initiatorSwap.IsRedeemed()
			if err != nil {
				fmt.Println("Failed to check if swap is redeemed: ", err)
				return
			}

			fmt.Println("Redeemed: ", redeemed)
			if redeemed {
				tx.Status = 3
				tx.Secret = hex.EncodeToString(secret)
				tx.InitiatorRedeemTxHash = txHash
				if err := s.store.UpdateTransaction(tx); err != nil {
					fmt.Println("Failed to update transaction: ", err)
					return
				}
			}
			if err != nil {
				fmt.Println("Failed to check if swap is redeemed: ", err)
				return
			}
		}

	}

	if tx.Status == 3 {
		secret, err := hex.DecodeString(tx.Secret)
		if err != nil {
			panic(fmt.Errorf("constraint violation, this check should have been done before storage into DB: %v", err))
		}

		txHash, err := redeemerSwap.Redeem(secret)
		if err != nil {
			fmt.Println("Failed to redeem swap: ", err)
			return
		}
		tx.Status = 5
		tx.FollowerRedeemTxHash = txHash
		if err := s.store.UpdateTransaction(tx); err != nil {
			fmt.Println("Failed to update transaction: ", err)
			return
		}
	}

	if tx.Status == 4 {
		txHash, err := initiatorSwap.Refund()
		if err != nil {
			fmt.Println("Failed to refund swap: ", err)
			return
		}
		tx.Status = 6
		tx.FollowerRefundTxHash = txHash
		if err := s.store.UpdateTransaction(tx); err != nil {
			fmt.Println("Failed to update transaction: ", err)
			return
		}
	}
}

func (s *executor) GetAccount() (model.Account, error) {
	bitcoinAddress, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(s.bitcoinPrivateKey.PubKey().SerializeCompressed()), s.client.Net())
	if err != nil {
		return model.Account{}, err
	}

	ethereumAddress := crypto.PubkeyToAddress(s.ethereumPrivateKey.PublicKey)
	_, btcBalance, _ := s.client.GetUTXOs(bitcoinAddress, 0)
	ethBalance, err := s.ethereumClient.GetERC20Balance(s.wbtcAddress, ethereumAddress, nil)
	if err != nil {
		fmt.Println("Failed to get WBTC balance: ", err)
		return model.Account{}, err
	}
	return model.Account{
		BtcAddress:       bitcoinAddress.EncodeAddress(),
		BtcPubKey:        hex.EncodeToString(s.bitcoinPrivateKey.PubKey().SerializeCompressed()),
		WbtcAddress:      ethereumAddress.String(),
		BtcBalance:       strconv.FormatUint(btcBalance, 10),
		WbtcBalance:      ethBalance.String(),
		WbtcTokenAddress: s.wbtcAddress.Hex(),
		Fee:              0.1,
	}, nil
}

func (s *executor) ExecuteSwap(from, to, secretHash string, wbtcExpiry int64, amount uint64) error {
	_, err := s.getInitiatorSwap(from, secretHash, wbtcExpiry, deductFee(amount))
	if err != nil {
		return err
	}
	redeemer, err := s.getRedeemerSwap(to, secretHash, wbtcExpiry, amount)
	if err != nil {
		return err
	}

	initiated, txHash, err := redeemer.IsInitiated()
	if err != nil {
		return err
	}

	if !initiated {
		return fmt.Errorf("precondition violation: swap has not been initiated")
	}

	return s.store.PutTransaction(model.Transaction{
		FromAddress:             from,
		ToAddress:               to,
		SecretHash:              secretHash,
		WBTCExpiry:              wbtcExpiry,
		Amount:                  amount,
		Fee:                     amount - deductFee(amount),
		Status:                  1,
		InitiatorInitiateTxHash: txHash,
	})
}

func deductFee(amount uint64) uint64 {
	return amount * 999 / 1000
}

func (s *executor) decodeAddress(addr string) (interface{}, error) {
	if len(addr) == 40 {
		return common.HexToAddress("0x" + addr), nil
	}
	if len(addr) == 42 && strings.HasPrefix(addr, "0x") {
		return common.HexToAddress(addr), nil
	}
	return hex.DecodeString(addr)
}

func (s *executor) getInitiatorSwap(addr, secretHash string, block int64, amount uint64) (swapper.InitiatorSwap, error) {
	secretHashBytes, err := hex.DecodeString(secretHash)
	if err != nil {
		return nil, err
	}

	address, err := s.decodeAddress(addr)
	if err != nil {
		return nil, err
	}

	switch address := address.(type) {
	case []byte:
		return bitcoin.NewInitiatorSwap(s.bitcoinPrivateKey, address, secretHashBytes, 144, amount, s.client)
	case common.Address:
		return ethereum.NewInitiatorSwap(s.ethereumPrivateKey, address, s.wbtcAddress, secretHashBytes, big.NewInt(block), big.NewInt(int64(amount)), s.ethereumClient)
	default:
		return nil, fmt.Errorf("unknown address type")
	}
}

func (s *executor) getRedeemerSwap(addr, secretHash string, wbtcExpiry int64, amount uint64) (swapper.RedeemerSwap, error) {
	secretHashBytes, err := hex.DecodeString(secretHash)
	if err != nil {
		return nil, err
	}

	address, err := s.decodeAddress(addr)
	if err != nil {
		return nil, err
	}

	switch address := address.(type) {
	case []byte:
		return bitcoin.NewRedeemerSwap(s.bitcoinPrivateKey, address, secretHashBytes, 288, amount, s.client)
	case common.Address:
		return ethereum.NewRedeemerSwap(s.ethereumPrivateKey, address, s.wbtcAddress, secretHashBytes, big.NewInt(wbtcExpiry), big.NewInt(int64(amount)), s.ethereumClient)
	default:
		return nil, fmt.Errorf("unknown address type")
	}
}
