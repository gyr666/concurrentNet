package util

import (
	"math"
)

type DynamicCounter interface {
	init()
	Push(uint64)
	Max() uint64
	Min() uint64
	Ave() uint64
	Size() uint64
	Mid() uint64
	Sum() uint64
	Boot()
}

func NewCounter() DynamicCounter {
	return &counterImpl{}
}

type counterImpl struct {
	size    uint64
	ave0    uint64
	sum0    uint64
	max0    uint64
	min0    uint64
	channel chan uint64
}

func (c *counterImpl) init() {
	c.channel = make(chan uint64, 100)
	c.min0 = math.MaxInt64
	c.max0 = 0
}
func (c *counterImpl) Push(u uint64) {
	c.size++
	c.channel <- u
}
func (c *counterImpl) Max() uint64 {
	return c.max0
}

func (c *counterImpl) Size() uint64 {
	return c.size
}
func (c *counterImpl) Min() uint64 {
	return c.min0
}
func (c *counterImpl) Ave() uint64 {
	return c.ave0
}
func (c *counterImpl) Mid() uint64 {
	return c.ave0
}
func (c *counterImpl) Sum() uint64 {
	return c.sum0
}

func (c *counterImpl) Boot() {
	c.init()
	go func() {
		for {
			select {
			case data := <-c.channel:
				c.sum0 += data
				c.ave0 = c.sum0 / c.size
				if c.max0 < data {
					c.max0 = data
				}
				if c.min0 > data {
					c.min0 = data
				}
			}
		}
	}()
}
