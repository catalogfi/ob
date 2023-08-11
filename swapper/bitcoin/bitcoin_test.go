package bitcoin_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/catalogfi/wbtc-garden/swapper"
	"github.com/catalogfi/wbtc-garden/swapper/bitcoin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func randomHex(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

var _ = Describe("Bitcoin", func() {
	It("should create a new swap", func() {
		// PRIV_KEY_1 := os.Getenv("PRIV_KEY_1")
		// PRIV_KEY_2 := os.Getenv("PRIV_KEY_2")
		// Skip("")
		PRIV_KEY_1 := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // mvb8yA23gtNPsBpd21Wq5J6YY4GEnfYQyX
		PRIV_KEY_2 := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" // myS2zesC4Va7ofV5MtnqZDct8iZdaBzULE

		privKeyBytes1, _ := hex.DecodeString(PRIV_KEY_1)
		privKey1, _ := btcec.PrivKeyFromBytes(privKeyBytes1)

		privKeyBytes2, _ := hex.DecodeString(PRIV_KEY_2)
		privKey2, _ := btcec.PrivKeyFromBytes(privKeyBytes2)

		pkAddr1, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey1.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr1:", pkAddr1.EncodeAddress())

		pkAddr2, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey2.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr2:", pkAddr2.EncodeAddress())

		client := bitcoin.NewClient("https://mempool.space/testnet/api", &chaincfg.RegressionNetParams)

		secret, _ := randomHex(32)
		secret_hash := sha256.Sum256(secret)

		iSwapA, err := bitcoin.NewInitiatorSwap(privKey1, pkAddr2, secret_hash[:], 1000, 0, 10000, client)
		Expect(err).To(BeNil())
		rSwapA, err := bitcoin.NewRedeemerSwap(privKey1, pkAddr2, secret_hash[:], 1000, 0, 10000, client)
		Expect(err).To(BeNil())

		iSwapB, err := bitcoin.NewInitiatorSwap(privKey2, pkAddr1, secret_hash[:], 1000, 0, 10000, client)
		Expect(err).To(BeNil())
		rSwapB, err := bitcoin.NewRedeemerSwap(privKey2, pkAddr1, secret_hash[:], 1000, 0, 10000, client)
		Expect(err).To(BeNil())

		go func() {
			defer GinkgoRecover()
			Expect(swapper.ExecuteAtomicSwapFirst(iSwapA, rSwapA, secret)).To(BeNil())
		}()
		Expect(swapper.ExecuteAtomicSwapSecond(iSwapB, rSwapB)).To(BeNil())
	})
	It("should create a new swap and refund", func() {
		// Skip("")
		PRIV_KEY_1 := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // mvb8yA23gtNPsBpd21Wq5J6YY4GEnfYQyX
		PRIV_KEY_2 := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" // myS2zesC4Va7ofV5MtnqZDct8iZdaBzULE

		privKeyBytes1, _ := hex.DecodeString(PRIV_KEY_1)
		privKey1, _ := btcec.PrivKeyFromBytes(privKeyBytes1)

		privKeyBytes2, _ := hex.DecodeString(PRIV_KEY_2)
		privKey2, _ := btcec.PrivKeyFromBytes(privKeyBytes2)

		pkAddr1, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey1.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr1:", pkAddr1.EncodeAddress())

		pkAddr2, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey2.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr2:", pkAddr2.EncodeAddress())

		client := bitcoin.NewClient("https://mempool.space/testnet/api", &chaincfg.RegressionNetParams)

		secret, _ := randomHex(32)
		secret_hash := sha256.Sum256(secret)

		iSwapA, err := bitcoin.NewInitiatorSwap(privKey1, pkAddr2, secret_hash[:], 5, 0, 10000, client)
		Expect(err).To(BeNil())

		_, err = iSwapA.Initiate()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 10)
		_, err = iSwapA.Refund()
		if err != nil {
			panic(err)
		}

	})
	It("should send via segwit", func() {
		PRIV_KEY_1 := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" // tb1q5428vq2uzwhm3taey9sr9x5vm6tk78ew9gs838
		PRIV_KEY_2 := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" // tb1qcjzphr67dug28rw9ueewrqllmxlqe5f0v7g34m

		privKeyBytes1, _ := hex.DecodeString(PRIV_KEY_1)
		privKey1, _ := btcec.PrivKeyFromBytes(privKeyBytes1)

		privKeyBytes2, _ := hex.DecodeString(PRIV_KEY_2)
		privKey2, _ := btcec.PrivKeyFromBytes(privKeyBytes2)

		pkAddr1, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(privKey1.PubKey().SerializeCompressed()), &chaincfg.TestNet3Params)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr1:", pkAddr1.EncodeAddress())
		pkAddr2, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(privKey2.PubKey().SerializeCompressed()), &chaincfg.TestNet3Params)
		Expect(err).To(BeNil())
		fmt.Println("pkAddr2:", pkAddr2.EncodeAddress())

		client := bitcoin.NewClient("https://mempool.space/testnet/api", &chaincfg.TestNet3Params)
		txhash, err := client.Send(pkAddr2, 10000, privKey1)
		Expect(txhash).NotTo(BeNil())
		Expect(err).To(BeNil())
		fmt.Printf("https://mempool.space/testnet/tx/%s", txhash)
	})
})
