package rest

import (
	"sync"

	"github.com/catalogfi/orderbook/model"
)

type socketPool struct {
	mu                *sync.RWMutex
	updatedOrdersPool map[string][]chan UpdatedOrders
	OpenOrdersPool    map[string][]chan OpenOrders
	orderUpdatesPool  map[uint][]chan UpdatedOrder
}

type SocketPool interface {
	FilterAndBufferOrder(order model.Order) error
	AddUpdatedOrdersChannel(creator string, channel chan UpdatedOrders)
	AddOpenOrdersChannel(orderPair string, channel chan OpenOrders)
	AddOrderUpdatesChannel(id uint, channel chan UpdatedOrder)
	RemoveUpdatedOrdersChannel(creator string, channel chan UpdatedOrders)
	RemoveOpenOrdersChannel(orderPair string, channel chan OpenOrders)
	RemoveOrderUpdatesChannel(id uint, channel chan UpdatedOrder)
}

func NewSocketPool() SocketPool {
	return &socketPool{
		mu:                new(sync.RWMutex),
		updatedOrdersPool: make(map[string][]chan UpdatedOrders),
		OpenOrdersPool:    make(map[string][]chan OpenOrders),
		orderUpdatesPool:  make(map[uint][]chan UpdatedOrder),
	}
}

func (s *socketPool) FilterAndBufferOrder(order model.Order) error {

	users := []string{order.Maker}
	if order.Taker != "" {
		users = append(users, order.Taker)
	} else {
		s.bufferOpenOrders(order.OrderPair, []model.Order{order})
	}
	s.bufferUpdatedOrders(users, []model.Order{order})
	s.bufferOrderUpdates(order.ID, order)
	return nil
}
func (s *socketPool) bufferUpdatedOrders(users []string, orders []model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range users {
		for _, chann := range (s.updatedOrdersPool)[user] {
			chann <- UpdatedOrders{
				Orders: orders,
			}
		}
	}
}
func (s *socketPool) bufferOrderUpdates(orderId uint, order model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, chann := range (s.orderUpdatesPool)[orderId] {
		chann <- UpdatedOrder{
			Order: order,
		}

	}
}
func (s *socketPool) bufferOpenOrders(orderPair string, orders []model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, chann := range (s.OpenOrdersPool)[orderPair] {
		chann <- OpenOrders{
			Orders: orders,
		}

	}
}
func (s *socketPool) AddUpdatedOrdersChannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	(s.updatedOrdersPool)[creator] = append((s.updatedOrdersPool)[creator], channel)
}
func (s *socketPool) AddOpenOrdersChannel(orderPair string, channel chan OpenOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	(s.OpenOrdersPool)[orderPair] = append((s.OpenOrdersPool)[orderPair], channel)
}

func (s *socketPool) AddOrderUpdatesChannel(id uint, channel chan UpdatedOrder) {
	s.mu.Lock()
	defer s.mu.Unlock()
	(s.orderUpdatesPool)[id] = append((s.orderUpdatesPool)[id], channel)
}

func (s *socketPool) RemoveUpdatedOrdersChannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for m, n := range (s.updatedOrdersPool)[creator] {
		if n == channel {
			(s.updatedOrdersPool)[creator] = append((s.updatedOrdersPool)[creator][0:m], (s.updatedOrdersPool)[creator][m+1:len((s.updatedOrdersPool)[creator])]...)
			return
		}
	}
}
func (s *socketPool) RemoveOpenOrdersChannel(orderPair string, channel chan OpenOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for m, n := range (s.OpenOrdersPool)[orderPair] {
		if n == channel {
			(s.OpenOrdersPool)[orderPair] = append((s.OpenOrdersPool)[orderPair][0:m], (s.OpenOrdersPool)[orderPair][m+1:len((s.OpenOrdersPool)[orderPair])]...)
			return
		}
	}
}
func (s *socketPool) RemoveOrderUpdatesChannel(id uint, channel chan UpdatedOrder) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for m, n := range (s.orderUpdatesPool)[id] {
		if n == channel {
			(s.orderUpdatesPool)[id] = append((s.orderUpdatesPool)[id][0:m], (s.orderUpdatesPool)[id][m+1:len((s.orderUpdatesPool)[id])]...)
			return
		}
	}
}
