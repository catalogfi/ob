package rest

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/catalogfi/orderbook/feehub"
	"github.com/catalogfi/orderbook/model"
	"github.com/catalogfi/orderbook/price"
	"github.com/catalogfi/orderbook/screener"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spruceid/siwe-go"
	"go.uber.org/zap"
)

const (
	GetOrderMessageRegex  = `^(?P<action>subscribe):(?P<orderID>\d+)$`
	GetOrdersMessageRegex = `^(?P<action>subscribe):(?P<address>0x[a-fA-F0-9]{40})$`
)

var upgrader = websocket.Upgrader{
	// check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		// return r.Header.Get("Origin") == "http://wbtcgarden"
		// TODO: add better origin checks
		return true
	},
}

type Server struct {
	router       *gin.Engine
	store        Store
	auth         Auth
	config       model.Config
	logger       *zap.Logger
	secret       string
	socketPool   SocketPool
	screener     screener.Screener
	feeHubClient feehub.FeehubClient
	priceClient  price.PriceFetcher
}

type Store interface {
	// get value locked in the given chain for the given user
	ValueLockedByChain(chain model.Chain, config model.Network) (*big.Int, error)
	// create order
	CreateOrder(creator, sendAddress, receiveAddress, orderPair, secretHash, userWalletBTCAddress string, sendAmount, receiveAmount, feeInBtc, feeInSeed *big.Int, IsDiscounted bool, config model.Config) (uint, error)
	// fill order
	FillOrder(orderID uint, filler, sendAddress, receiveAddress string, config model.Network) error
	// get order by id
	GetOrder(orderID uint) (*model.Order, error)
	// get order by atomic swap id
	GetOrderBySwapID(swapID uint) (*model.Order, error)
	// get orders by address
	GetOrdersByAddress(address string) ([]model.Order, error)

	GetPendingOrdersForAddress(address string) ([]model.Order, error)
	// cancel order by id
	CancelOrder(creator string, orderID uint) error
	// get all orders for the given user
	FilterOrders(maker, taker, orderPair, secretHash string, status model.Status, minPrice, maxPrice float64, minAmount, maxAmount float64, page, perPage int, verbose bool) ([]model.Order, error)

	GetSecrets(lastUpdated time.Time) ([]model.SecretRevealed, error)
}

func NewServer(store Store, config model.Config, logger *zap.Logger, secret string, socketPool SocketPool, screener screener.Screener, fhClient *feehub.FeehubClient, pc price.PriceFetcher) *Server {
	childLogger := logger.With(zap.String("service", "rest"))
	return &Server{
		router:       gin.Default(),
		store:        store,
		secret:       secret,
		logger:       childLogger,
		auth:         NewAuth(config.Network),
		config:       config,
		socketPool:   socketPool,
		screener:     screener,
		feeHubClient: *fhClient,
		priceClient:  pc,
	}
}

func (s *Server) Run(ctx context.Context, addr string) error {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authRoutes := s.router.Group("/")
	authRoutes.Use(s.authenticate)

	// websocket
	s.router.GET("/", s.socket())

	s.router.GET("/health", s.health())
	s.router.GET("/orders/:id", s.getOrder())
	s.router.GET("/orders", s.getOrders())
	s.router.GET("/nonce", s.nonce())
	s.router.GET("/assets", s.supportedAssets())
	s.router.GET("/chains/:chain/value", s.getValueByChain())
	s.router.GET("/secrets", s.secrets())
	s.router.POST("/verify", s.verify())
	{
		authRoutes.POST("/orders", s.postOrders())
		authRoutes.PUT("/orders/:id", s.fillOrder())
		authRoutes.DELETE("/orders/:id", s.cancelOrder())
	}

	server := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		fmt.Println("stopped")
	}()
	<-ctx.Done()
	return server.Shutdown(ctx)
}

type CreateOrder struct {
	SendAddress          string                    `json:"sendAddress" binding:"required"`
	ReceiveAddress       string                    `json:"receiveAddress" binding:"required"`
	OrderPair            string                    `json:"orderPair" binding:"required"`
	SendAmount           string                    `json:"sendAmount" binding:"required"`
	ReceiveAmount        string                    `json:"receiveAmount" binding:"required"`
	SecretHash           string                    `json:"secretHash" binding:"required"`
	UserWalletBTCAddress string                    `json:"userWalletBTCAddress" binding:"required"`
	FeePayment           feehub.ConditionalPayment `json:"feePayment"`
	Filler               string                    `json:"filler"`
	IsDiscounted         bool                      `json:"isDiscounted"`
}

