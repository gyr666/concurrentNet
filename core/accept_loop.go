package core

import (
	"log"
	"sync"

	"gunplan.top/concurrentNet/threading"

	"gunplan.top/concurrentNet/core/netpoll"
)

type eventLoop struct {
	poller    *netpoll.Poller
	lg        *ioLoopGroup
	listeners map[int]listener
	lk        sync.Mutex
	active    bool
	index     int
	t         threading.ThreadPool
	wg        sync.WaitGroup
	h         ChannelInCallback
}

func NewEventLoop(sum int, h ChannelInCallback) (*eventLoop, error) {
	poller, err := netpoll.NewPoller()

	if err != nil {
		return nil, err
	}

	lg := NewIOLoopGroup()
	err = lg.create(sum)
	if err != nil {
		return nil, err
	}

	slp := &eventLoop{
		h:         h,
		poller:    poller,
		listeners: make(map[int]listener),
		lg:        lg,
		index:     -1,
		t:         threading.CreateNewStealPool(sum, 2, nil),
	}
	return slp, nil
}

func (alp *eventLoop) start() error {
	alp.lg.iterate(func(loop *ioLoop) bool {
		alp.t.Exec(loop.start)
		return true
	})
	return alp.poller.Polling(alp.eventHandler)
}

func (alp *eventLoop) stop() {
	if err := alp.poller.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d,%v\n", alp.index, err)
	}
	alp.lg.iterate(func(lp *ioLoop) bool {
		lp.stop()
		return true
	})
	alp.wg.Wait()
}

func (alp *eventLoop) close() {
	for _, listener := range alp.listeners {
		listener.close()
	}
	alp.poller.Close()
	if alp.lg == nil {
		return
	}
	alp.lg.iterate(func(lp *ioLoop) bool {
		lp.close()
		return true
	})
}

func (alp *eventLoop) eventHandler(fd int, ev uint32) error {
	//TODO accept
	//nfd, sa, err := unix.Accept(fd)
	//if err != nil {
	//	if err == unix.EAGAIN {
	//		return nil
	//	}
	//	return err
	//}
	//if err := unix.SetNonblock(nfd, true); err != nil {
	//	return err
	//}
	//lp := alp.lg.next()
	//
	//channel := alp.factory.NewChannel(nfd, sa)
	//
	//lp.poller.Trigger(func() error {
	//	if err := lp.poller.AddReader(nfd); err != nil {
	//		log.Println(err)
	//		return nil
	//	}
	//	lp.channels[nfd] = channel
	//	return nil
	//})

	return nil
}
