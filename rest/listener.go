package rest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type DBListener struct {
	dsn        string
	socketPool SocketPool
	logger     *zap.Logger
	store      Store
}

func NewDBListener(dsn string,
	socketPool SocketPool,
	logger *zap.Logger, store Store) DBListener {
	return DBListener{
		dsn:        dsn,
		socketPool: socketPool,
		logger:     logger,
		store:      store,
	}

}

func (listener *DBListener) Start(ordersUpdatechan string, swapsUpdatechan string, ordersAddedChan string) {

	logError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			listener.logger.Error(err.Error())
		}
	}
	orderUpdatesListener := pq.NewListener(listener.dsn, 10*time.Second, time.Minute, logError)
	err := orderUpdatesListener.Listen(ordersUpdatechan)
	if err != nil {
		panic(err)
	}
	orderAddedListener := pq.NewListener(listener.dsn, 10*time.Second, time.Minute, logError)
	err = orderAddedListener.Listen(ordersAddedChan)
	if err != nil {
		panic(err)
	}

	swapsListener := pq.NewListener(listener.dsn, 10*time.Second, time.Minute, logError)
	err = swapsListener.Listen(swapsUpdatechan)
	if err != nil {
		panic(err)
	}

	listener.logger.Info("Started listening to postgres events...")

	listener.waitForEvent(orderUpdatesListener, orderAddedListener, swapsListener)

}

// oid ->  orderid
// sid -> swap id
// on -> orders notification
// sn -> swaps notification
// uol -> updated orders listener
// nol -> new orders listener
// sl -> swaps listener
func (listener *DBListener) waitForEvent(uol, nol, sl *pq.Listener) {
	for {
		select {
		case on := <-uol.Notify:
			if on == nil {
				continue
			}
			listener.logger.Info(fmt.Sprint("Received data from channel [", on.Channel, "] :"))
			oid, err := strconv.ParseUint(on.Extra, 10, 64)
			if err != nil {
				listener.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				continue
			}
			listener.logger.Info("received order:", zap.Uint64("order id:", oid))
			order, err := listener.store.GetOrder(uint(oid))
			if err != nil {
				listener.logger.Error("Failed to get order", zap.String("error", err.Error()), zap.Uint64("order id:", oid))
				continue
			}
			err = listener.socketPool.FilterAndBufferOrder(*order)
			if err != nil {
				listener.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
				continue
			}
		case on := <-nol.Notify:
			if on == nil {
				continue
			}
			listener.logger.Info(fmt.Sprint("Received data from channel [", on.Channel, "] :"))
			oid, err := strconv.ParseUint(on.Extra, 10, 64)
			if err != nil {
				listener.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				continue
			}
			listener.logger.Info("received order:", zap.Uint64("order id:", oid))
			order, err := listener.store.GetOrder(uint(oid))
			if err != nil {
				listener.logger.Error("Failed to get order", zap.String("error", err.Error()), zap.Uint64("order id:", oid))
				continue
			}
			err = listener.socketPool.FilterAndBufferOrder(*order)
			if err != nil {
				listener.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
				continue
			}
		case sn := <-sl.Notify:
			if sn == nil {
				continue
			}
			listener.logger.Info(fmt.Sprint("Received data from channel [", sn.Channel, "] :"))
			sid, err := strconv.ParseUint(sn.Extra, 10, 64)
			if err != nil {
				listener.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				continue
			}
			order, err := listener.store.GetOrderBySwapID(uint(sid))
			if err != nil {
				listener.logger.Error("Failed to get order by swap id", zap.String("error", err.Error()), zap.Uint64("swap id:", sid))
				continue
			}
			listener.logger.Info("received order:", zap.Uint64("order id:", uint64(order.ID)), zap.Uint64("swap id:", sid))
			err = listener.socketPool.FilterAndBufferOrder(*order)
			if err != nil {
				listener.logger.Error("Failed to write order to channel", zap.String("error", err.Error()), zap.Uint64("order id:", uint64(order.ID)))
			}
		case <-time.After(90 * time.Second):
			listener.logger.Info("Received no events for 90 seconds, checking connection")
			go func() {
				uol.Ping()
				nol.Ping()
				sl.Ping()
			}()
		}
	}
}
