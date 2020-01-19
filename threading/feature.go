package threading

type Future interface {
	get() interface{}
}
