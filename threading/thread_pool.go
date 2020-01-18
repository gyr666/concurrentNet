package threading

import "time"

type PoolState uint8
type ControlType uint8

const (
	STOP     PoolState = 1
	RUNNING  PoolState = 1 << 1
	WAITING  PoolState = 1 << 2
	SHUTDOWNING PoolState = 1 << 3
)

const (
	STOPANY  ControlType = 1
	STOPALL  ControlType = 1 << 1
)

type FutureImpl struct{
	wait *chan interface{}
}

func (f *FutureImpl) get() interface{}{
	return <- *f.wait
}

func newFuture(c* chan interface{})Future{
	return &FutureImpl{wait:c}
}
func NewThreadPool(core, ext int, wait time.Time, strategy func(interface{})) ThreadPool  {
	tp := threadpoolImpl{}
	tp.Init(core,ext,wait,strategy)
	tp.Boot()
	return &tp
}
type ThreadPool interface {
	Booter
	PoolExecutor
	Status() PoolState
	Init(core, ext int, wait time.Time, strategy func(interface{}))
}

type threadpoolImpl struct {
	status    PoolState
	core      int
	workQueue chan *Task
	controlChannel chan ControlType
}

func (t *threadpoolImpl) Init(core, ext int, wait time.Time, strategy func(interface{})){
	t.core = core
	t.workQueue = make(chan *Task,1000)
	t.controlChannel = make(chan ControlType,10)
}
func (t *threadpoolImpl) Boot() {
	t.status = RUNNING
	for i := 0; i < t.core; i++ {
		go func() {
			for ;t.status == RUNNING; {
				select {
				case task := <-t.workQueue:
					task.rev <- task.t(task.param)
				case op   := <-t.controlChannel:
					if op == STOPANY {
						return
					} else if op == STOPALL {
						t.controlChannel <- STOPALL
						return
					}
				}
			}
		}()
	}

}

func (t *threadpoolImpl) Status() PoolState{
	return t.status
}

func (t *threadpoolImpl) Shutdown() {
	t.status = SHUTDOWNING
}

func (t *threadpoolImpl) ShutdownNow() {
	t.controlChannel <- STOPALL
}
func (t *threadpoolImpl) addQueue(task *Task) {
	switch t.status {
	case RUNNING:
		t.workQueue <- task
	case SHUTDOWNING:
		panic("pool has been close")
	}
}
func (t *threadpoolImpl) Exec(f func()) {
	t.addQueue(&Task{param: nil, t: func(i ...interface{}) interface{} {
		f()
		return nil
	}})
}
func (t *threadpoolImpl) Execwr(f func() interface{}) Future {
	tsk :=Task{param: nil, t: func(i ...interface{}) interface{} {
		return f()
	}}
	tsk.init()
	t.addQueue(&tsk)
	return newFuture(&tsk.rev)
}
func (t *threadpoolImpl) Execwp(f func(...interface{}), p ...interface{}) {
	t.addQueue(&Task{param: p, t: func(i ...interface{}) interface{} {
		f(i...)
		return nil
	}})
}

func (t *threadpoolImpl) Execwpr(f func(...interface{}) interface{}, p ...interface{}) Future {
	tsk :=Task{param: p, t: f}
	tsk.init()
	t.addQueue(&tsk)
	return newFuture(&tsk.rev)
}
