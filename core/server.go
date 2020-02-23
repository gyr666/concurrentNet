package core

import (
	"gunplan.top/concurrentNet/config"
	"log"
	"runtime"
	"sync"
)

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
	lg  subLoopGroup
	wg  sync.WaitGroup
	once sync.Once
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
	// Todo : check server status, only called when status in [BOOTING,RUNNING]
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
	s.once.Do(func(){
		s.u <- 1
		s.status = STOPPING

	})
	s.o.OnStopping()
}

func (s *ServerImpl) Join() {
	<-s.u
	s.stopLoops()
	s.wg.Wait()
	s.closeLoops()
	s.status = STOPPED
	s.o.OnStopped()
}

func (s *ServerImpl) Sync() uint8 {
	s.status=BOOTING
	s.o.OnBooting()
	if err:=s.startLoops();err != nil{
		log.Println(err)
		s.closeLoops()
		return -1
	}
	s.status=RUNNING
	s.o.OnRunning(s.n)
	s.Join()
	return 0
}

func (s *ServerImpl) startLoops() error {

	cpuNum := runtime.NumCPU()
	for i := 0; i < cpuNum; i++ {
		slp, err := NewSubLoop(i)
		if err != nil {
			return err
		}
		s.lg.registe(slp)
	}

	mlp, err := NewMainLoop()
	if err != nil {
		return err
	}
	s.lg.registe(mlp)

	//note : mainLoop at the last of the loopGroup
	s.lg.iterate(false,func(lp Loop) bool{
		s.wg.Add(1)
		go func() {
			lp.start()
			s.wg.Done()
		}()
		return true
	})
	return nil
}

func (s *ServerImpl) stopLoops() {
	//note : mainLoop at the last of the loopGroup , we should stop it first
	s.lg.iterate(true,func(lp Loop) bool {
		lp.stop()
		return true
	})
}

func (s *ServerImpl) closeLoops() {

	s.lg.iterate(true,func(lp Loop) bool {
		lp.close()
		return true
	})
}