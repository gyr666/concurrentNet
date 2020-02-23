package core

import (
	"gunplan.top/concurrentNet/core/netpoll"
	"log"
)

type subLoop struct {
	poller   *netpoll.Poller
	channels map[int]ChildChannel
	index int
}

func NewSubLoop(index int) (*subLoop, error) {
	poller, err := netpoll.NewPoller()
	if err != nil {
		return nil, err
	}
	slp := &subLoop{
		index: index,
		poller:   poller,
		channels: make(map[int]ChildChannel),
	}
	return slp, nil
}

func (slp *subLoop) start() {

}

func (slp *subLoop) stop() {
	if err := slp.poller.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d , %v", slp.index, err)
	}
}

func (slp *subLoop) close() {
	slp.poller.Close()
}

func (slp *subLoop) eventHandler() {

}
