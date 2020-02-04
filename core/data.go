package core

type Data interface{
	Transfer(interface{})
	To() interface{}
}
