package core

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/config"
	"log"
	"runtime"
	"sync"
)

type ChannelInCallback func(c Channel, p Pipeline)

type Server interface {
	OnChannelConnect(ChannelInCallback) Server
	Option(config.ConfigStrategy) Server
	WaitType(config.WaitType) Server
	RegObserve(ServerObserve) Server
	Stop()
	Sync() error
	Join()
}

type ServerImpl struct {
	u   chan uint8
	c   ChannelInCallback
	cfj config.Config
	n   []NetworkInet64
	o   ServerObserve
	s   bool
	once   sync.Once
	lk     sync.Mutex
	loop   *acceptLoop
	alloc  buffer.Allocator
	status ServerStatus
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

func (s *ServerImpl) Option(strategy config.GetConfig) Server {
	s.cfj = strategy.Get()
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

	s.loop.stop()
	s.loop.close()
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
		log.Println(err)
		return errBoot
	}

	s.lk.Lock()
	s.status = RUNNING
	s.lk.Unlock()
	s.o.OnRunning(s.n)

	s.Join()
	return nil
}

func (s *ServerImpl) startLoops() (err error) {
	var (
		lg  ioLoopGroup
		alp *acceptLoop
	)
	defer func() {
		if err != nil {
			if alp != nil {
				alp.close()
			} else {
				lg.iterate(func(lp *ioLoop) bool {
					lp.close()
					return true
				})
			}
		}
	}()

	var lp *ioLoop
	cpuNum := runtime.NumCPU()
	for i := 0; i < cpuNum; i++ {
		lp, err = NewIOLoop(i, s.alloc)
		if err != nil {
			return
		}
		lg.registe(lp)
	}

	alp, err = NewAcceptLoop(&lg)
	if err != nil {
		return
	}
	s.loop = alp
	err = s.loop.start()
	if err != nil {
		return
	}
	return nil
}

func (s *ServerImpl) startLoops() error {
	var lps []Loop
	cpuNum := runtime.NumCPU()
	for i := 0; i < cpuNum; i++ {
		slp, err := NewSubLoop()
		if err != nil {
			return err
		}
		lps = append(lps, slp)
	}
	mlp, err := NewMainLoop()
	if err != nil {
		return err
	}
	lps = append(lps, mlp)

	//To make mlp at the first of the loopGroup , when use iterate close loops , will close mlp first
	for i := len(lps) - 1; i >= 0; i-- {
		s.lg.registe(lps[i])
	}

	s.lg.iterate(func(lp Loop) bool {
		s.wg.Add(1)
		go func() {
			lp.start()
			s.wg.Done()
		}()
		return true
	})
	return nil
}
