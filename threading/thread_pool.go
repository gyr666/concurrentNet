package threading

import "time"

type PoolType uint8
type ControlType uint8

const (
	STOP     PoolType = 1
	RUNNING  PoolType = 1 << 1
	WAITING  PoolType = 1 << 2
	BLOCKING PoolType = 1 << 3
)
const (
	STOPANY  ControlType = 1
	STOPALL  ControlType = 1 << 1
)

type ThreadPool interface {
	Booter
	PoolExecutor
	Status() PoolType
	Init(core, ext int, wait time.Time, strategy func(interface{}))
}

type threadpoolImpl struct {
	status    uint8
	core      int
	workQueue chan *Task
	contorlChannel chan int
}

func (t *threadpoolImpl) Boot() {
	for i := 0; i < t.core; i++ {
		go func() {
			for {
				select {
				case task := <-t.workQueue:
					task.rev = task.t(task.param)
				case op   := <-t.contorlChannel:
					if op == STOPANY {
						return
					}else if == STOPALL {
						t.controlChannel <- STOPALL
						return
					}
				}
			}
		}()
	}

}

func (t *threadpoolImpl) Shutdown() {
	t.contorlChannel <- 
}
func (t *threadpoolImpl) addQueue(task *Task) {
	t.workQueue <- task
}
func (t *threadpoolImpl) Exec(f func()) {
	t.addQueue(&Task{param: nil, t: func(i ...interface{}) interface{} {
		f()
		return nil
	}})
}
func (t *threadpoolImpl) Execwr(f func() interface{}) Furture {
	t.addQueue(&Task{param: nil, t: func(i ...interface{}) interface{} {
		return f()
	}})
}
func (t *threadpoolImpl) EXecwp(f func(...interface{}), p ...interface{}) {
	t.addQueue(&Task{param: p, t: func(i ...interface{}) interface{} {
		f(i...)
		return nil
	}})
}
func (t *threadpoolImpl) EXecwpr(f func(...interface{}) interface{}, p ...interface{}) Furture {
	t.addQueue(&Task{param: p, t: f})
}