type Auth interface {
	Verify(req model.VerifySiwe) (*jwt.Token, error)
}

func (s *Server) authenticate(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("missing authorization token")))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		ctx.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("invalid signing method")))
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		s.logger.Debug("authorization failure", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userWallet, exists := claims["userWallet"]; exists {
			ctx.Set("userWallet", strings.ToLower(userWallet.(string)))
			ctx.Set("token", tokenString)
		} else {
			s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("invalid token claims")))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (s *Server) health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "online",
		})
	}
}

func (s *Server) getValueByChain() gin.HandlerFunc {
	return func(c *gin.Context) {
		chain, err := model.ParseChain(c.Param("chain"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "chain not supported"})
			return
		}

		valueLocked, err := s.store.ValueLockedByChain(chain, s.config.Network)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"value": valueLocked,
		})
	}
}

func (s *Server) supportedAssets() gin.HandlerFunc {
	assets := map[model.Chain][]model.Asset{}
	for chain, netConf := range s.config.Network {
		assets[chain] = []model.Asset{}
		for asset := range netConf.Assets {
			assets[chain] = append(assets[chain], asset)
		}
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, assets)
	}
}

func (s *Server) postOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		creator, exists := c.Get("userWallet")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		token, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
			return
		}
		req := CreateOrder{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the addresses is blacklisted
		senderChain, receiverChain, _, _, err := model.ParseOrderPair(req.OrderPair)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		addrs := map[string]model.Chain{
			creator.(string):   model.Ethereum,
			req.ReceiveAddress: receiverChain,
			req.SendAddress:    senderChain,
		}

		blacklisted, err := s.screener.IsBlacklisted(addrs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if blacklisted {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address is blacklisted from database",
			})
			return
		}

		sendAmount, ok := new(big.Int).SetString(req.SendAmount, 10)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid send amount: %s", sendAmount),
			})
			return
		}

		receiveAmount, ok := new(big.Int).SetString(req.ReceiveAmount, 10)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid receive amount: %s", sendAmount),
			})
			return
		}

		feeInBtc := big.NewInt(0)
		if req.IsDiscounted {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			btcSeedPrice, err := s.priceClient.GetPrice(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("failed to get btc seed price : %v", err.Error()),
				})
				return
			}

			fee := new(big.Float).Quo(new(big.Float).SetInt(req.FeePayment.HTLC.RecvAmount.Int), new(big.Float).SetFloat64(btcSeedPrice.Price))
			feeInBtc, _ = fee.Int(nil)
			feeInBtc = new(big.Int).Div(feeInBtc, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

			maxFee := new(big.Int).Div(sendAmount, big.NewInt(100))
			maxFee = maxFee.Add(maxFee, big.NewInt(30000))

			if feeInBtc.Cmp(maxFee) > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("fee is greater than : %s", maxFee.String()),
				})
				return
			}
			ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = s.feeHubClient.ProcessPayFillerDelegates(ctx, req.FeePayment, req.Filler, token.(string))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("failed to process pay filler : %v", err.Error()),
				})
				return
			}
		}

		oid, err := s.store.CreateOrder(
			strings.ToLower(creator.(string)),
			req.SendAddress,
			req.ReceiveAddress,
			req.OrderPair,
			req.SecretHash,
			req.UserWalletBTCAddress,
			sendAmount,
			receiveAmount,
			feeInBtc,
			req.FeePayment.HTLC.RecvAmount.Int,
			req.IsDiscounted,
			s.config)
		if err != nil {
			errorMessage := fmt.Sprintf("failed to create order: %v", err.Error())
			// fmt.Println(errorMessage, "error")
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
	SendAddress    string `json:"sendAddress" binding:"required"`
	ReceiveAddress string `json:"receiveAddress" binding:"required"`
}

