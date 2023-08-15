package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/wbtc-garden/swapper"
	"go.uber.org/zap"
)

type initiatorSwap struct {
	logger         *zap.Logger
	initiator      *btcec.PrivateKey
	initiateTxHash string
	htlcScript     []byte
	waitBlocks     int64
	amount         uint64
	scriptAddr     btcutil.Address
	watcher        swapper.Watcher
	client         Client
}

func GetExpiry(goingFirst bool) int64 {
	if goingFirst {
		return 288
	}
	return 144
}

func NewInitiatorSwap(logger *zap.Logger, initiator *btcec.PrivateKey, redeemerAddr btcutil.Address, secretHash []byte, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.InitiatorSwap, error) {
	initiatorAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(initiator.PubKey().SerializeCompressed()), client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create initiator address: %w", err)
	}

	htlcScript, err := NewHTLCScript(initiatorAddr, redeemerAddr, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	watcher, err := NewWatcher(initiatorAddr, redeemerAddr, secretHash, waitBlocks, minConfirmations, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	if logger == nil {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}
	childLogger := logger.With(zap.String("service", "initSwap"))
	return &initiatorSwap{
		logger:     childLogger,
		initiator:  initiator,
		htlcScript: htlcScript,
		watcher:    watcher,
		scriptAddr: scriptAddr,
		amount:     amount,
		waitBlocks: waitBlocks,
		client:     client,
	}, nil
}

func (s *initiatorSwap) Initiate() (string, error) {
	txHash, err := s.client.Send(s.scriptAddr, s.amount, s.initiator)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	s.initiateTxHash = txHash
	s.logger.Info("Initiated", zap.String("txHash", txHash))
	return txHash, nil
}

func (s *initiatorSwap) Expired() (bool, error) {
	return s.watcher.Expired()
}

func (s *initiatorSwap) Refund() (string, error) {
	script := NewHTLCRefundWitness(s.initiator.PubKey().SerializeCompressed())
	txHash, err := s.client.Spend(s.htlcScript, script, s.initiator, uint(s.waitBlocks))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	s.logger.Info("Refunded", zap.String("txHash", txHash))
	return txHash, nil
}

func (s *initiatorSwap) WaitForRedeem() ([]byte, string, error) {
	for {
		s.logger.Debug("Waiting for Redemption", zap.String("address", s.scriptAddr.EncodeAddress()))
		redeemed, secret, tx, err := s.IsRedeemed()
		if redeemed {
			s.logger.Info("Redeemed", zap.String("secret", hex.EncodeToString(secret)))
			return secret, tx, nil
		}

		if err != nil {
			s.logger.Error("check redemption", zap.Error(err))
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *initiatorSwap) IsRedeemed() (bool, []byte, string, error) {
	return s.watcher.IsRedeemed()
}

type redeemerSwap struct {
	logger     *zap.Logger
	amount     uint64
	redeemer   *btcec.PrivateKey
	htlcScript []byte
	scriptAddr btcutil.Address
	watcher    swapper.Watcher
	client     Client
}

func NewRedeemerSwap(logger *zap.Logger, redeemer *btcec.PrivateKey, initiator btcutil.Address, secretHash []byte, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.RedeemerSwap, error) {
	redeemerAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(redeemer.PubKey().SerializeCompressed()), client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create redeemer address: %w", err)
	}

	htlcScript, err := NewHTLCScript(initiator, redeemerAddr, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	watcher, err := NewWatcher(initiator, redeemerAddr, secretHash, waitBlocks, minConfirmations, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Initialise the logger
	if logger == nil {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}
	childLogger := logger.With(zap.String("service", "redeemSwap"))

	return &redeemerSwap{
		logger:     childLogger,
		redeemer:   redeemer,
		watcher:    watcher,
		htlcScript: htlcScript,
		scriptAddr: scriptAddr,
		amount:     amount,
		client:     client,
	}, nil
}

func (s *redeemerSwap) Redeem(secret []byte) (string, error) {
	script := NewHTLCRedeemWitness(s.redeemer.PubKey().SerializeCompressed(), secret)
	txHash, err := s.client.Spend(s.htlcScript, script, s.redeemer, 0)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	s.logger.Info("Redeemed", zap.String("txHash", txHash))
	return txHash, nil
}

func (s *redeemerSwap) WaitForInitiate() ([]string, error) {
	for {
		initiated, txHashes, err := s.IsInitiated()
		if initiated {
			return txHashes, nil
		}
		if err != nil {
			s.logger.Error("wait for initiation", zap.Error(err))
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *redeemerSwap) IsInitiated() (bool, []string, error) {
	return s.watcher.IsInitiated()
}

type watcher struct {
	client           Client
	scriptAddr       btcutil.Address
	amount           uint64
	waitBlocks       int64
	minConfirmations uint64
}

func NewWatcher(initiator, redeemerAddr btcutil.Address, secretHash []byte, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.Watcher, error) {
	htlcScript, err := NewHTLCScript(initiator, redeemerAddr, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	return &watcher{scriptAddr: scriptAddr, amount: amount, waitBlocks: waitBlocks, minConfirmations: minConfirmations, client: client}, nil
}

func (w *watcher) Expired() (bool, error) {
	currentBlock, err := w.client.GetTipBlockHeight()
	if err != nil {
		return false, err
	}
	initiated, txHashes, err := w.IsInitiated()
	if err != nil || !initiated {
		return false, err
	}
	initiateBlockHeight, err := w.client.GetBlockHeight(txHashes[0])
	if err != nil {
		return false, err
	}
	return currentBlock-initiateBlockHeight+1 >= uint64(w.waitBlocks), nil
}

func (w *watcher) IsInitiated() (bool, []string, error) {
	utxos, bal, err := w.client.GetUTXOs(w.scriptAddr, 0)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if bal >= w.amount && len(utxos) > 0 {
		txHashes := make([]string, len(utxos))
		for i, utxo := range utxos {
			final, err := w.client.IsFinal(utxo.TxID, w.minConfirmations)
			if err != nil {
				return false, nil, fmt.Errorf("failed to check if final: %w", err)
			}
			if !final {
				return false, nil, nil
			}
			txHashes[i] = utxo.TxID
		}
		return true, txHashes, nil
	}
	return false, nil, nil
}

func (w *watcher) IsRedeemed() (bool, []byte, string, error) {
	witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	if err != nil {
		return false, nil, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if len(witness) == 5 {
		// inputs are [ 0 : sig, 1 : spender.PubKey().SerializeCompressed(),2 : secret, 3 :[]byte{0x1}, script]
		secretString := witness[2]
		secretBytes := make([]byte, hex.DecodedLen(len(secretString)))
		_, err := hex.Decode(secretBytes, []byte(secretString))
		if err != nil {
			return false, nil, "", fmt.Errorf("failed to decode secret: %w", err)
		}
		return true, secretBytes, tx, nil
	}
	return false, nil, "", nil
}

func (w *watcher) IsRefunded() (bool, string, error) {
	_, bal, err := w.client.GetUTXOs(w.scriptAddr, 0)
	if err != nil {
		return false, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	if err != nil {
		return false, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if len(witness) == 4 && bal == 0 {
		fmt.Println("Refunded:", witness)
		// inputs are [ 0 : sig, 1 : spender.PubKey().SerializeCompressed(), 2 :[]byte{}, script]
		return true, tx, nil

	}
	return false, "", nil
}
