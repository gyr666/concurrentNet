package threading

import (
	"sync"
	"time"
)

type futureImpl struct {
	wait *chan interface{}
}

func (f *futureImpl) Get() interface{} {
	return <-*f.wait
}

func (f *futureImpl) isDone() bool {
	return len(*f.wait) == 1
}

func newFuture(c *chan interface{}) Future {
	return &futureImpl{wait: c}
}

func NewThreadPool(core, ext int, span time.Duration, w uint64, strategy func(interface{})) ThreadPool {
	tp := threadPoolImpl{}
	tp.Init(core, ext, span, w, strategy)
	tp.Boot()
	return &tp
}

type ThreadPool interface {
	Launcher
	PoolExecutor
	Status() PoolState
	Init(core, ext int, d time.Duration, w uint64, strategy func(interface{}))
	WaitForStop()
}

type threadPoolImpl struct {
	status         PoolState
	s              func(interface{})
	core           int
	workQueue      chan *Task
	controlChannel chan ControlType
	g              sync.WaitGroup
	w              sync.Mutex
}

func (t *threadPoolImpl) Init(core, ext int, span time.Duration, w uint64, strategy func(interface{})) {
	t.core = core
	t.workQueue = make(chan *Task, 1000)
	t.controlChannel = make(chan ControlType, t.core)
	t.s = strategy
	t.g.Add(t.core)
}
func (t *threadPoolImpl) Boot() {
	t.status = RUNNING
	for i := 0; i < t.core; i++ {
		go t.LaunchWork()
	}

}

func (t *threadPoolImpl) LaunchWork() {
	for {
		select {
		case task := <-t.workQueue:
			// core execute unit
			task.rev <- task.t(task.param...)
			// control unit
		case op := <-t.controlChannel:
			switch op {

			case SHUTDOWN:
				t.controlChannel <- op
				if len(t.workQueue) != 0 {
					t.consumeRemain()
				}
				fallthrough

			case SHUTDOWNNOW:
				t.controlChannel <- op
				fallthrough

			case STOPANY:
				t.g.Done()
				return
			}
		}
	}
}

func (t *threadPoolImpl) consumeRemain() {
	for len(t.workQueue) != 0 {
		// The necessity og locing here is that
		// we have to make sure operator get length
		// and operator consume the channel is an
		// atomic operation.
		t.w.Lock()
		if len(t.workQueue) != 0 {
			task := <-t.workQueue
			t.w.Unlock()
			task.rev <- task.t(task.param...)
		} else {
			t.w.Unlock()
		}
	}
}

func (t *threadPoolImpl) LaunchWorkExt() {
	t.g.Add(1)
	t.LaunchWork()
}

func (t *threadPoolImpl) Status() PoolState {
	return t.status
}

func (t *threadPoolImpl) WaitForStop() {
	t.g.Wait()
}

func (t *threadPoolImpl) Shutdown() {
	t.waitStop(SHUTDOWN)
}

func (t *threadPoolImpl) ShutdownNow() {
	close(t.workQueue)
	t.waitStop(SHUTDOWNNOW)
}

func (t *threadPoolImpl) ShutdownAny() {
	t.controlChannel <- STOPANY
}

func (t *threadPoolImpl) waitStop(c ControlType) {
	t.status = STOPPING
	t.controlChannel <- c
	go func() {
		t.g.Wait()
		close(t.controlChannel)
		t.status = STOPPED
	}()
}

func (t *threadPoolImpl) addQueue(task *Task) {
	switch t.status {
	case RUNNING:
		t.workQueue <- task
	case STOPPED:
		fallthrough
	case STOPPING:
		panic("pool has been close")
	}
}

func (t *threadPoolImpl) Exec(f func()) {
	t.addQueue(&Task{param: nil, t: func(i ...interface{}) interface{} {
		f()
		return nil
	}})
}

func (t *threadPoolImpl) Execwr(f func() interface{}) Future {
	tsk := Task{param: nil, t: func(i ...interface{}) interface{} {
		return f()
	}}
	tsk.init()
	t.addQueue(&tsk)
	return newFuture(&tsk.rev)
}

func (t *threadPoolImpl) Execwp(f func(...interface{}), p ...interface{}) {
	t.addQueue(&Task{param: p, t: func(i ...interface{}) interface{} {
		f(i...)
		return nil
	}})
}

func (t *threadPoolImpl) Execwpr(f func(...interface{}) interface{}, p ...interface{}) Future {
	tsk := Task{param: p, t: f}
	tsk.init()
	t.addQueue(&tsk)
	return newFuture(&tsk.rev)
}
