package rest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Listner struct {
	dsn        string
	socketPool SocketPool
	logger     *zap.Logger
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

func (listner *Listner) Start(ordersUpdatechan string, swapsUpdatechan string) {

	logError := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			listner.logger.Error(err.Error())
		}
	}
	ordersListener := pq.NewListener(listner.dsn, 10*time.Second, time.Minute, logError)
	err := ordersListener.Listen(ordersUpdatechan)
	if err != nil {
		panic(err)
	}

	swapsListener := pq.NewListener(listner.dsn, 10*time.Second, time.Minute, logError)
	err = swapsListener.Listen(swapsUpdatechan)
	if err != nil {
		panic(err)
	}

	listner.logger.Info("Started listening to postgres events...")

	for {
		listner.waitForEvent(ordersListener, swapsListener)
	}

}

// oid ->  orderid
// sid -> swap id
// on -> orders notification
// sn -> swaps notification
// ol -> orders listener
// sl -> swaps listener
func (listner *Listner) waitForEvent(ol *pq.Listener, sl *pq.Listener) {
	for {
		select {
		case on := <-ol.Notify:
			listner.logger.Info(fmt.Sprint("Received data from channel [", on.Channel, "] :"))
			oid, err := strconv.ParseUint(on.Extra, 10, 64)
			if err != nil {
				listner.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				return
			}
			listner.logger.Info("received order:", zap.Uint64("order id:", oid))
			err = listner.socketPool.FilterAndBufferOrder(oid)
			if err != nil {
				listner.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
			}
			return
		case sn := <-sl.Notify:
			listner.logger.Info(fmt.Sprint("Received data from channel [", sn.Channel, "] :"))
			sid, err := strconv.ParseUint(sn.Extra, 10, 64)
			if err != nil {
				listner.logger.Error(fmt.Sprintf("Error processing id: %v", err))
				return
			}
			if sid&1 == 1 {
				sid += 1
			}
			oid := sid >> 1
			listner.logger.Info("received order:", zap.Uint64("order id:", oid), zap.Uint64("swap id:", sid))
			err = listner.socketPool.FilterAndBufferOrder(oid)
			if err != nil {
				listner.logger.Error("Failed to write order to channel", zap.Uint64("order id:", oid))
			}
			return
		case <-time.After(90 * time.Second):
			listner.logger.Info("Received no events for 90 seconds, checking connection")
			go func() {
				ol.Ping()
				sl.Ping()
			}()
			return
		}
	}
}
