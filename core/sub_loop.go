package core

import (
	"gunplan.top/concurrentNet/core/netpoll"
)

type subLoop struct {
	poller   *netpoll.Poller
	channels map[int]ChildChannel
}

func NewSubLoop() (*subLoop, error) {
	poller, err := netpoll.NewPoller()
	if err != nil {
		return nil, err
	}
	slp := &subLoop{
		poller:   poller,
		channels: make(map[int]ChildChannel),
	}
	return slp, nil
}

func (slp *subLoop) start() {

}

func (slp *subLoop) stop() {
	slp.poller.Close()
}

func (slp *subLoop) eventHandler() {

}
