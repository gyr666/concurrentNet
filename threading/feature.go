package threading

type Future interface {
	isDone() bool
	get() interface{}
}
