package core

type TimeTrigger interface {
	Interval() uint
	Next() bool
}
