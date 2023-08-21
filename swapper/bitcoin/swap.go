package bitcoin

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/wbtc-garden/swapper"
)

type initiatorSwap struct {
	watcher swapper.Watcher
	client  Client

	initiator  *btcec.PrivateKey
	waitBlocks int64
	amount     uint64
	script     []byte
	scriptAddr btcutil.Address
}

// TODO : Naming is very confusing, Bob (Buy BTC/ Sell WBTC) needs to create a initiator swap for WBTC and a redeemer swap
// for btc ???
func NewInitiatorSwap(initiator *btcec.PrivateKey, redeemerAddr btcutil.Address, secretHash []byte, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.InitiatorSwap, error) {
	initiatorAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(initiator.PubKey().SerializeCompressed()), client.Net())
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

	watcher, err := NewWatcher(scriptAddr, waitBlocks, minConfirmations, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &initiatorSwap{
		initiator:  initiator,
		script:     htlcScript,
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
	return txHash, nil
}

func (s *initiatorSwap) Expired() (bool, error) {
	return s.watcher.Expired()
}

func (s *initiatorSwap) Refund() (string, error) {
	witness := NewHTLCRefundWitness(s.initiator.PubKey().SerializeCompressed())
	txHash, err := s.client.Spend(s.script, witness, s.initiator, uint(s.waitBlocks))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	return txHash, nil
}

func (s *initiatorSwap) WaitForRedeem() ([]byte, string, error) {
	for {
		redeemed, secret, tx, err := s.IsRedeemed()
		if err != nil {
			time.Sleep(30 * time.Second)
			continue
		}

		if redeemed {
			return secret, tx, nil
		}
		time.Sleep(30 * time.Second)
	}
}

func (s *initiatorSwap) IsRedeemed() (bool, []byte, string, error) {
	return s.watcher.IsRedeemed()
}

type redeemerSwap struct {
	amount     uint64
	redeemer   *btcec.PrivateKey
	htlcScript []byte
	scriptAddr btcutil.Address
	watcher    swapper.Watcher
	client     Client
}

func NewRedeemerSwap(redeemer *btcec.PrivateKey, initiator btcutil.Address, secretHash []byte, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.RedeemerSwap, error) {
	redeemerAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(redeemer.PubKey().SerializeCompressed()), client.Net())
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

	watcher, err := NewWatcher(scriptAddr, waitBlocks, minConfirmations, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &redeemerSwap{
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
	return txHash, nil
}

func (s *redeemerSwap) WaitForInitiate() (string, error) {
	for {
		initiated, txHash, _, _ := s.IsInitiated()
		if initiated {
			return txHash, nil
		}
		time.Sleep(30 * time.Second)
	}
}

func (s *redeemerSwap) IsInitiated() (bool, string, uint64, error) {
	return s.watcher.IsInitiated()
}
