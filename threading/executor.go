package threading

type Executor interface {
	// execute a task on `Executor`
	Exec(func())
}

// PoolExecutor is a pooled executor
type PoolExecutor interface {
	Executor
	// execute a task with result
	Execwr(func() interface{}) Future
	// execute a task with parameter and result
	Execwpr(func(...interface{}) interface{}, ...interface{}) Future
	// execute a task with parameter
	Execwp(func(...interface{}), ...interface{})
}

type Booter interface {
	Boot()
	Shutdown()
	ShutdownNow()
}

type Task struct {
	param []interface{}
	rev   chan interface{}
	t     ExecTask
}

func (t *Task) init() *Task {
	t.rev = make(chan interface{}, 1)
	return t
}

type ExecTask func(...interface{}) interface{}
