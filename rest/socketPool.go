package rest

import (
	"sync"

	"github.com/catalogfi/wbtc-garden/model"
)

type socketPool struct {
	mu   *sync.RWMutex
	pool map[string][]chan UpdatedOrders
}

type SocketPool interface {
	FilterAndBufferOrder(order model.Order) error
	AddSocketChannel(creator string, channel chan UpdatedOrders)
	RemoveSocketChannel(creator string, channel chan UpdatedOrders)
}

func NewSocketPool(pool map[string][]chan UpdatedOrders) SocketPool {
	return &socketPool{
		mu:   new(sync.RWMutex),
		pool: pool,
	}
}

func (s *socketPool) FilterAndBufferOrder(order model.Order) error {

	var users []string
	users = append(users, order.Maker)
	if order.Taker != "" {
		users = append(users, order.Taker)
	}
	var orders []model.Order
	orders = append(orders, order)

	s.bufferOrders(users, orders)
	return nil
}
func (s *socketPool) bufferOrders(users []string, orders []model.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range users {
		for _, chann := range (s.pool)[user] {
			chann <- UpdatedOrders{
				Orders: orders,
			}
		}
	}
}
func (s *socketPool) AddSocketChannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	(s.pool)[creator] = append((s.pool)[creator], channel)
}
func (s *socketPool) RemoveSocketChannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for m, n := range (s.pool)[creator] {
		if n == channel {
			(s.pool)[creator] = append((s.pool)[creator][0:m], (s.pool)[creator][m+1:len((s.pool)[creator])]...)
			return
		}
	}
}