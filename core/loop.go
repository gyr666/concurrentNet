package core

type Loop interface {
	start()
	close()
	stop()
	eventHandler()
}
