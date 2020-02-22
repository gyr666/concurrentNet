package core

type Loop interface {
	start()
	stop()
	eventHandler()
}


