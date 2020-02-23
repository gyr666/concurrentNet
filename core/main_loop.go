package core

import (
	"gunplan.top/concurrentNet/core/netpoll"
	"log"
	"sync"
)

type mainLoop struct {
	poller   *netpoll.Poller
	channels map[int]ParentChannel
	mux sync.Mutex
	active bool
	index int
}

func NewMainLoop() (*mainLoop, error) {
	poller, err := netpoll.NewPoller()
	if err != nil {
		return nil, err
	}
	slp := &mainLoop{
		poller:   poller,
		channels: make(map[int]ParentChannel),
		index:-1,
	}
	return slp, nil
}

func (mlp *mainLoop) start() {

}

func (mlp *mainLoop) stop() {
	if err := mlp.poller.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d,%v\n", mlp.index, err)
	}
}

func (mlp *mainLoop) close() {

	for _,channel:=range mlp.channels{
		channel.Close()
	}
	mlp.poller.Close()
}

func (mlp *mainLoop) eventHandler() {

}
