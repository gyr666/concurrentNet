package threading

type Future interface {
	isDone() bool
	Get() interface{}
}
