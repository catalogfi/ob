package rest_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	// s *rest.Server
	c              rest.Client
	jwtToken       string
	CurrentOrderID uint
)

var _ = BeforeSuite(func() {
	// done := make(chan bool)
	StartServer()
	// <-done
	time.Sleep(3 * time.Second) // await server to start
	c = rest.NewClient("http://localhost:8080", os.Getenv("PRIVATE_KEY"))
	jwtToken = ""
	CurrentOrderID = 0
})

var _ = Describe("Rest", func() {
	It("check health of server", func() {
		Expect(c).NotTo(BeNil())
		Expect(jwtToken).To(Equal(""))
		Expect(CurrentOrderID).To(Equal(uint(0)))
		Expect(c.Health()).To(Equal("ok"))
	})

	It("check nonce", func() {
		nonce, err := c.GetNonce()
		Expect(err).NotTo(HaveOccurred())
		Expect(nonce).ToNot(BeNil())
	})

	It("check verify", func() {
		verified, err := c.Login()
		Expect(err).NotTo(HaveOccurred())
		Expect(verified).ToNot(BeNil())
		fmt.Println("verified: ", verified)
		c.SetJwt(verified)
	})

	It("should create Order", func() {
		OrderID, err := c.CreateOrder("mg54DDo5jfNkx5tF4d7Ag6G6VrJaSjr7ES", "0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da", "bitcoin:primary-ethereum:primary", "1", "10", "0xd87c01599e0f31a714ca73e5de993e274430101d4675d80da19d84b2bf19817d")
		CurrentOrderID = OrderID
		Expect(err).NotTo(HaveOccurred())
		Expect(OrderID).To(BeNumerically(">=", 0))
	})

	It("fill the order", func() {
		err := c.FillOrder(CurrentOrderID, "0xF403cE7776B22B74EcA871EcDaBeAA2103CD4A49", "mg54DDo8jfNkx5tF4d7Ag6G6VrJaSjr7ES")
		Expect(err).NotTo(HaveOccurred())
	})

	It("get Initiator Initiate Orders", func() {
		orders, err := c.GetInitiatorInitiateOrders()
		Expect(err).NotTo(HaveOccurred())
		Expect(orders).ToNot(BeNil())
		Expect(len(orders)).To(BeNumerically(">=", 1))
		CurrentOrder := orders[len(orders)-1]
		Expect(CurrentOrder.ID).To(Equal(CurrentOrderID))
		// Expect(CurrentOrder.Status).To(Equal("FILLED"))
		Expect(CurrentOrder.OrderPair).To(Equal("bitcoin:primary-ethereum:primary"))
		Expect(CurrentOrder.Taker).To(Equal("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da"))
		Expect(CurrentOrder.Maker).To(Equal("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da"))
		Expect(CurrentOrder.Status).To(Equal(model.Status(2)))
		// fmt.Println("CurrentOrder: ", CurrentOrder.Taker, CurrentOrder.Status, CurrentOrder.OrderPair)
	})

	It("get Followers init Orders", func() {
		orders, err := c.GetFollowerInitiateOrders()
		Expect(err).NotTo(HaveOccurred())
		Expect(orders).ToNot(BeNil())
		Expect(len(orders)).To(BeNumerically(">=", 0))

		// fmt.Println("CurrentOrder: ", CurrentOrder.Taker, CurrentOrder.Status, CurrentOrder.OrderPair)
	})

	It("get Initiator Redeem Orders", func() {
		order, err := c.GetInitiatorRedeemOrders()
		Expect(err).NotTo(HaveOccurred())
		Expect(order).ToNot(BeNil())
		Expect(len(order)).To(BeNumerically("==", 0)) // as atomic swap is not implemented yet
	})

	It("get Followers Redeem Orders", func() {
		order, err := c.GetFollowerRedeemOrders()
		Expect(err).NotTo(HaveOccurred())
		Expect(order).ToNot(BeNil())
		Expect(len(order)).To(BeNumerically("==", 0)) // as atomic swap is not implemented yet
	})

})

func StartServer() {
	go func() {
		store, err := store.New(sqlite.Open("gorm.db"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		config := model.Config{
			RPC: map[model.Chain]string{
				model.BitcoinRegtest:   "http://localhost:30000",
				model.EthereumLocalnet: "http://localhost:8545",
			},
		}
		s := rest.NewServer(store, config, "PANTHER")
		s.Run(":8080")
	}()
}
