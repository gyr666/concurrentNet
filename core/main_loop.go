package core

import "gunplan.top/concurrentNet/core/netpoll"

type mainLoop struct {
	poller *netpoll.Poller
	channels map[int]ParentChannel
}

func NewMainLoop()(*mainLoop,error){
	poller,err:=netpoll.NewPoller()
	if err!=nil{
		return nil,err
	}
	slp:=&mainLoop{
		poller:   poller,
		channels: make(map[int]ParentChannel),
	}
	return slp,nil
}

func (mp *mainLoop)start(){

}

func (mp *mainLoop)stop(){

}

func (mp *mainLoop)eventHandler(){

}