package core

type TimeTrigger interface {
	Interval() uint
	next() bool
}
