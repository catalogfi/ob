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
}

func NewDBListener(dsn string,
	socketPool SocketPool,
	logger *zap.Logger) DBListener {
	return DBListener{
		dsn:        dsn,
		socketPool: socketPool,
		logger:     logger,
	}

}

func (listener *DBListener) Start(ordersUpdatechan string, swapsUpdatechan string) {

	logError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			listener.logger.Error(err.Error())
		}
	}
	ordersListener := pq.NewListener(listener.dsn, 10*time.Second, time.Minute, logError)
	err := ordersListener.Listen(ordersUpdatechan)
	if err != nil {
		panic(err)
	}

	swapsListener := pq.NewListener(listener.dsn, 10*time.Second, time.Minute, logError)
	err = swapsListener.Listen(swapsUpdatechan)
	if err != nil {
		panic(err)
	}

	listener.logger.Info("Started listening to postgres events...")

	listener.waitForEvent(ordersListener, swapsListener)

}

// oid ->  orderid
// sid -> swap id
// on -> orders notification
// sn -> swaps notification
// ol -> orders listener
// sl -> swaps listener
func (listener *DBListener) waitForEvent(ol *pq.Listener, sl *pq.Listener) {
	for {
		select {
		case on := <-ol.Notify:
			listener.logger.Info(fmt.Sprint("Received data from channel [", on.Channel, "] :"))
			oid, err := strconv.ParseUint(on.Extra, 10, 64)
			if err != nil {
				listener.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				continue
			}
			listener.logger.Info("received order:", zap.Uint64("order id:", oid))
			err = listener.socketPool.FilterAndBufferOrder(oid)
			if err != nil {
				listener.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
			}
		case sn := <-sl.Notify:
			listener.logger.Info(fmt.Sprint("Received data from channel [", sn.Channel, "] :"))
			sid, err := strconv.ParseUint(sn.Extra, 10, 64)
			if err != nil {
				listener.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				continue
			}
			if sid&1 == 1 {
				sid += 1
			}
			oid := sid >> 1
			listener.logger.Info("received order:", zap.Uint64("order id:", oid), zap.Uint64("swap id:", sid))
			err = listener.socketPool.FilterAndBufferOrder(oid)
			if err != nil {
				listener.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
			}
		case <-time.After(90 * time.Second):
			listener.logger.Info("Received no events for 90 seconds, checking connection")
			go func() {
				ol.Ping()
				sl.Ping()
			}()
		}
	}
}
