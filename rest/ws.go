package rest

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) socket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// upgrade get request to websocket protocol
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to upgrade to websocket %v", err)})
			return
		}
		defer ws.Close()

		pinger := time.NewTicker(time.Second * 60)
		for {
			// Read Message from client
			_, message, err := ws.ReadMessage()
			if err != nil {
				s.logger.Debug("failed to read a message", zap.Error(err))
				err = ws.WriteJSON(map[string]interface{}{
					"type":  fmt.Sprintf("%T", WebsocketError{Code: 101, Error: err.Error()}),
					"error": WebsocketError{Code: 101, Error: err.Error()},
				})
				if err != nil {
					ws.Close()
				}
			}
			subscription := s.subscribe(message)

			go func() {
				for {
					select {
					case resp, ok := <-subscription:
						if !ok {
							return
						}
						ws.WriteJSON(map[string]interface{}{
							"type": fmt.Sprintf("%T", resp),
							"msg":  resp,
						})
						if err != nil {
							s.logger.Debug("failed to write message", zap.Error(err))
							ws.Close()
						}
					case <-pinger.C:
						err = ws.WriteJSON(map[string]interface{}{
							"type": "ping",
							"msg":  "ping",
						})
						if err != nil {
							s.logger.Debug("failed to write ping message", zap.Error(err))
							ws.Close()
						}
					}
				}
			}()
		}
	}
}

func (s *Server) subscribe(msg []byte) <-chan interface{} {
	responses := make(chan interface{})
	fmt.Println("subscribing to ", string(msg))

	go func() {
		defer close(responses)

		values := strings.Split(string(msg), "_")
		if len(values) != 2 || strings.ToLower(values[0]) != "subscribe" {
			responses <- WebsocketError{Code: 1, Error: fmt.Sprintf("invalid message %s", msg)}
			return
		}

		isAddress, err := regexp.Match("0x[0-9a-fA-F]{40}", []byte(values[1]))
		if err == nil && isAddress {
			for order := range s.subscribeToUpdatedOrders(values[1]) {
				responses <- order
			}
			return
		}

		isOrderID, err := regexp.Match("[0-9]+", []byte(values[1]))
		if err == nil && isOrderID {
			orderID, err := strconv.ParseUint(values[1], 10, 64)
			if err != nil {
				responses <- WebsocketError{Code: 2, Error: fmt.Sprintf("failed to parse order id %s: %v", values[1], err)}
				return
			}
			for order := range s.subscribeToOrderUpdates(uint(orderID)) {
				responses <- order
			}
			return
		}

		isOrderPair, err := regexp.Match("^[a-zA-Z:]+-[a-zA-Z:]+", []byte(values[1]))
		if err == nil && isOrderPair {
			for response := range s.subscribeToOpenOrders(values[1]) {
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

func (s *Server) subscribeToOrderUpdates(id uint) <-chan UpdatedOrder {
	responses := make(chan UpdatedOrder)
	go func() {
		defer close(responses)

		order, err := s.store.GetOrder(id)
		if err != nil {
			responses <- UpdatedOrder{Error: fmt.Sprintf("failed to get order %d: %v", id, err)}
			s.logger.Error("failed to get order", zap.Error(err))
			return
		}
		responses <- UpdatedOrder{Order: *order}
		if order.Status >= model.Executed {
			return
		}

		for {
			newOrder, err := s.store.GetOrder(id)
			if err != nil {
				responses <- UpdatedOrder{Error: fmt.Sprintf("failed to get order %d: %v", id, err)}
				s.logger.Error("failed to get order", zap.Error(err))
				return
			}
			if order.Status == newOrder.Status && order.FollowerAtomicSwap.Status == newOrder.FollowerAtomicSwap.Status && order.InitiatorAtomicSwap.Status == newOrder.InitiatorAtomicSwap.Status {
				time.Sleep(time.Second * 2)
				continue
			}
			order = newOrder
			fmt.Println(order)
			responses <- UpdatedOrder{Order: *order}
			if order.Status >= model.Executed {
				return
			}
		}
	}()
	return responses
}

func (s *Server) subscribeToUpdatedOrders(creator string) <-chan UpdatedOrders {
	responses := make(chan UpdatedOrders)

	go func() {
		defer close(responses)
		orderMap := map[uint]model.Order{}

		for {
			orders, err := s.store.GetOrdersByAddress(creator)
			if err != nil {
				responses <- UpdatedOrders{Error: fmt.Sprintf("failed to get orders for %s: %v", creator, err)}
				s.logger.Error("failed to get order", zap.Error(err))
				return
			}

			// hasUpdated will always be true
			newOrders, _ := updatedOrders(orderMap, orders)
			responses <- newOrders

			for {
				newOrdersByAddr, err := s.store.GetOrdersByAddress(creator)
				if err != nil {
					responses <- UpdatedOrders{Error: fmt.Sprintf("failed to get orders for %s: %v", creator, err)}
					s.logger.Error("failed to get order", zap.Error(err))
					return
				}

				newOrders, hasUpdated := updatedOrders(orderMap, newOrdersByAddr)
				if !hasUpdated {
					time.Sleep(time.Second * 2)
					continue
				}

				responses <- newOrders
			}
		}
	}()
	return responses
}

func (s *Server) subscribeToOpenOrders(orderPair string) <-chan OpenOrder {
	responses := make(chan OpenOrder)
	go func() {
		defer close(responses)
		processed := map[uint]bool{}

		for {
			orders, err := s.store.FilterOrders("", "", "", orderPair, "", model.Created, 0.0, 0.0, 0.0, 0.0, 0, 0, true)
			if err != nil {
				responses <- OpenOrder{Error: fmt.Sprintf("failed to get orders on pair %s: %v", orderPair, err)}
				s.logger.Error("failed to get open orders", zap.Error(err))
				break
			}

			for _, order := range orders {
				if _, ok := processed[order.ID]; ok {
					continue
				}

				processed[order.ID] = true
				responses <- OpenOrder{Order: order}
			}
		}
	}()
	return responses
}

type OpenOrder struct {
	Order model.Order `json:"order"`
	Error string      `json:"error"`
}

type UpdatedOrders struct {
	Orders []model.Order `json:"orders"`
	Error  string        `json:"error"`
}

type UpdatedOrder struct {
	Order model.Order `json:"order"`
	Error string      `json:"error"`
}

func updatedOrders(orders map[uint]model.Order, newOrders []model.Order) (UpdatedOrders, bool) {
	hasUpdated := false
	// Remove unchanged orders
	for i := 0; i < len(newOrders); i++ {
		exist, ok := orders[newOrders[i].ID]
		if !ok || isDifferent(exist, newOrders[i]) {
			orders[newOrders[i].ID] = newOrders[i]
			hasUpdated = true
		} else {
			newOrders = append(newOrders[:i], newOrders[i+1:]...)
			i--
		}
	}

	fmt.Println(newOrders, hasUpdated)
	return UpdatedOrders{Orders: newOrders}, hasUpdated
}

func isDifferent(a, b model.Order) bool {
	return a.Status != b.Status ||
		(a.FollowerAtomicSwap.CurrentConfirmations != b.FollowerAtomicSwap.CurrentConfirmations || a.InitiatorAtomicSwap.CurrentConfirmations != b.InitiatorAtomicSwap.CurrentConfirmations) ||
		(a.FollowerAtomicSwap.FilledAmount != b.FollowerAtomicSwap.FilledAmount || a.InitiatorAtomicSwap.FilledAmount != b.InitiatorAtomicSwap.FilledAmount)
}
