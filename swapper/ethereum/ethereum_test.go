package ethereum_test

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/susruth/wbtc-garden/swapper"
	"github.com/susruth/wbtc-garden/swapper/ethereum"
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

		PRIV_KEY_1 := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" //0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
		PRIV_KEY_2 := "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" //0x70997970C51812dc3A010C7d01b50e0d17dc79C8
		DEPLOYER := common.HexToAddress("0x2D914D5a5F9d66c16475B2F5cA9a820308F8B35a")
		TOKEN := common.HexToAddress("0x401dDf5FD514c7C3AD5bFF8A70221ff7d091163F")

		privKey1, err := crypto.HexToECDSA(PRIV_KEY_1)
		Expect(err).To(BeNil())
		privKey2, err := crypto.HexToECDSA(PRIV_KEY_2)
		Expect(err).To(BeNil())
		pkAddr1 := crypto.PubkeyToAddress(privKey1.PublicKey)
		pkAddr2 := crypto.PubkeyToAddress(privKey2.PublicKey)

		fmt.Println(pkAddr1.Hex())
		fmt.Println(pkAddr2.Hex())

		client, err := ethereum.NewClient("http://localhost:8545")
		Expect(err).To(BeNil())

		secret, _ := randomHex(32)
		secret_hash := sha256.Sum256(secret)

		aExpiry, err := ethereum.GetExpiry(client, true)
		Expect(err).To(BeNil())

		bExpiry, err := ethereum.GetExpiry(client, false)
		Expect(err).To(BeNil())

		iSwapA, err := ethereum.NewInitiatorSwap(privKey1, pkAddr2, DEPLOYER, TOKEN, secret_hash[:], aExpiry, big.NewInt(100000), client)
		Expect(err).To(BeNil())
		rSwapA, err := ethereum.NewRedeemerSwap(privKey1, pkAddr2, DEPLOYER, TOKEN, secret_hash[:], bExpiry, big.NewInt(100000), client)
		Expect(err).To(BeNil())

		iSwapB, err := ethereum.NewInitiatorSwap(privKey2, pkAddr1, DEPLOYER, TOKEN, secret_hash[:], bExpiry, big.NewInt(100000), client)
		Expect(err).To(BeNil())
		rSwapB, err := ethereum.NewRedeemerSwap(privKey2, pkAddr1, DEPLOYER, TOKEN, secret_hash[:], aExpiry, big.NewInt(100000), client)
		Expect(err).To(BeNil())

		go func() {
			defer GinkgoRecover()
			Expect(swapper.ExecuteAtomicSwapFirst(iSwapA, rSwapA, secret)).To(BeNil())
		}()
		Expect(swapper.ExecuteAtomicSwapSecond(iSwapB, rSwapB)).To(BeNil())
	})
})
