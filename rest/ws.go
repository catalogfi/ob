package rest

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/catalogfi/orderbook/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) socket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade get request to websocket protocol
		ctx, cancel := context.WithCancel(context.Background())
		mx := new(sync.RWMutex)
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to upgrade to websocket %v", err)})
			cancel()
			return
		}
		pinger := time.NewTicker(time.Second * 60)
		defer func() {
			cancel()
			pinger.Stop()
			ws.Close()
		}()

		go func() {
			for range pinger.C {
				mx.Lock()
				err = ws.WriteJSON(map[string]interface{}{
					"type": "ping",
					"msg":  "ping",
				})
				mx.Unlock()
				if err != nil {
					s.logger.Debug("failed to write ping message", zap.Error(err))
					cancel()
					return
				}
			}
		}()
		for {
			// Read Message from client
			_, message, err := ws.ReadMessage()
			if err != nil {
				s.logger.Debug("failed to read a message", zap.Error(err))
				cancel()
				return
			}
			subscription := s.subscribe(message, ctx)

			go func() {
				for resp := range subscription {

					mx.Lock()
					err = ws.WriteJSON(map[string]interface{}{
						"type": fmt.Sprintf("%T", resp),
						"msg":  resp,
					})
					mx.Unlock()
					if err != nil {
						s.logger.Debug("failed to write message", zap.Error(err))
						cancel()
						return
					}
				}
			}()
		}
	}
}

func (s *Server) subscribe(msg []byte, ctx context.Context) <-chan interface{} {
	responses := make(chan interface{})
	fmt.Println("subscribing to ", string(msg))

	go func() {
		defer func() {
			close(responses)
		}()

		values := strings.Split(string(msg), "::")
		if len(values) != 2 || strings.ToLower(values[0]) != "subscribe" {
			responses <- WebsocketError{Code: 1, Error: fmt.Sprintf("invalid message %s", msg)}
			return
		}

		isOnlyPendingPending, err := regexp.Match("^0x[0-9a-fA-F]{40}-onlyPending$", []byte(values[1]))
		if err == nil && isOnlyPendingPending {
			addrValues := strings.Split(string(values[1]), "-")
			if len(addrValues) != 2 {
				responses <- WebsocketError{Code: 1, Error: fmt.Sprintf("invalid subscribe message %s", values[1])}
				return
			}
			for order := range s.subscribeToPendingAndUpdatedOrders(strings.ToLower(addrValues[0]), ctx) {
				responses <- order
			}
			return
		}

		isAddress, err := regexp.Match("^0x[0-9a-fA-F]{40}$", []byte(values[1]))
		if err == nil && isAddress {
			for order := range s.subscribeToUpdatedOrders(strings.ToLower(values[1]), ctx) {
				responses <- order
			}
			return
		}

		isOrderID, err := regexp.Match("[0-9]+$", []byte(values[1]))
		if err == nil && isOrderID {
			orderID, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				responses <- WebsocketError{Code: 2, Error: fmt.Sprintf("failed to parse order id %s: %v", values[1], err)}
				return
			}
			for order := range s.subscribeToOrderUpdates(uint(orderID), ctx) {
				responses <- order
			}
			return
		}

		isOrderPair, err := regexp.Match("^[a-zA-Z_]+(:0x[0-9a-fA-F]{40})?-[a-zA-Z_]+(:0x[0-9a-fA-F]{40})?", []byte(values[1]))
		if err == nil && isOrderPair {
			for response := range s.subscribeToOpenOrders(values[1], ctx) {
				responses <- response
			}
			return
		}
		responses <- WebsocketError{Code: 3, Error: fmt.Sprintf("invalid subscribe message %s", values[1])}
	}()
	return responses
}

type WebsocketError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func (s *Server) subscribeToOrderUpdates(id uint, ctx context.Context) <-chan UpdatedOrder {
	responses := make(chan UpdatedOrder)
	go func() {
		defer func() {
			s.socketPool.RemoveOrderUpdatesChannel(id, responses)
			close(responses)
		}()

		s.socketPool.AddOrderUpdatesChannel(id, responses)
		order, err := s.store.GetOrder(id)
		if err != nil {
			responses <- UpdatedOrder{Error: fmt.Sprintf("failed to get orders for %s: %v", id, err)}
			s.logger.Error("failed to get order", zap.Error(err))
			return
		}

		currentState := UpdatedOrder{
			Order: *order,
		}
		responses <- currentState

		<-ctx.Done()
	}()
	return responses
}

func (s *Server) subscribeToUpdatedOrders(creator string, ctx context.Context) <-chan UpdatedOrders {
	responses := make(chan UpdatedOrders)

	go func() {
		defer func() {
			s.socketPool.RemoveUpdatedOrdersChannel(creator, responses)
			close(responses)
		}()

		s.socketPool.AddUpdatedOrdersChannel(creator, responses)
		orders, err := s.store.GetOrdersByAddress(creator)
		if err != nil {
			responses <- UpdatedOrders{Error: fmt.Sprintf("failed to get orders for %s: %v", creator, err)}
			s.logger.Error("failed to get orders", zap.Error(err))
			return
		}

		newOrders := UpdatedOrders{
			Orders: orders,
		}
		responses <- newOrders

		<-ctx.Done()

	}()
	return responses
}
func (s *Server) subscribeToPendingAndUpdatedOrders(creator string, ctx context.Context) <-chan UpdatedOrders {
	responses := make(chan UpdatedOrders)

	go func() {
		defer func() {
			s.socketPool.RemoveUpdatedOrdersChannel(creator, responses)
			close(responses)
		}()

		s.socketPool.AddUpdatedOrdersChannel(creator, responses)
		orders, err := s.store.GetPendingOrdersForAddress(creator)
		if err != nil {
			responses <- UpdatedOrders{Error: fmt.Sprintf("failed to get orders for %s: %v", creator, err)}
			s.logger.Error("failed to get orders", zap.Error(err))
			return
		}

		newOrders := UpdatedOrders{
			Orders: orders,
		}
		responses <- newOrders

		<-ctx.Done()

	}()
	return responses
}

func (s *Server) subscribeToOpenOrders(orderPair string, ctx context.Context) <-chan OpenOrders {
	responses := make(chan OpenOrders)
	go func() {
		defer func() {
			s.socketPool.RemoveOpenOrdersChannel(orderPair, responses)
			close(responses)
		}()

		s.socketPool.AddOpenOrdersChannel(orderPair, responses)
		orders, err := s.store.FilterOrders("", "", orderPair, "", model.Created, 0.0, 0.0, 0.0, 0.0, 0, 0, true)
		if err != nil {
			responses <- OpenOrders{Error: fmt.Sprintf("failed to get orders for %s: %v", orderPair, err)}
			s.logger.Error("failed to get open orders", zap.Error(err))
			return
		}

		newOrders := OpenOrders{
			Orders: orders,
		}
		responses <- newOrders

		<-ctx.Done()
	}()
	return responses
}

type OpenOrders struct {
	Orders []model.Order `json:"orders"`
	Error  string        `json:"error"`
}

type UpdatedOrders struct {
	Orders []model.Order `json:"orders"`
	Error  string        `json:"error"`
}

type UpdatedOrder struct {
	Order model.Order `json:"order"`
	Error string      `json:"error"`
}
