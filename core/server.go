package core

import (
	"runtime"
	"sync"

	"gunplan.top/concurrentNet/threading"

	"gunplan.top/concurrentNet/config"
)

func NewConcurrentNet() Server {
	s := ServerImpl{}
	return s.Init()
}

type ChannelInCallback func(c Channel, p Pipeline)

type Server interface {
	OnChannelConnect(ChannelInCallback) Server
	Option(*config.GetConfig) Server
	WaitType(config.WaitType) Server
	RegObserve(ServerObserve) Server
	Stop()
	Sync() error
	Join()
}

type ServerImpl struct {
	BaseServer
	comp   threading.Future
	n      []NetworkInet64
	t      threading.ThreadPool
	lk     sync.Mutex
	status ServerStatus
	alp    *eventLoop
}

func (s *ServerImpl) Init() Server {
	s.BaseServer.init(s)
	s.t = threading.CreateNewThreadPool(1, 0, 0, 1, nil)
	s.o = &DefaultObserve{}
	return s
}

func (s *ServerImpl) Stop() {
	s.o.OnStopping()
	s.status = STOPPING
	s.lk.Lock()
	s.alp.stop()
	s.status = STOPPED
	s.lk.Unlock()
	s.o.OnStopped()
}

func (s *ServerImpl) Join() {
	s.comp.Get()
}

func (s *ServerImpl) Sync() error {
	s.lk.Lock()
	defer s.lk.Unlock()
	s.status = BOOTING
	s.o.OnBooting()

	s.comp = s.t.Execwpr(func(i ...interface{}) interface{} {
		return (i[0].(func() error))()
	}, s.startLoops)

	s.status = RUNNING
	s.o.OnRunning(s.n)

	if s.cfj.WaitType == config.SYNC {
		s.Join()
	}
	return nil
}

func (s *ServerImpl) startLoops() error {
	var err error = nil

	defer func(err0 error) {
		s.failClean(err0)
	}(err)

	s.alp, err = NewEventLoop(runtime.NumCPU(), s.c)
	if err != nil {
		return err
	}

	err = s.alp.start()
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerImpl) failClean(err error) {
	if err == nil && s.alp != nil {
		s.alp.close()
	}
}
