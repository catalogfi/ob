package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/susruth/wbtc-garden/swapper"
)

type initiatorSwap struct {
	initiator      *btcec.PrivateKey
	initiateTxHash string
	htlcScript     []byte
	waitBlocks     int64
	amount         uint64
	scriptAddr     btcutil.Address
	watcher        swapper.Watcher
	client         Client
}

func GetAddress(client Client, redeemerAddr, initiatorAddr btcutil.Address, secretHash []byte, waitBlocks int64) (btcutil.Address, error) {
	htlcScript, err := NewHTLCScript(initiatorAddr, redeemerAddr, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}
	return scriptAddr, nil
}

func GetAmount(client Client, redeemerAddr, initiatorAddr btcutil.Address, secretHash []byte, waitBlocks int64) (uint64, error) {
	scriptAddr, err := GetAddress(client, redeemerAddr, initiatorAddr, secretHash, waitBlocks)
	if err != nil {
		return 0, fmt.Errorf("failed to get script address: %w", err)
	}
	_, balance, err := client.GetUTXOs(scriptAddr, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get UTXOs: %w", err)
	}
	return balance, nil
}

func NewInitiatorSwap(initiator *btcec.PrivateKey, redeemerAddr btcutil.Address, secretHash []byte, waitBlocks int64, amount uint64, client Client) (swapper.InitiatorSwap, error) {
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

	fmt.Println("script address:", scriptAddr.EncodeAddress())

	watcher, err := NewWatcher(initiatorAddr, redeemerAddr, secretHash, waitBlocks, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	return &initiatorSwap{initiator: initiator, htlcScript: htlcScript, watcher: watcher, scriptAddr: scriptAddr, amount: amount, waitBlocks: waitBlocks, client: client}, nil
}

func (s *initiatorSwap) Initiate() (string, error) {
	txHash, err := s.client.Send(s.scriptAddr, s.amount, s.initiator)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	s.initiateTxHash = txHash
	fmt.Println("Successfully Initiated:", txHash)
	return txHash, nil
}

func (initiatorSwap *initiatorSwap) Expired() (bool, error) {
	return initiatorSwap.watcher.Expired()
}

func (s *initiatorSwap) Refund() (string, error) {
	script := NewHTLCRefundScript(s.initiator.PubKey().SerializeCompressed())
	txHash, err := s.client.Spend(s.htlcScript, script, s.initiator, uint(s.waitBlocks))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	fmt.Println("Successfully Refunded:", txHash)
	return txHash, nil
}

func (s *initiatorSwap) WaitForRedeem() ([]byte, string, error) {
	for {
		fmt.Println("Waiting for Redemption on:", s.scriptAddr)
		redeemed, secret, tx, err := s.IsRedeemed()
		if redeemed {
			return secret, tx, nil
		}

		if err != nil {
			fmt.Println("failed to check if redeemed:", err)
		}
		time.Sleep(5 * time.Second)
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

func NewRedeemerSwap(redeemer *btcec.PrivateKey, initiator btcutil.Address, secretHash []byte, waitTime int64, amount uint64, client Client) (swapper.RedeemerSwap, error) {
	redeemerAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(redeemer.PubKey().SerializeCompressed()), client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create redeemer address: %w", err)
	}

	htlcScript, err := NewHTLCScript(initiator, redeemerAddr, secretHash, waitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	fmt.Println("script address:", scriptAddr.EncodeAddress())
	watcher, err := NewWatcher(initiator, redeemerAddr, secretHash, waitTime, amount, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}
	return &redeemerSwap{redeemer: redeemer, watcher: watcher, htlcScript: htlcScript, scriptAddr: scriptAddr, amount: amount, client: client}, nil
}

func (s *redeemerSwap) Redeem(secret []byte) (string, error) {
	script := NewHTLCRedeemScript(s.redeemer.PubKey().SerializeCompressed(), secret)
	txHash, err := s.client.Spend(s.htlcScript, script, s.redeemer, 0)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	fmt.Println("Successfully Redeemed:", txHash)
	return txHash, nil
}

func (s *redeemerSwap) WaitForInitiate() (string, error) {
	for {
		initiated, txHash, err := s.IsInitiated()
		if initiated {
			return txHash, nil
		}
		if err != nil {
			fmt.Println("failed to check if initiated:", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *redeemerSwap) IsInitiated() (bool, string, error) {
	return s.watcher.IsInitiated()
}

type watcher struct {
	client     Client
	scriptAddr btcutil.Address
	amount     uint64
	waitBlocks int64
}

func NewWatcher(initiator, redeemerAddr btcutil.Address, secretHash []byte, waitBlocks int64, amount uint64, client Client) (swapper.Watcher, error) {
	htlcScript, err := NewHTLCScript(initiator, redeemerAddr, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	witnessProgram := sha256.Sum256(htlcScript)
	scriptAddr, err := btcutil.NewAddressWitnessScriptHash(witnessProgram[:], client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	fmt.Println("script address:", scriptAddr.EncodeAddress())
	return &watcher{scriptAddr: scriptAddr, amount: amount, waitBlocks: waitBlocks, client: client}, nil
}

func (w *watcher) Expired() (bool, error) {
	currentBlock, err := w.client.GetTipBlockHeight()
	if err != nil {
		return false, err
	}
	initiated, txHash, err := w.IsInitiated()
	if err != nil || !initiated {
		return false, err
	}
	initiateBlockHeight, err := w.client.GetBlockHeight(txHash)
	if err != nil {
		return false, err
	}
	if currentBlock > initiateBlockHeight+uint64(w.waitBlocks) {
		return true, nil
	}
	return false, nil
}

func (w *watcher) IsInitiated() (bool, string, error) {
	utxos, bal, err := w.client.GetUTXOs(w.scriptAddr, 0)
	if err != nil {
		return false, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if bal >= w.amount && len(utxos) > 0 {
		final, err := w.client.IsFinal(utxos[0].TxID)
		if err != nil {
			return false, "", fmt.Errorf("failed to check if final: %w", err)
		}
		if !final {
			return false, utxos[0].TxID, nil
		}
		return true, utxos[0].TxID, nil
	}
	return false, "", nil
}

func (w *watcher) IsRedeemed() (bool, []byte, string, error) {
	witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	if err != nil {
		return false, nil, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if len(witness) != 0 {
		fmt.Println("Redeemed:", witness)
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
