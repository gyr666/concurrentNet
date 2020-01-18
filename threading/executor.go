package threading

type Executor interface {
	Exec(func())
}

type PoolExecutor interface {
	Executor
	Execwr(func() interface{})
	EXecwpr(func(...interface{}) interface{}, ...interface{})
	Execwp(func(...interface{}))
}
type Booter interface {
	Shutdown()
	ShutdownNow()
}

type Task struct {
	param []interface{}
	rev   interface{}
	t     ExecTask
}

type ExecTask func(...interface{}) interface{}
