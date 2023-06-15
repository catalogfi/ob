package bitcoin

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/susruth/wbtc-garden-server/swapper"
)

type initiatorSwap struct {
	initiator             *btcec.PrivateKey
	initiateTxHash        string
	initiateTxBlockHeight uint64
	htlcScript            []byte
	waitBlocks            int64
	scriptAddr            btcutil.Address
	client                Client
}

func NewInitiatorSwap(initiator *btcec.PrivateKey, redeemerPublicKey, secretHash []byte, waitBlocks int64, amount uint64, client Client) (swapper.InitiatorSwap, error) {
	htlcScript, err := NewHTLCScript(initiator.PubKey().SerializeCompressed(), redeemerPublicKey, secretHash, waitBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	scriptAddr, err := btcutil.NewAddressScriptHash(htlcScript, client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	fmt.Println("script address:", scriptAddr.EncodeAddress())
	return &initiatorSwap{initiator: initiator, htlcScript: htlcScript, scriptAddr: scriptAddr, client: client}, nil
}

func (s *initiatorSwap) Initiate() (string, error) {
	txHash, err := s.client.Send(s.scriptAddr, 100000, s.initiator)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	s.initiateTxHash = txHash
	fmt.Println("Successfully Initiated:", txHash)
	return txHash, nil
}

func (initiatorSwap *initiatorSwap) Expired() (bool, error) {
	currentBlock, err := initiatorSwap.client.GetTipBlockHeight()
	if err != nil {
		return false, err
	}

	initiateBlockHeight, err := initiatorSwap.client.GetBlockHeight(initiatorSwap.initiateTxHash)
	if err != nil {
		return false, err
	}

	expiryBlockHeight := initiateBlockHeight + uint64(initiatorSwap.waitBlocks)
	if currentBlock > expiryBlockHeight {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *initiatorSwap) Refund() (string, error) {
	script, err := NewHTLCRefundScript()
	if err != nil {
		return "", fmt.Errorf("failed to create redeem script: %w", err)
	}
	txHash, err := s.client.Spend(s.htlcScript, script, s.initiator)
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
	scriptSig, tx, err := s.client.GetSpendingScriptSig(s.scriptAddr)
	if err != nil {
		return false, nil, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if scriptSig != "" {
		fmt.Println("Redeemed:", scriptSig)
		secret, err := hex.DecodeString(strings.Split(scriptSig, " ")[3])
		if err != nil {
			return false, nil, "", fmt.Errorf("failed to decode %s", strings.Split(scriptSig, " ")[3])
		}
		return true, secret, tx, nil
	}
	if err := s.checkTimeout(); err != nil {
		return false, nil, "", err
	}
	return false, nil, "", nil
}

func (s *initiatorSwap) checkTimeout() error {
	if s.initiateTxBlockHeight == 0 {
		blockHeight, err := s.client.GetBlockHeight(s.initiateTxHash)
		if err != nil {
			return fmt.Errorf("failed to get block height: %w", err)
		}
		if blockHeight != 0 {
			s.initiateTxBlockHeight = blockHeight
		}
	} else {
		tipBlockHeight, err := s.client.GetTipBlockHeight()
		if err != nil {
			return fmt.Errorf("failed to get block height: %w", err)
		}
		if tipBlockHeight-s.initiateTxBlockHeight > uint64(s.waitBlocks)+1 {
			return swapper.ErrRedeemTimeout
		}
	}
	return nil
}

type redeemerSwap struct {
	amount     uint64
	redeemer   *btcec.PrivateKey
	htlcScript []byte
	scriptAddr btcutil.Address
	client     Client
}

func NewRedeemerSwap(redeemer *btcec.PrivateKey, initiatorPublicKeyBytes, secretHash []byte, waitTime int64, amount uint64, client Client) (swapper.RedeemerSwap, error) {
	htlcScript, err := NewHTLCScript(initiatorPublicKeyBytes, redeemer.PubKey().SerializeCompressed(), secretHash, waitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTLC script: %w", err)
	}
	scriptAddr, err := btcutil.NewAddressScriptHash(htlcScript, client.Net())
	if err != nil {
		return nil, fmt.Errorf("failed to create script address: %w", err)
	}

	fmt.Println("script address:", scriptAddr.EncodeAddress())

	return &redeemerSwap{redeemer: redeemer, htlcScript: htlcScript, scriptAddr: scriptAddr, amount: amount, client: client}, nil
}

func (s *redeemerSwap) Redeem(secret []byte) (string, error) {
	script, err := NewHTLCRedeemScript(secret)
	if err != nil {
		return "", fmt.Errorf("failed to create redeem script: %w", err)
	}
	txHash, err := s.client.Spend(s.htlcScript, script, s.redeemer)
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
	utxos, bal, err := s.client.GetUTXOs(s.scriptAddr, 0)
	if err != nil {
		return false, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if bal >= s.amount {
		return true, utxos[0].TxID, nil
	}
	return false, "", nil
}
