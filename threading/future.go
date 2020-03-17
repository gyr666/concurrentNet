package threading

type Future interface {
	isDone() bool
	Get() interface{}
}

type futureImpl struct {
	wait *chan interface{}
}

func (f *futureImpl) Get() interface{} {
	defer func() {
		close(*f.wait)
	}()
	return <-*f.wait
}

func (f *futureImpl) isDone() bool {
	return len(*f.wait) == 1
}

func newFuture(c *chan interface{}) Future {
	return &futureImpl{wait: c}
}
