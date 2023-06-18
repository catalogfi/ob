package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/susruth/wbtc-garden/model"
)

type Server struct {
	router   *gin.Engine
	swappers map[string]Swapper
}

type Swapper interface {
	GetAccount() (model.Account, error)
	GetAddresses(from, to, secretHash string, wbtcExpiry int64) (model.HTLCAddresses, error)
	ExecuteSwap(from, to, secretHash string, wbtcExpiry int64) error
	Transactions(address string) ([]model.Transaction, error)
}

func NewServer(swappers map[string]Swapper) *Server {
	return &Server{
		router:   gin.Default(),
		swappers: swappers,
	}
}

func (s *Server) Run(addr string) error {
	s.router.Use(cors.Default())
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	chains := []string{}
	for chain, swapper := range s.swappers {
		s.router.GET(fmt.Sprintf("/%s", chain), s.GetAccount(swapper))
		s.router.GET(fmt.Sprintf("/%s/addresses/:from/:to/:secretHash/:expiry", chain), s.GetAccount(swapper))
		s.router.POST(fmt.Sprintf("/%s/transactions", chain), s.PostTransactions(swapper))
		s.router.GET(fmt.Sprintf("/%s/transactions/:address", chain), s.GetTransactions(swapper))
		chains = append(chains, chain)
	}
	s.router.GET("/chains", func(c *gin.Context) {
		c.JSON(http.StatusOK, chains)
	})
	return s.router.Run(addr)
}

func (s *Server) GetAccount(swapper Swapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		account, err := swapper.GetAccount()
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
}

func (s *Server) PostTransactions(swapper Swapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := PostTransactionReq{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := swapper.ExecuteSwap(req.From, req.To, req.SecretHash, int64(req.WBTCExpiry)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to execute the swap",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}

type GetTransactionsResp []model.Transaction

func (s *Server) GetTransactions(swapper Swapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAddress := c.Param("address")
		transactions, err := swapper.Transactions(userAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to retrieve user transactions",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, transactions)
	}
}

func (s *Server) GetAddresses(swapper Swapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		from := c.Param("from")
		to := c.Param("to")
		secretHash := c.Param("secretHash")
		wbtcExpiry := c.Param("expiry")
		expiry, err := strconv.ParseInt(wbtcExpiry, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to parse expiry",
				"message": err.Error(),
			})
		}

		addrs, err := swapper.GetAddresses(from, to, secretHash, expiry)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to retrieve user transactions",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, addrs)
	}
}
