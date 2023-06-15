package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const FEE = uint64(6000)
const BTC_VERSION = 2

type UTXO struct {
	Amount uint64 `json:"value"`
	TxID   string `json:"txid"`
	Vout   uint32 `json:"vout"`
}

type UTXOs []UTXO

type Client interface {
	GetSpendingScriptSig(address btcutil.Address) (string, string, error)
	GetBlockHeight(txhash string) (uint64, error)
	GetTipBlockHeight() (uint64, error)
	GetUTXOs(address btcutil.Address, amount uint64) (UTXOs, uint64, error)
	Send(to btcutil.Address, amount uint64, from *btcec.PrivateKey) (string, error)
	Spend(script []byte, scriptSig []byte, spender *btcec.PrivateKey) (string, error)
	Net() *chaincfg.Params
}

type client struct {
	url string
	net *chaincfg.Params
}

func NewClient(url string, net *chaincfg.Params) Client {
	return &client{url: url, net: net}
}

func (client *client) Net() *chaincfg.Params {
	return client.net
}

func (client *client) GetTipBlockHeight() (uint64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/blocks/tip/height", client.url))
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction: %w", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	return strconv.ParseUint(string(data), 10, 64)
}

func (client *client) GetBlockHeight(txhash string) (uint64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/tx/%s", client.url, txhash))
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction: %w", err)
	}

	var tx Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return 0, fmt.Errorf("failed to decode transaction: %w", err)
	}

	if !tx.Status.Confirmed {
		return 0, nil
	}
	return tx.Status.BlockHeight, nil
}

func (client *client) GetSpendingScriptSig(address btcutil.Address) (string, string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/address/%s/txs", client.url, address.EncodeAddress()))
	if err != nil {
		return "", "", fmt.Errorf("failed to get transactions: %w", err)
	}
	var txs []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&txs); err != nil {
		return "", "", fmt.Errorf("failed to decode transactions: %w", err)
	}
	for _, tx := range txs {
		for _, vin := range tx.VINs {
			if vin.Prevout.ScriptPubKeyType == "p2sh" {
				return vin.ScriptSigAsm, tx.TxID, nil
			}
		}
	}
	return "", "", nil
}

func (client *client) GetUTXOs(address btcutil.Address, amount uint64) (UTXOs, uint64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/address/%s/utxo", client.url, address.EncodeAddress()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get UTXOs: %w", err)
	}
	utxos := UTXOs{}
	if err := json.NewDecoder(resp.Body).Decode(&utxos); err != nil {
		return nil, 0, fmt.Errorf("failed to decode UTXOs: %w", err)
	}

	var balance uint64
	for _, utxo := range utxos {
		balance += utxo.Amount
	}

	if amount == 0 {
		return utxos, balance, nil
	}
	if balance < amount {
		return nil, 0, fmt.Errorf("insufficient balance in %s", address.EncodeAddress())
	}

	var selected UTXOs
	var total uint64
	for _, utxo := range utxos {
		if total >= amount {
			break
		}
		selected = append(selected, utxo)
		total += utxo.Amount
	}
	return selected, total, nil
}

func (client *client) Send(to btcutil.Address, amount uint64, from *btcec.PrivateKey) (string, error) {
	tx := wire.NewMsgTx(BTC_VERSION)

	fromAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(from.PubKey().SerializeCompressed()), client.Net())
	if err != nil {
		return "", fmt.Errorf("failed to create address from private key: %w", err)
	}

	utxos, selectedAmount, err := client.GetUTXOs(fromAddr, amount+FEE)
	if err != nil {
		return "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return "", fmt.Errorf("failed to parse txid in the utxo: %w", err)
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txid, utxo.Vout), nil, nil))
	}

	toScript, err := txscript.PayToAddrScript(to)
	if err != nil {
		return "", fmt.Errorf("failed to create script for address: %w", err)
	}

	fromScript, err := txscript.PayToAddrScript(fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to create script for address: %w", err)
	}

	tx.AddTxOut(wire.NewTxOut(int64(amount), toScript))
	tx.AddTxOut(wire.NewTxOut(int64(selectedAmount-amount-FEE), fromScript))

	for i := range tx.TxIn {
		sigScript, err := txscript.SignatureScript(tx, i, fromScript, txscript.SigHashAll, from, true)
		if err != nil {
			return "", fmt.Errorf("failed to sign transaction: %w", err)
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	return client.SubmitTx(tx)
}

func (client *client) Spend(script, redeemScript []byte, spender *btcec.PrivateKey) (string, error) {
	tx := wire.NewMsgTx(BTC_VERSION)

	scriptAddr, err := btcutil.NewAddressScriptHash(script, client.Net())
	if err != nil {
		return "", fmt.Errorf("failed to create script address: %w", err)
	}
	utxos, balance, err := client.GetUTXOs(scriptAddr, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get UTXOs: %w", err)
	}
	for _, utxo := range utxos {
		txid, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return "", fmt.Errorf("failed to parse txid in the utxo: %w", err)
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(txid, utxo.Vout), nil, nil))
	}

	spenderAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(spender.PubKey().SerializeCompressed()), client.Net())
	if err != nil {
		return "", fmt.Errorf("failed to create address from private key: %w", err)
	}
	spenderToScript, err := txscript.PayToAddrScript(spenderAddr)
	if err != nil {
		return "", fmt.Errorf("failed to create script for address: %w", err)
	}
	tx.AddTxOut(wire.NewTxOut(int64(balance-FEE), spenderToScript))

	for i := range tx.TxIn {
		sig, err := txscript.RawTxInSignature(tx, i, script, txscript.SigHashAll, spender)
		if err != nil {
			return "", err
		}
		sigScript, err := txscript.NewScriptBuilder().AddData(sig).AddOps(redeemScript).AddData(script).Script()
		if err != nil {
			return "", err
		}

		tx.TxIn[i].SignatureScript = sigScript
	}
	return client.SubmitTx(tx)
}

func (client *client) SubmitTx(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/tx", client.url), "application/text", bytes.NewBufferString(hex.EncodeToString(buf.Bytes())))
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to send transaction: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read transaction id: %w", err)
	}
	return string(data), nil
}

type Transaction struct {
	TxID   string `json:"txid"`
	VINs   []VIN  `json:"vin"`
	Status Status `json:"status"`
}

type VIN struct {
	TxID         string  `json:"txid"`
	Vout         int     `json:"vout"`
	Prevout      Prevout `json:"prevout"`
	ScriptSigAsm string  `json:"scriptsig_asm"`
}

type Prevout struct {
	ScriptPubKeyType string `json:"scriptpubkey_type"`
}

type Status struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint64 `json:"block_height"`
}
