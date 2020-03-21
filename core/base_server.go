package core

import "gunplan.top/concurrentNet/config"

type BaseServer struct {
	r   Server
	c   ChannelInCallback
	cfj config.Config
	o   ServerObserve
}

func (s *BaseServer) init(r Server) {
	s.r = r
}

func (s *BaseServer) OnChannelConnect(c ChannelInCallback) Server {
	s.c = c
	return s.r
}

func (s *BaseServer) Option(strategy *config.GetConfig) Server {
	s.cfj = *strategy.Get()
	return s.r
}

func (s *BaseServer) WaitType(w config.WaitType) Server {
	s.cfj.WaitType = w
	return s.r
}

func (s *BaseServer) RegObserve(o ServerObserve) Server {
	s.o = o
	return s.r
}
