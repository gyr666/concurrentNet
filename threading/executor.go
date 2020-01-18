package threading

type Executor interface {
	Exec(func())
}

type PoolExecutor interface {
	Executor
	Execwr(func() interface{}) Future
	Execwpr(func(...interface{}) interface{}, ...interface{}) Future
	Execwp(func(...interface{}),...interface{})
}
type Booter interface {
	Shutdown()
	ShutdownNow()
}

type Task struct {
	param []interface{}
	rev   chan interface{}
	t     ExecTask
}

func (t *Task) init() *Task{
	t.rev = make(chan interface{},1)
	return t
}
type ExecTask func(...interface{}) interface{}

type Future interface{
	get() interface{}
}
