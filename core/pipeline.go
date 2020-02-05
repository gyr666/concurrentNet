// pipeline.go
package core

type Pipeline interface{
	AddLast(func(Data) Data)
	AddFirst(func(Data) Data)
}
