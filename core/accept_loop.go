package core

import (
	"gunplan.top/concurrentNet/config"
	"log"
	"sync"

	"gunplan.top/concurrentNet/threading"

	"gunplan.top/concurrentNet/core/netpoll"
)

type eventLoop struct {
	poll      netpoll.UnixPoll
	lg        *ioLoopGroup
	listeners map[int]listener
	lk        sync.Mutex
	active    bool
	index     int
	t         threading.ThreadPool
	wg        sync.WaitGroup
	h         ChannelInCallback
	sc        netpoll.ServerChannel
}

func NewEventLoop(sum int, h ChannelInCallback, c *config.Config) (*eventLoop, error) {
	sc := netpoll.NewServerChannel(c.Listen.Port, int(c.Backlog))
	slp := &eventLoop{
		sc:        sc,
		h:         h,
		poll:      netpoll.NewServerUnixPoll(sc),
		listeners: make(map[int]listener),
		lg:        NewIOLoopGroup().create(sum),
		index:     -1,
		t:         threading.CreateNewStealPool(sum, 2, nil),
	}
	slp.poll.Handle(nil)
	return slp, nil
}

func (alp *eventLoop) start() {
	alp.lg.iterate(func(loop *ioLoop) bool {
		alp.t.Exec(loop.start)
		return true
	})
	alp.poll.Select()
}

func (alp *eventLoop) stop() {
	if err := alp.polling.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d,%v\n", alp.index, err)
	}
	alp.lg.iterate(func(lp *ioLoop) bool {
		alp.t.Exec(lp.stop)
		alp.wg.Done()
		return true
	})
	alp.wg.Wait()
}

func (alp *eventLoop) addChannel(c Channel) {

}

func (alp *eventLoop) close() {
	for _, listener := range alp.listeners {
		listener.close()
	}
	alp.poll.Close()
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
	//lp.polling.Trigger(func() error {
	//	if err := lp.polling.AddReader(nfd); err != nil {
	//		log.Println(err)
	//		return nil
	//	}
	//	lp.channels[nfd] = channel
	//	return nil
	//})

	return nil
}
