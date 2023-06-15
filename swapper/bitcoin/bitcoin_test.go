package bitcoin_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/susruth/wbtc-garden-server/swapper"
	"github.com/susruth/wbtc-garden-server/swapper/bitcoin"
)

var _ = Describe("Bitcoin", func() {
	It("should create a new swap", func() {

		PRIV_KEY_1 := os.Getenv("PRIV_KEY_1")
		PRIV_KEY_2 := os.Getenv("PRIV_KEY_2")

		privKeyBytes1, _ := hex.DecodeString(PRIV_KEY_1)
		privKey1, _ := btcec.PrivKeyFromBytes(privKeyBytes1)

		privKeyBytes2, _ := hex.DecodeString(PRIV_KEY_2)
		privKey2, _ := btcec.PrivKeyFromBytes(privKeyBytes2)

		pkAddr1, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey1.PubKey().SerializeCompressed()), &chaincfg.TestNet3Params)
		Expect(err).To(BeNil())
		fmt.Println(pkAddr1.EncodeAddress())

		pkAddr2, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(privKey2.PubKey().SerializeCompressed()), &chaincfg.TestNet3Params)
		Expect(err).To(BeNil())
		fmt.Println(pkAddr2.EncodeAddress())

		client := bitcoin.NewClient("https://blockstream.info/testnet/api", &chaincfg.TestNet3Params)

		secret := []byte("super_secret_____swap")
		secret_hash := sha256.Sum256(secret)

		iSwapA, err := bitcoin.NewInitiatorSwap(privKey1, privKey2.PubKey().SerializeCompressed(), secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())
		rSwapA, err := bitcoin.NewRedeemerSwap(privKey1, privKey2.PubKey().SerializeCompressed(), secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())

		iSwapB, err := bitcoin.NewInitiatorSwap(privKey2, privKey1.PubKey().SerializeCompressed(), secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())
		rSwapB, err := bitcoin.NewRedeemerSwap(privKey2, privKey1.PubKey().SerializeCompressed(), secret_hash[:], 1000, 10000, client)
		Expect(err).To(BeNil())

		go func() {
			defer GinkgoRecover()
			Expect(swapper.ExecuteAtomicSwapFirst(iSwapA, rSwapA, secret)).To(BeNil())
		}()
		Expect(swapper.ExecuteAtomicSwapSecond(iSwapB, rSwapB)).To(BeNil())
	})
})
