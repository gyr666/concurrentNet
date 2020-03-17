package threading

// Executor, the instance can execute work at background
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

// Launcher, boot the Executor
type Launcher interface {
	// boot
	Boot()
	// shutdown, pool forbidden execute new works, but execute all of works in queue
	Shutdown()
	// showdown right away. pool forbidden execute new works, and all of works in queue
	ShutdownNow()
	// shutdown any of thread
	ShutdownAny()
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
