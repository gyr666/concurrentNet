package threading

type BasePool struct {
	// this is a pointer
	addQueue0 func(task *Task)
	status    PoolStatusTransfer
}

func (b *BasePool) Init(a func(*Task)) {
	b.status = PoolStatusTransfer{}
	b.addQueue0 = a
}

func (b *BasePool) Exec(f func()) {
	_Exec(f, b.addQueue)
}

func (b *BasePool) Execwr(f func() interface{}) Future {
	return _Execwr(f, b.addQueue)
}

func (b *BasePool) Execwp(f func(...interface{}), p ...interface{}) {
	_Execwp(f, b.addQueue, p...)
}

func (b *BasePool) Execwpr(f func(...interface{}) interface{}, p ...interface{}) Future {
	return _Execwpr(f, b.addQueue, p...)
}

func (b *BasePool) Status() PoolState {
	return b.status.get()
}

func (b *BasePool) addQueue(task *Task) {
	switch b.status.get() {
	case BOOTING:
		fallthrough
	case RUNNING:
		b.addQueue0(task)
	case STOPPED:
		fallthrough
	case STOPPING:
		panic("pool has been close")
	case WAITING:
		// spin add
		b.addQueue(task)
	}
}
