package core

import (
	"log"
	"runtime"
	"sync"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/config"
)

func NewConcurrentNet() Server {
	return &ServerImpl{}
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
	u      chan uint8
	c      ChannelInCallback
	cfj    config.Config
	n      []NetworkInet64
	o      ServerObserve
	s      bool
	once   sync.Once
	lk     sync.Mutex
	alloc  buffer.Allocator
	status ServerStatus
	ilg    *ioLoopGroup
	alp    *acceptLoop
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

func (s *ServerImpl) Option(strategy *config.GetConfig) Server {
	s.cfj = *strategy.Get()
	return s
}

func (s *ServerImpl) WaitType(w config.WaitType) Server {
	s.cfj.WaitType = w
	return s
}

func (s *ServerImpl) Stop() {
	s.once.Do(func() {
		s.u <- 1
		s.lk.Lock()
		s.status = STOPPING
		s.lk.Unlock()

	})
	s.o.OnStopping()
}

func (s *ServerImpl) Join() {
	<-s.u

	s.alp.stop()
	s.alp.close()
	if err := s.alloc.Destroy(); err != nil {
		log.Println(err)
	}

	s.lk.Lock()
	s.status = STOPPED
	s.lk.Unlock()
	s.o.OnStopped()
}

func (s *ServerImpl) Sync() error {
	s.lk.Lock()
	s.status = BOOTING
	s.lk.Unlock()
	s.o.OnBooting()

	s.alloc = buffer.NewLikedBufferAllocator()
	if err := s.startLoops(); err != nil {
		log.Fatal(err)
		return errBoot
	}

	s.lk.Lock()
	s.status = RUNNING
	s.lk.Unlock()
	s.o.OnRunning(s.n)
	s.Join()
	return nil
}

func (s *ServerImpl) startLoops() error {
	var err error = nil

	defer func(err0 error) {
		s.failClean(err0)
	}(err)

	var lp *ioLoop
	cpuNum := runtime.NumCPU()
	for i := 0; i < cpuNum; i++ {
		lp, err = NewIOLoop(i, s.alloc)
		if err != nil {
			return err
		}
		s.ilg.registe(lp)
	}

	s.alp, err = NewAcceptLoop(s.ilg)
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
	if err != nil {
		if s.alp != nil {
			s.alp.close()
		}
		if s.ilg != nil {
			s.ilg.iterate(func(lp *ioLoop) bool {
				lp.close()
				return true
			})
		}
	}
}
