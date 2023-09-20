package rest

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/catalogfi/wbtc-garden/model"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Listner struct {
	dsn        string
	socketPool SocketPool
	logger     *zap.Logger
}

type socketPool struct {
	mu   *sync.RWMutex
	pool map[string][]chan UpdatedOrders
}

func NewListner(dsn string,
	socketPool SocketPool,
	logger *zap.Logger) Listner {
	return Listner{
		dsn:        dsn,
		socketPool: socketPool,
		logger:     logger,
	}

}

func (listner *Listner) Start(pgChannel string) {

	logError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			listner.logger.Error(err.Error())
		}
	}
	listener := pq.NewListener(listner.dsn, 10*time.Second, time.Minute, logError)
	err := listener.Listen(pgChannel)

	listner.logger.Info("Started listening to postgres events...")

	if err != nil {
		panic(err)
	}
	for {
		listner.waitForEvent(listener)
	}
}

func (listner *Listner) waitForEvent(l *pq.Listener) {
	for {
		select {
		case n := <-l.Notify:
			listner.logger.Info(fmt.Sprint("Received data from channel [", n.Channel, "] :"))
			var order model.Order
			err := json.Unmarshal([]byte(n.Extra), &order)
			if err != nil {
				listner.logger.Error(fmt.Sprintf("Error processing JSON: %v", err))
				return
			}
			listner.logger.Info("received order:", zap.Any("order", order))
			listner.socketPool.FilterAndBufferOrder(order)
			return
		case <-time.After(90 * time.Second):
			listner.logger.Info("Received no events for 90 seconds, checking connection")
			go func() {
				l.Ping()
			}()
			return
		}
	}
}
