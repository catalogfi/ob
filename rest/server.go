package rest

import (
	"fmt"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spruceid/siwe-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	router *gin.Engine
	store  Store
	auth   Auth
	config model.Config
	logger *zap.Logger
	secret string
}

type Store interface {
	// get value locked in the given chain for the given user
	GetValueLocked(config model.Config, chain model.Chain) (*big.Int, error)
	// create order
	CreateOrder(creator, sendAddress, receiveAddress, orderPair, sendAmount, receiveAmount, secretHash string, userWalletBTCAddress string, config model.Config) (uint, error)
	// fill order
	FillOrder(orderID uint, filler, sendAddress, receiveAddress string, config model.Config) error
	// get order by id
	GetOrder(orderID uint) (*model.Order, error)
	// cancel order by id
	CancelOrder(creator string, orderID uint) error
	// get all orders for the given user
	FilterOrders(maker, taker, orderPair, secretHash, sort string, status model.Status, minPrice, maxPrice float64, minAmount, maxAmount float64, page, perPage int, verbose bool) ([]model.Order, error)
}

func NewServer(store Store, config model.Config, logger *zap.Logger, secret string) *Server {
	childLogger := logger.With(zap.String("service", "rest"))
	return &Server{
		router: gin.Default(),
		store:  store,
		secret: secret,
		logger: childLogger,
		auth:   NewAuth(config),
		config: config,
	}
}

func (s *Server) Run(addr string) error {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	authRoutes := s.router.Group("/")
	authRoutes.Use(s.authenticateJWT)

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "online",
		})
	})
	s.router.GET("/ws/order", s.GetOrderBySocket())
	s.router.GET("/ws/orders", s.GetOrdersSocket())
	s.router.GET("/orders/:id", s.GetOrder())
	s.router.GET("/orders", s.GetOrders())
	s.router.GET("/nonce", s.Nonce())
	s.router.GET("/assets", s.SupportedAssets())
	s.router.GET("/chains/:chain/value", s.GetValueByChain())
	s.router.POST("/verify", s.Verify())
	{
		authRoutes.POST("/orders", s.PostOrders())
		authRoutes.PUT("/orders/:id", s.FillOrder())
		authRoutes.DELETE("/orders/:id", s.CancelOrder())
	}
	return s.router.Run(addr)
}

type CreateOrder struct {
	SendAddress          string `json:"sendAddress" binding:"required"`
	ReceiveAddress       string `json:"receiveAddress" binding:"required"`
	OrderPair            string `json:"orderPair" binding:"required"`
	SendAmount           string `json:"sendAmount" binding:"required"`
	ReceiveAmount        string `json:"receiveAmount" binding:"required"`
	SecretHash           string `json:"secretHash" binding:"required"`
	UserWalletBTCAddress string `json:"userWalletBTCAddress" binding:"required"`
}

type Auth interface {
	Verify(req model.VerifySiwe) (*jwt.Token, error)
}

func (s *Server) authenticateJWT(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("missing authorization token")))
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
		ctx.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("invalid signing method")))
			return nil, fmt.Errorf("Invalid signing method")
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
		} else {
			s.logger.Debug("authorization failure", zap.Error(fmt.Errorf("invalid token claims")))
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

func (s *Server) GetOrderBySocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade get request to websocket protocol
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			s.logger.Error("failed to upgrade '/ws/order' to websocket")
			return
		}
		defer ws.Close()

		// Wait for the initiating message
		_, message, err := ws.ReadMessage()
		if err != nil {
			s.logger.Debug("read message", zap.Error(err))
			return
		}

		// Verify the user message
		action, orderStr, match := ParseGetOrderMessage(string(message))
		if !match {
			res := map[string]string{
				"error": "invalid action message",
			}
			ws.WriteJSON(res)
			return
		}
		orderID, err := strconv.ParseUint(orderStr, 10, 64)
		if err != nil {
			res := map[string]string{
				"error": "invalid order ID",
			}
			ws.WriteJSON(res)
			return
		}

		switch action {
		case "subscribe":
			// Fetch the current status of the order and send it to client
			order, err := s.store.GetOrder(uint(orderID))
			if err != nil {
				s.logger.Debug("get order", zap.String("api", "/ws/order"))
				return
			}
			if err := ws.WriteJSON(*order); err != nil {
				return
			}

			// Check the order status periodically and write updates to client
			ticker := time.NewTicker(15 * time.Second)

			for ; true; <-ticker.C {
				if order.Status >= model.OrderExecuted {
					break
				}
				if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
				newOrder, err := s.store.GetOrder(uint(orderID))
				if err != nil {
					s.logger.Error("get order", zap.String("api", "/ws/order"))
					continue
				}
				if newOrder.Status != order.Status {
					if err := ws.WriteJSON(*newOrder); err != nil {
						return
					}
					order = newOrder
				}
			}
		default:
			// ignore all unknown actions
		}
	}
}

