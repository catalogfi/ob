package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/susruth/wbtc-garden/model"
)

type Server struct {
	router  *gin.Engine
	store   Store
	swapper Swapper
}

type Store interface {
	Transactions(address string) ([]model.Transaction, error)
}

type Swapper interface {
	GetAccount() (model.Account, error)
	ExecuteSwap(from, to, secretHash string, wbtcExpiry int64, amount uint64) error
}

func NewServer(store Store, swapper Swapper) *Server {
	return &Server{
		router:  gin.Default(),
		store:   store,
		swapper: swapper,
	}
}

func (s *Server) Run(addr string) error {
	s.router.GET("/", s.GetAccount())
	s.router.POST("/transactions", s.PostTransactions())
	s.router.GET("/transactions/:address", s.GetTransactions())

	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Include this line if you need to allow credentials (e.g., cookies)
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		c.Next()
	})
	return s.router.Run(addr)
}

func (s *Server) GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		account, err := s.swapper.GetAccount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, account)
	}
}

type PostTransactionReq struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	SecretHash string  `json:"secretHash"`
	WBTCExpiry float64 `json:"wbtcExpiry"`
	Amount     float64 `json:"amount"`
}

func (s *Server) PostTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := PostTransactionReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.swapper.ExecuteSwap(req.From, req.To, req.SecretHash, int64(req.WBTCExpiry), uint64(req.Amount*100000000)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to execute the swap",
				"message": err.Error(),
			})
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}

type GetTransactionsResp []model.Transaction

func (s *Server) GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddress := c.Param("address")
		transactions, err := s.store.Transactions(userAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to retrieve user transactions",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, transactions)
	}
}
