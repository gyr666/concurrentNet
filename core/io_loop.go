package core

import (
	"log"
	"sync"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/core/netpoll"
)

type ioLoop struct {
	polling  *netpoll.Polling
	channels map[int]Channel
	index    int
	lk       sync.Mutex
	p        chan Channel
	alloc    buffer.Allocator
	l        Pipeline
}

func NewIOLoop(index int, alloc buffer.Allocator) *ioLoop {
	lp := &ioLoop{
		index:    index,
		channels: make(map[int]Channel),
		alloc:    alloc,
		p:        make(chan Channel, 200),
	}
	lp.polling = netpoll.NewPolling(lp.eventHandler)
	return lp
}

func (lp *ioLoop) start() {
	if err := lp.polling.Polling(); err != nil {
		log.Println(err)
	}
}

func (lp *ioLoop) stop() {
	if err := lp.polling.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d , %v", lp.index, err)
	}
	lp.close()
}

func (lp *ioLoop) close() {
	lp.polling.Close()
}

func (lp *ioLoop) addChannel(c Channel) {
	lp.p <- c
}

func (lp *ioLoop) eventHandler(fd int, events uint32) error {
	//if channel,ok:=lp.channels[fd];ok{
	//	//switch {
	//	//
	//	//}
	//}
	return nil
}

func (lp *ioLoop) Read(buffer.ByteBuffer, error) {

}
