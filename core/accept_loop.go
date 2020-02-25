package core

import (
	"gunplan.top/concurrentNet/core/netpoll"
	"log"
	"sync"
)

type acceptLoop struct {
	poller    *netpoll.Poller
	lg        *ioLoopGroup
	listeners map[int]listener
	lk        sync.Mutex
	active    bool
	index     int
	wg        sync.WaitGroup
}

func NewAcceptLoop(lg *ioLoopGroup) (*acceptLoop, error) {
	poller, err := netpoll.NewPoller()
	if err != nil {
		return nil, err
	}
	slp := &acceptLoop{
		poller:    poller,
		listeners: make(map[int]listener),
		lg:        lg,
		index:     -1,
	}
	return slp, nil
}

func (alp *acceptLoop) start() error {
	return nil
}

func (alp *acceptLoop) stop() {
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

func (alp *acceptLoop) close() {
	for _, listener := range alp.listeners {
		listener.close()
	}
	alp.poller.Close()
	alp.lg.iterate(func(lp *ioLoop) bool {
		lp.close()
		return true
	})

}

func (alp *acceptLoop) eventHandler(fd int, ev uint32) error {
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
