package core

import "gunplan.top/concurrentNet/config"

type ChannelInCallback func(c Channel, p Pipeline)

type Server interface {
	OnChannelConnect(ChannelInCallback) Server
	SetServerSocketChannel(ParentChannel) Server
	Option(config.ConfigStrategy) Server
	AddListen(*NetworkInet64) Server
	WaitType(config.WaitType) Server
	RegObserve(ServerObserve) Server
	Stop()
	Sync() uint8
	Join()
}

type ServerImpl struct {
	u   chan uint8
	c   ChannelInCallback
	i   ParentChannel
	cfj config.Config
	n   []NetworkInet64
	o   ServerObserve
	s   bool
}

func (s *ServerImpl) Init() Server {
	s.u = make(chan uint8, 1)
	s.o = &DefaultObserve{}
	return s
}

func (s *ServerImpl) RegObserve(o ServerObserve) Server {
	s.o = o
	return s
}

func (s *ServerImpl) OnChannelConnect(c ChannelInCallback) Server {
	s.c = c
	return s
}

func (s *ServerImpl) SetServerSocketChannel(i ParentChannel) Server {
	s.i = i
	return s
}

func (s *ServerImpl) Option(strategy config.ConfigStrategy) Server {
	s.cfj = config.Config{}
	strategy.Fill(&s.cfj)
	return s
}

func (s *ServerImpl) AddListen(n *NetworkInet64) Server {
	if s.i == nil {
		panic("please set parent channel")
	}
	s.n = append(s.n, *n)
	s.i.Listen(n)
	return s
}

func (s *ServerImpl) WaitType(w config.WaitType) Server {
	s.cfj.WaitType = w
	return s
}

func (s *ServerImpl) Stop() {
	s.o.OnStopping()
	s.u <- 1
	s.o.OnStopped()
}

func (s *ServerImpl) Join() {
	<-s.u
}

func (s *ServerImpl) Sync() uint8 {
	s.o.OnBooting()
	//todo use callback
	if s.s {
		s.o.OnBooted(s.n)
		s.Join()
	}
	s.o.OnBooted(s.n)
	return 0
}
