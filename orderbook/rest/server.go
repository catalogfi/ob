package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/susruth/wbtc-garden/orderbook/model"
)

type Server struct {
	router *gin.Engine
	store  Store
}

type Store interface {
	// create order
	CreateOrder(creator, sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) (uint, error)
	// fill order
	FillOrder(orderID uint, filler, sendAddress, recieveAddress string, initiateAtomicSwapTimelock, followerAtomicSwapTimelock uint64) error
	// get order by id
	GetOrder(orderID uint) (*model.Order, error)
	// cancel order by id
	CancelOrder(creator string, orderID uint) error
	// get all orders for the given user
	FilterOrders(maker, taker, orderPair, secretHash, sort string, status model.Status, minPrice, maxPrice float64, page, perPage int, verbose bool) ([]model.Order, error)
}

func NewServer(store Store) *Server {
	return &Server{
		router: gin.Default(),
		store:  store,
	}
}

func (s *Server) Run(addr string) error {
	s.router.Use(cors.Default())
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	s.router.GET("/orders/:id", s.GetOrder())
	s.router.GET("/orders", s.GetOrders())
	s.router.POST("/orders", s.PostOrders())        // TODO: add auth middleware
	s.router.PUT("/orders/:id", s.FillOrder())      // TODO: add auth middleware
	s.router.DELETE("/orders/:id", s.CancelOrder()) // TODO: add auth middleware
	return s.router.Run(addr)
}

type CreateOrder struct {
	SendAddress    string `json:"sendAddress"`
	RecieveAddress string `json:"recieveAddress"`
	OrderPair      string `json:"orderPair"`
	SendAmount     string `json:"sendAmount"`
	RecieveAmount  string `json:"recieveAmount"`
	SecretHash     string `json:"secretHash"`
}

func (s *Server) PostOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: extract from auth token
		creator := ""
		req := CreateOrder{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		oid, err := s.store.CreateOrder(creator, req.SendAddress, req.RecieveAddress, req.OrderPair, req.SendAmount, req.RecieveAmount, req.SecretHash)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to create order",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"orderId": oid,
		})
	}
}

type FillOrder struct {
	SendAddress    string `json:"sendAddress"`
	RecieveAddress string `json:"recieveAddress"`
}

func (s *Server) FillOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode id has to be a number: %v", err.Error())})
			return
		}

		// TODO: extract from auth token
		filler := ""

		// TODO: calculate initiator and follower timelocks
		initiatorTimeLock := uint64(0)
		followerTimeLock := uint64(0)

		req := FillOrder{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.store.FillOrder(uint(orderID), filler, req.SendAddress, req.RecieveAddress, initiatorTimeLock, followerTimeLock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusAccepted, gin.H{})
	}
}

func (s *Server) GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode id has to be a number: %v", err.Error())})
			return
		}
		order, err := s.store.GetOrder(uint(orderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, order)
	}
}

func (s *Server) CancelOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: extract from auth token
		maker := ""

		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode id has to be a number: %v", err.Error())})
			return
		}
		if err := s.store.CancelOrder(maker, uint(orderID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func (s *Server) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		maker := c.DefaultQuery("maker", "")
		taker := c.DefaultQuery("taker", "")
		orderPair := c.DefaultQuery("order_pair", "")
		secretHash := c.DefaultQuery("secret_hash", "")
		orderBy := c.DefaultQuery("sort", "")
		verbose, err := strconv.ParseBool(c.DefaultQuery("verbose", "false"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode verbose has to be a boolean: %v", err.Error())})
			return
		}

		status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
		if err != nil && status < int(model.Unknown) || status > int(model.OrderFailed) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode status has to be a number between %d and %d", model.Unknown, model.OrderFailed)})
			return
		}

		minPrice, err := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode minPrice has to be a number: %v", err.Error())})
			return
		}
		maxPrice, err := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode maxPrice has to be a number: %v", err.Error())})
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode page has to be a number: %v", err.Error())})
			return
		}
		perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode per_page has to be a number: %v", err.Error())})
			return
		}

		orders, err := s.store.FilterOrders(maker, taker, orderPair, secretHash, orderBy, model.Status(status), minPrice, maxPrice, page, perPage, verbose)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}
