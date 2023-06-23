package rest_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/susruth/wbtc-garden/orderbook/model"
	"github.com/susruth/wbtc-garden/orderbook/rest"
	"github.com/susruth/wbtc-garden/orderbook/store"
	"github.com/susruth/wbtc-garden/orderbook/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest Suite")
}

var _ = Describe("Rest", func() {
	var (
		// s *rest.Server
		c              rest.Client
		jwtToken       string
		CurrentOrderID uint
	)

	BeforeSuite(func() {
		// done := make(chan bool)
		StartServer()
		// <-done
		time.Sleep(3 * time.Second) // await server to start
		c = rest.NewClient("http://localhost:8080")
		jwtToken = ""
		CurrentOrderID = 0

	})

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
		verified, err := c.Verify("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da")
		Expect(err).NotTo(HaveOccurred())
		Expect(verified).ToNot(BeNil())
		fmt.Println("verified: ", verified)
		c.SetJwt(verified)
	})

	It("check create order", func() {
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
		orders, err := c.GetInitiatorInitiateOrders("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da")
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
		orders, err := c.GetFollowerInitiateOrders("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da")
		Expect(err).NotTo(HaveOccurred())
		Expect(orders).ToNot(BeNil())
		Expect(len(orders)).To(BeNumerically(">=", 0))
		
		// fmt.Println("CurrentOrder: ", CurrentOrder.Taker, CurrentOrder.Status, CurrentOrder.OrderPair)
	})

	It("get Initiator Redeem Orders", func() {
		order, err := c.GetInitiatorRedeemOrders("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da")
		Expect(err).NotTo(HaveOccurred())
		Expect(order).ToNot(BeNil())
		Expect(len(order)).To(BeNumerically("==", 0)) // as atomic swap is not implemented yet
	})

	It("get Followers Redeem Orders", func() {
		order, err := c.GetFollowerRedeemOrders("0x17100301bB2FF58aE6B5ca5B8f9Ec6F872E0F2da")
		Expect(err).NotTo(HaveOccurred())
		Expect(order).ToNot(BeNil())
		Expect(len(order)).To(BeNumerically("==", 0)) // as atomic swap is not implemented yet
	})

})

func StartServer() {
	go func() {
		store, err := store.NewStore(sqlite.Open("gorm.db"), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		auth := user.NewAuth()
		s := rest.NewServer(store, auth, "PANTHER")
		s.Run(":8080")
	}()
}
