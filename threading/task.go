package threading

func _Exec(f func(), addQueue func(*Task)) {
	addQueue(&Task{param: nil, t: func(i ...interface{}) interface{} {
		f()
		return nil
	}})
}

func _Execwr(f func() interface{}, addQueue func(*Task)) Future {
	tsk := Task{param: nil, t: func(i ...interface{}) interface{} {
		return f()
	}}
	tsk.init()
	addQueue(&tsk)
	return newFuture(&tsk.rev)
}

func _Execwp(f func(...interface{}), addQueue func(*Task), p ...interface{}) {
	addQueue(&Task{param: p, t: func(i ...interface{}) interface{} {
		f(i...)
		return nil
	}})
}

func _Execwpr(f func(...interface{}) interface{}, addQueue func(*Task), p ...interface{}) Future {
	tsk := Task{param: p, t: f}
	tsk.init()
	addQueue(&tsk)
	return newFuture(&tsk.rev)
}
