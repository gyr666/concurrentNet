package test

import (
	"fmt"
	"testing"

	"gunplan.top/concurrentNet/util"
)

func TestCounter(t *testing.T) {
	c := util.NewCounter()
	c.Boot()
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)

	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(88)
	c.Push(120000)
	c.Push(27)
	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)
	c.Push(99)
	c.Push(88)

	c.Push(88)
	c.Push(120000)
	c.Push(27)

	c.Push(99)
	c.Push(99)
	c.Push(88)
	c.Push(99)

	c.Push(99)

	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())

	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
	fmt.Println(c.Sum(), c.Max(), c.Min())
}
