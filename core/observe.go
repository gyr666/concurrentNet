package core

import "fmt"

type ServerObserve interface {
	OnBooting()
	OnBooted(l []NetworkInet64)
	OnStopping()
	OnStopped()
}

type DefaultObserve struct {
}

func (d *DefaultObserve) OnBooting() {
	fmt.Println("On Booting")
}

func (d *DefaultObserve) OnBooted(l []NetworkInet64) {
	fmt.Printf("On Booted \nlisten:%#v\n", l)
}

func (d *DefaultObserve) OnStopping() {
	fmt.Println("On Stopping")
}

func (d *DefaultObserve) OnStopped() {
	fmt.Println("On Stopped")
}
