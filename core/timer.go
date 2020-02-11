package core

type TimeTigger interface {
	Intervale() uint
	next() bool
}
