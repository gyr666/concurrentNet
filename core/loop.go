package core

type Loop interface {
	start()
	stop()
	addChannel(Channel)
}