func (s *Server) fillOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode id has to be a number: %v", err.Error())})
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

		// Get details of the order to fill
		order, err := s.store.GetOrder(uint(orderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("error getting the order %v", err.Error()),
			})
			return
		}

		// Check if the addresses is blacklisted
		senderChain, receiverChain, _, _, err := model.ParseOrderPair(order.OrderPair)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		addrs := map[string]model.Chain{
			filler.(string):    model.Ethereum,
			req.ReceiveAddress: senderChain,
			req.SendAddress:    receiverChain,
		}

		blacklisted, err := s.screener.IsBlacklisted(addrs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if blacklisted {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Address is blacklisted from database",
			})
			return
		}

		if err := s.store.FillOrder(uint(orderID), strings.ToLower(filler.(string)), req.SendAddress, req.ReceiveAddress, s.config.Network); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to fill the Order %v", err.Error()),
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{})
	}
}

func (s *Server) getOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode id has to be a number: %v", err.Error())})
			return
		}
		order, err := s.store.GetOrder(uint(orderID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to get order %s", err.Error()),
			})
		}
		c.JSON(http.StatusOK, order)
	}
}

func (s *Server) cancelOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: extract from auth token
		maker, exists := c.Get("userWallet")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode id has to be a number: %v", err.Error())})
			return
		}
		if err := s.store.CancelOrder(strings.ToLower(maker.(string)), uint(orderID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to get account details",
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

func (s *Server) getOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		maker := c.DefaultQuery("maker", "")
		taker := c.DefaultQuery("taker", "")
		orderPair := c.DefaultQuery("order_pair", "")
		secretHash := c.DefaultQuery("secret_hash", "")

		maker = strings.ToLower(maker)
		taker = strings.ToLower(taker)

		verbose, err := strconv.ParseBool(c.DefaultQuery("verbose", "false"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode verbose has to be a boolean: %v", err.Error())})
			return
		}

		status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
		if err != nil && status < int(model.Unknown) || status > int(model.FailedSoft) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode status has to be a number between %d and %d", model.Unknown, model.FailedSoft)})
			return
		}

		minPrice, err := strconv.ParseFloat(c.DefaultQuery("min_price", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode minPrice has to be a number: %v", err.Error())})
			return
		}
		maxPrice, err := strconv.ParseFloat(c.DefaultQuery("max_price", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode maxPrice has to be a number: %v", err.Error())})
			return
		}
		page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode page has to be a number: %v", err.Error())})
			return
		}
		perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode per_page has to be a number: %v", err.Error())})
			return
		}

		minAmount, err := strconv.ParseFloat(c.DefaultQuery("min_amount", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode minAmount has to be a number: %v", err.Error())})
			return
		}

		maxAmount, err := strconv.ParseFloat(c.DefaultQuery("max_amount", "0"), 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode maxAmount has to be a number: %v", err.Error())})
			return
		}

		orders, err := s.store.FilterOrders(maker, taker, orderPair, secretHash, model.Status(status), minPrice, maxPrice, minAmount, maxAmount, page, perPage, verbose)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to get orders %s", err.Error()),
			})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func (s *Server) nonce() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"nonce": siwe.GenerateNonce(),
		})
	}
}

func (s *Server) verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := model.VerifySiwe{}
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := s.auth.Verify(req)
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

func (s *Server) secrets() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lastUpdatedUnixNano, err := strconv.ParseInt(ctx.DefaultQuery("lastUpdated", "0"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode last_updated has to be a number: %v", err.Error())})
			return
		}
		lastUpdated := time.Unix(0, lastUpdatedUnixNano)
		secrets, err := s.store.GetSecrets(lastUpdated)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to get secrets %s", err.Error()),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"secrets": secrets})
	}
}

func ParseGetOrderMessage(message string) (string, string, bool) {
	messageRegex := regexp.MustCompile(GetOrderMessageRegex)
	if !messageRegex.MatchString(message) {
		return "", "", false
	}
	matches := messageRegex.FindStringSubmatch(message)
	if len(matches) != 3 {
		return "", "", false
	}
	return matches[1], matches[2], true
}

func ParseGetOrdersMessage(message string) (string, string, bool) {
	messageRegex := regexp.MustCompile(GetOrdersMessageRegex)
	if !messageRegex.MatchString(message) {
		return "", "", false
	}
	matches := messageRegex.FindStringSubmatch(message)
	if len(matches) != 3 {
		return "", "", false
	}
	return matches[1], matches[2], true
}
