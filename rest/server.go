package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spruceid/siwe-go"
	"github.com/susruth/wbtc-garden/model"
)

type Server struct {
	router *gin.Engine
	store  Store
	auth   Auth
	config model.Config
	secret string
}

type Store interface {
	// create order
	CreateOrder(creator, sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string, urls map[model.Chain]string) (uint, error)
	// fill order
	FillOrder(orderID uint, filler, sendAddress, recieveAddress string, urls map[model.Chain]string) error
	// get order by id
	GetOrder(orderID uint) (*model.Order, error)
	// cancel order by id
	CancelOrder(creator string, orderID uint) error
	// get all orders for the given user
	FilterOrders(maker, taker, orderPair, secretHash, sort string, status model.Status, minPrice, maxPrice float64, page, perPage int, verbose bool) ([]model.Order, error)
}

func NewServer(store Store, config model.Config, secret string) *Server {
	return &Server{
		router: gin.Default(),
		store:  store,
		secret: secret,
		auth:   NewAuth(),
		config: config,
	}
}

func (s *Server) Run(addr string) error {
	s.router.Use(cors.Default())

	authRoutes := s.router.Group("/")
	authRoutes.Use(s.authenticateJWT)

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	s.router.GET("/orders/:id", s.GetOrder())
	s.router.GET("/orders", s.GetOrders())
	s.router.GET("/nonce", s.Nonce())
	s.router.POST("/verify", s.Verify())
	{
		authRoutes.POST("/orders", s.PostOrders())
		authRoutes.PUT("/orders/:id", s.FillOrder())
		authRoutes.DELETE("/orders/:id", s.CancelOrder())
	}
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

type Auth interface {
	Verfiy(req model.VerifySiwe) (*jwt.Token, error)
}

func (s *Server) authenticateJWT(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		ctx.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method")
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userWallet, exists := claims["userWallet"]; exists {
			ctx.Set("userWallet", userWallet.(string))

		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (s *Server) PostOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		creator, exists := c.Get("userWallet")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		req := CreateOrder{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		oid, err := s.store.CreateOrder(creator.(string), req.SendAddress, req.RecieveAddress, req.OrderPair, req.SendAmount, req.RecieveAmount, req.SecretHash, s.config.RPC)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to create order: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMessage,
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
		filler, exists := c.Get("userWallet")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		req := FillOrder{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.store.FillOrder(uint(orderID), filler.(string), req.SendAddress, req.RecieveAddress, s.config.RPC); err != nil {
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
		maker, exists := c.Get("userWallet")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode id has to be a number: %v", err.Error())})
			return
		}
		if err := s.store.CancelOrder(maker.(string), uint(orderID)); err != nil {
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
		if err != nil && status < int(model.Unknown) || status > int(model.OrderFailedSoft) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("failed to decode status has to be a number between %d and %d", model.Unknown, model.OrderFailedSoft)})
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

func (s *Server) Nonce() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"nonce": siwe.GenerateNonce(),
		})
	}
}

func (s *Server) Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.VerifySiwe{}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := s.auth.Verfiy(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tokenString, err := token.SignedString([]byte(s.secret))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": tokenString})

	}
}
