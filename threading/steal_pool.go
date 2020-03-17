package threading

import (
	"sync"
	"time"

	"gunplan.top/concurrentNet/util"
)

type StealPool interface {
	Launcher
	PoolExecutor
	Status() PoolState
	Init(core int, w int, strategy func(interface{}))
	Join()
}

type stealPoolImpl struct {
	BasePool
	status         PoolStatusTransfer
	s              func(interface{})
	core           int
	index          util.Sequence
	workQueues     []chan *Task
	controlChannel chan ControlType
	g              sync.WaitGroup
	w              sync.Mutex
}

func (t *stealPoolImpl) Init(core int, w int, strategy func(interface{})) {
	t.BasePool.Init(t.addQueue0)
	t.core = core
	t.workQueues = make([]chan *Task, t.core)
	for i := range t.workQueues {
		t.workQueues[i] = make(chan *Task, w)
	}
	t.controlChannel = make(chan ControlType, t.core)
	t.s = strategy
	t.index = util.Sequence{Max: core}
	t.status = PoolStatusTransfer{}
	t.g.Add(t.core)
}

func (t *stealPoolImpl) Boot() {
	t.status.whenThreadBooting()
	for i := 0; i < t.core; i++ {
		go t.LaunchWork(i)
	}
	t.status.whenThreadBooted()
}

func (t *stealPoolImpl) LaunchWork(i int) {
	for {
		select {
		case task, ok := <-t.workQueues[i]:
			// core execute unit
			if ok {
				task.rev <- task.t(task.param...)
			}
			// control unit
		case op := <-t.controlChannel:
			switch op {

			case SHUTDOWN:
				t.controlChannel <- op
				if len(t.workQueues[i]) != 0 {
					t.consumeRemain(i)
				}
				fallthrough

			case ShutdownNow:
				t.controlChannel <- op
				t.g.Done()
				return
			}
		default:
			t.consumeOther()
		}
	}
}

func (t *stealPoolImpl) consumeOther() {
	for {
		select {
		case task, ok := <-t.workQueues[t.index.Next()]:
			if ok {
				task.rev <- task.t(task.param...)
			}
		case <-time.After(10 * time.Microsecond):
			continue
		}
	}

}

func (t *stealPoolImpl) consumeRemain(i int) {
	for len(t.workQueues[i]) != 0 {
		// The necessity og locking here is that
		// we have to make sure operator get length
		// and operator consume the channel is an
		// atomic operation.
		t.w.Lock()
		if len(t.workQueues[i]) != 0 {
			task := <-t.workQueues[i]
			t.w.Unlock()
			task.rev <- task.t(task.param...)
		} else {
			t.w.Unlock()
		}
	}
	close(t.workQueues[i])
}

func (t *stealPoolImpl) ShutdownNow() {
	for i := range t.workQueues {
		close(t.workQueues[i])
	}
	t.waitStop(ShutdownNow)
}

func (t *stealPoolImpl) Shutdown() {
	t.waitStop(SHUTDOWN)
}

func (t *stealPoolImpl) Join() {
	t.g.Wait()
}
func (t *stealPoolImpl) waitStop(c ControlType) {
	t.status.whenThreadStopping()
	t.controlChannel <- c
	go func() {
		t.g.Wait()
		close(t.controlChannel)
		t.status.whenThreadStopped()
	}()
}

func (t *stealPoolImpl) addQueue0(task *Task) {
	t.workQueues[t.index.Next()] <- task
}
