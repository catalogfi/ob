package bitcoin

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
)

type instantClient struct {
	url           string
	indexerClient Client
}

func InstantWalletWrapper(url string, client Client) Client {
	return &instantClient{url: url, indexerClient: client}
}

func (client *instantClient) Net() *chaincfg.Params {
	return client.indexerClient.Net()
}

func (client *instantClient) GetTipBlockHeight() (uint64, error) {
	return client.indexerClient.GetTipBlockHeight()
}

func (client *instantClient) GetBlockHeight(txhash string) (uint64, error) {
	return client.indexerClient.GetBlockHeight(txhash)
}

func (client *instantClient) GetUTXOs(address btcutil.Address, amount uint64) (UTXOs, uint64, error) {
	return client.indexerClient.GetUTXOs(address, amount)
}

func (client *instantClient) GetSpendingWitness(address btcutil.Address) ([]string, string, error) {
	return client.indexerClient.GetSpendingWitness(address)
}

func (client *instantClient) Send(to btcutil.Address, amount uint64, from *btcec.PrivateKey) (string, error) {
	panic("not implemented")
}

func (client *instantClient) Spend(script []byte, scriptSig wire.TxWitness, spender *btcec.PrivateKey, waitBlocks uint) (string, error) {
	panic("not implemented")
}