func (s *Server) GetOrdersSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade get request to websocket protocol
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to upgrade to websocket %v", err)})
			return
		}
		defer ws.Close()

		// Read message from client
		_, message, err := ws.ReadMessage()
		if err != nil {
			s.logger.Debug("failed to read a message", zap.Error(err))
			return
		}

		// Verify the user message
		action, userAddr, match := ParseGetOrdersMessage(string(message))
		if !match {
			res := map[string]string{
				"error": "invalid action message",
			}
			ws.WriteJSON(res)
			return
		}

		switch action {
		case "subscribe":
			userAddr = strings.ToLower(userAddr)
			orders := map[uint]model.Order{}

			// Check the order status periodically and write updates to client
			ticker := time.NewTicker(15 * time.Second)

			for first := true; true; <-ticker.C {

				// Send a Ping message
				if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}

				// Fetch all orders has the userAddress
				makerOrders, err := s.store.FilterOrders(userAddr, "", "", "", "", model.Status(0), 0.0, 0.0, 0.0, 0.0, 0, 0, true)
				if err != nil && err != gorm.ErrRecordNotFound {
					s.logger.Error("load maker orders", zap.Error(err))
					continue
				}
				takerOrders, err := s.store.FilterOrders("", userAddr, "", "", "", model.Status(0), 0.0, 0.0, 0.0, 0.0, 0, 0, true)
				if err != nil && err != gorm.ErrRecordNotFound {
					s.logger.Error("load taker orders", zap.Error(err))
					continue
				}
				newOrders := append(makerOrders, takerOrders...)

				// Remove unchanged orders
				for i := 0; i < len(newOrders); i++ {
					exist, ok := orders[newOrders[i].ID]
					if !ok || !model.CompareOrder(exist, newOrders[i]) {
						orders[newOrders[i].ID] = newOrders[i]
					} else {
						newOrders = append(newOrders[:i], newOrders[i+1:]...)
						i--
					}
				}
				// Write all orders which has new updates (or the initial message after connection)
				if len(newOrders) != 0 || first {
					if err := ws.WriteJSON(newOrders); err != nil {
						return
					}
					first = false
				}
			}
		default:
			// ignore all unknown actions
		}
	}
}

func (s *Server) GetValueByChain() gin.HandlerFunc {
	return func(c *gin.Context) {
		chain, err := model.ParseChain(c.Param("chain"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "chain not supported"})
			return
		}

		valueLocked, err := s.store.GetValueLocked(s.config, chain)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"value": valueLocked,
		})
	}
}

func (s *Server) SupportedAssets() gin.HandlerFunc {
	assets := map[model.Chain][]model.Asset{}
	for chain, netConf := range s.config {
		assets[chain] = []model.Asset{}
		for asset := range netConf.Assets {
			assets[chain] = append(assets[chain], asset)
		}
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusCreated, assets)
	}
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

		oid, err := s.store.CreateOrder(strings.ToLower(creator.(string)), req.SendAddress, req.ReceiveAddress, req.OrderPair, req.SendAmount, req.ReceiveAmount, req.SecretHash, req.UserWalletBTCAddress, s.config)
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

func (s *Server) FillOrder() gin.HandlerFunc {
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

		if err := s.store.FillOrder(uint(orderID), strings.ToLower(filler.(string)), req.SendAddress, req.ReceiveAddress, s.config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to fill the Order %v", err.Error()),
			})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{})
	}
}

func (s *Server) GetOrder() gin.HandlerFunc {
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

func (s *Server) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		maker := c.DefaultQuery("maker", "")
		taker := c.DefaultQuery("taker", "")
		orderPair := c.DefaultQuery("order_pair", "")
		secretHash := c.DefaultQuery("secret_hash", "")
		orderBy := c.DefaultQuery("sort", "")

		maker = strings.ToLower(maker)
		taker = strings.ToLower(taker)

		verbose, err := strconv.ParseBool(c.DefaultQuery("verbose", "false"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode verbose has to be a boolean: %v", err.Error())})
			return
		}

		status, err := strconv.Atoi(c.DefaultQuery("status", "0"))
		if err != nil && status < int(model.Unknown) || status > int(model.OrderFailedSoft) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to decode status has to be a number between %d and %d", model.Unknown, model.OrderFailedSoft)})
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

		orders, err := s.store.FilterOrders(maker, taker, orderPair, secretHash, orderBy, model.Status(status), minPrice, maxPrice, minAmount, maxAmount, page, perPage, verbose)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("failed to get orders %s", err.Error()),
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
