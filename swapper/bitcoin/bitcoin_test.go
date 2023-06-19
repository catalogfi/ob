package bitcoin_test

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/bitcoin"
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
		PRIV_KEY_1 := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" //mvb8yA23gtNPsBpd21Wq5J6YY4GEnfYQyX
		PRIV_KEY_2 := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" //myS2zesC4Va7ofV5MtnqZDct8iZdaBzULE

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

		client := bitcoin.NewClient("http://localhost:30000", &chaincfg.RegressionNetParams)

		secret, _ := randomHex(32)
		secret_hash := sha256.Sum256(secret)

		iSwapA, err := bitcoin.NewInitiatorSwap(privKey1, pkAddr2, secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())
		rSwapA, err := bitcoin.NewRedeemerSwap(privKey1, pkAddr2, secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())

		iSwapB, err := bitcoin.NewInitiatorSwap(privKey2, pkAddr1, secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())
		rSwapB, err := bitcoin.NewRedeemerSwap(privKey2, pkAddr1, secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())

		go func() {
			defer GinkgoRecover()
			Expect(swapper.ExecuteAtomicSwapFirst(iSwapA, rSwapA, secret)).To(BeNil())
		}()
		Expect(swapper.ExecuteAtomicSwapSecond(iSwapB, rSwapB)).To(BeNil())
	})
})
