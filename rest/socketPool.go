package rest

import (
	"sync"

	"github.com/catalogfi/wbtc-garden/model"
)

type SocketPool interface {
	FilterAndBufferOrder(order model.Order)
	AddSocketChannel(creator string, channel chan UpdatedOrders)
	RemoveSocketchannel(creator string, channel chan UpdatedOrders)
}

func NewSocketPool(pool map[string][]chan UpdatedOrders) SocketPool {
	return &socketPool{
		mu:   new(sync.RWMutex),
		pool: pool,
	}
}

func (s *socketPool) FilterAndBufferOrder(order model.Order) {
	creator := order.Maker
	var orders []model.Order
	for _, chans := range (s.pool)[creator] {
		chans <- UpdatedOrders{
			Orders: append(orders, order),
		}
	}
}
func (s *socketPool) AddSocketChannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	(s.pool)[creator] = append((s.pool)[creator], channel)
	s.mu.Unlock()
}
func (s *socketPool) RemoveSocketchannel(creator string, channel chan UpdatedOrders) {
	s.mu.Lock()
	for m, n := range (s.pool)[creator] {
		if n == channel {
			(s.pool)[creator] = append((s.pool)[creator][0:m], (s.pool)[creator][m+1:len((s.pool)[creator])]...)
		}
	}
	s.mu.Unlock()
}
