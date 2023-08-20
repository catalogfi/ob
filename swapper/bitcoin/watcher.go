package bitcoin

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/catalogfi/wbtc-garden/swapper"
)

// Watcher implements the `swapper.Watcher` interface. It watches a particular HTLC contract address and detect state
// change by querying the bitcoin network. People can use this to check what kind of the state the atomic swap is at.
// It assumes the HTLC contract address is disposable and will only be used once.
type watcher struct {
	client           Client
	scriptAddr       btcutil.Address
	amount           uint64
	minConfirmations uint64
	waitBlocks       int64
}

func NewWatcher(scriptAddr btcutil.Address, waitBlocks int64, minConfirmations, amount uint64, client Client) (swapper.Watcher, error) {
	return &watcher{
		scriptAddr:       scriptAddr,
		amount:           amount,
		waitBlocks:       waitBlocks,
		minConfirmations: minConfirmations,
		client:           client,
	}, nil
}

func (w *watcher) Expired() (bool, error) {
	// Check if the swap has been initiated
	initiated, txHash, _, err := w.IsInitiated()
	if err != nil {
		return false, err
	}
	if !initiated {
		return false, nil
	}
	height, _, err := w.Status(txHash)
	if err != nil {
		return false, err
	}
	// Check the number of blocks has been passed since the initiation
	latest, err := w.client.GetTipBlockHeight()
	if err != nil {
		return false, err
	}
	diff := latest - height + 1
	return diff >= uint64(w.waitBlocks), nil
}

func (w *watcher) IsDetected() (bool, string, string, error) {
	// Fetch all utxos
	utxos, bal, err := w.client.GetUTXOs(w.scriptAddr, 0)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if len(utxos) == 0 {
		return false, "", "", nil
	}
	txHashes := make([]string, len(utxos))
	for i, utxo := range utxos {
		txHashes[i] = utxo.TxID
	}
	return bal >= w.amount, strings.Join(txHashes, ","), strconv.FormatUint(bal, 10), nil
}

func (w *watcher) IsInitiated() (bool, string, uint64, error) {
	// Fetch all utxos
	utxos, bal, err := w.client.GetUTXOs(w.scriptAddr, 0)
	if err != nil {
		return false, "", 0, fmt.Errorf("failed to get UTXOs: %w", err)
	}

	// Check all utxos are confirmed and greater than or equal to the required amount
	if bal >= w.amount && len(utxos) > 0 {
		latest, err := w.client.GetTipBlockHeight()
		if err != nil {
			return false, "", 0, fmt.Errorf("failed to get tip block: %w", err)
		}
		txHashes := make([]string, len(utxos))
		lastConfirmedTxBlock := uint64(0)
		for i, utxo := range utxos {
			confirmation := latest - utxo.Status.BlockHeight + 1
			if utxo.Status == nil || !utxo.Status.Confirmed || confirmation < w.minConfirmations {
				if confirmation < w.minConfirmations && utxo.Status.BlockHeight > 0 {
					return false, "", confirmation, nil
				}
				return false, "", 0, nil
			}
			txHashes[i] = utxo.TxID
			if utxo.Status.BlockHeight > lastConfirmedTxBlock {
				lastConfirmedTxBlock = utxo.Status.BlockHeight
			}
		}
		return true, strings.Join(txHashes, ","), w.minConfirmations, nil
	}
	return false, "", 0, nil
}

func (w *watcher) Status(initateTxHash string) (uint64, uint64, error) {
	txHashes := strings.Split(initateTxHash, ",")
	if len(txHashes) == 0 {
		return 0, 0, fmt.Errorf("empty initiate txhash list")
	}
	blockHeight, conf, err := w.client.GetConfirmations(txHashes[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get confirmations: %w", err)
	}
	if len(txHashes) > 1 {
		for _, txHash := range txHashes[1:] {
			nextBlockHeight, nextConf, err := w.client.GetConfirmations(txHash)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to get confirmations: %w", err)
			}
			if nextBlockHeight < blockHeight {
				blockHeight = nextBlockHeight
			}
			if nextConf < conf {
				conf = nextConf
			}
		}
	}
	return blockHeight, conf, err
}

// IsRedeemed checks if the secret has been revealed on-chain.
func (w *watcher) IsRedeemed() (bool, []byte, string, error) {
	// witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	// if err != nil {
	// 	return false, nil, "", fmt.Errorf("failed to get UTXOs: %w", err)
	// }
	// if len(witness) == 5 {
	// 	// Check if the redeem tx is confirmed
	// 	// latest, err := w.client.GetTipBlockHeight()
	// 	// if err != nil {
	// 	// 	return false, nil, "", fmt.Errorf("failed to get tip block: %w", err)
	// 	// }
	// 	// if !tx.Status.Confirmed || latest-tx.Status.BlockHeight+1 < w.minConfirmations {
	// 	// 	return false, nil, "", nil
	// 	// }

	// 	// inputs are [
	// 	// 0 : sig,
	// 	// 1 : spender.PubKey().SerializeCompressed(),
	// 	// 2 : secret,
	// 	// 3 :[]byte{0x1},
	// 	// 4 : script
	// 	// ]
	// 	secretBytes, err := hex.DecodeString(witness[2])
	// 	if err != nil {
	// 		return false, nil, "", fmt.Errorf("failed to decode secret: %w", err)
	// 	}

	// 	return true, secretBytes, tx.TxID, nil
	// }
	// return false, nil, "", nil
	witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	if err != nil {
		return false, nil, "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	if len(witness) == 5 {
		// fmt.Println("Redeemed:", witness)
		// inputs are [ 0 : sig, 1 : spender.PubKey().SerializeCompressed(),2 : secret, 3 :[]byte{0x1}, script]
		secretString := witness[2]
		secretBytes := make([]byte, hex.DecodedLen(len(secretString)))
		_, err := hex.Decode(secretBytes, []byte(secretString))
		if err != nil {
			return false, nil, "", fmt.Errorf("failed to decode secret: %w", err)
		}
		return true, secretBytes, tx.TxID, nil
	}
	return false, nil, "", nil
}

// IsRefunded checks if there's any refund tx with the script address. It will return `true` even if the script
// address is not fully refunded.
func (w *watcher) IsRefunded() (bool, string, error) {
	// witness, tx, err := w.client.GetSpendingWitness(w.scriptAddr)
	// if err != nil {
	// 	return false, "", fmt.Errorf("failed to get UTXOs: %w", err)
	// }

	// // Check if the refund tx is confirmed
	// latest, err := w.client.GetTipBlockHeight()
	// if err != nil {
	// 	return false, "", fmt.Errorf("failed to get tip block: %w", err)
	// }
	// if !tx.Status.Confirmed || latest-tx.Status.BlockHeight+1 < w.minConfirmations {
	// 	return false, "", nil
	// }

	// // inputs are [
	// // 0 : sig,
	// // 1 : spender.PubKey().SerializeCompressed(),
	// // 2 :[]byte{},
	// // script
	// // ]
	// return len(witness) == 4, tx.TxID, nil
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
		return true, tx.TxID, nil

	}
	return false, "", nil
}
