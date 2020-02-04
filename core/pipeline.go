// pipeline.go
package core

type Pipeline interface{
	addLast(func(Data) Data)
	addFirst(func(Data) Data)
}
