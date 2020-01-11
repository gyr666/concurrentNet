package event_loop

import (
	"golang.org/x/sys/unix"
)

type AbstractEventLoop struct {
	port int
}

func init0(port int) ExecCode {
	err := initKqueue()
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	if err != nil && !sendToErrorCenter(err) {
		return ExecCode(ERROR)
	}
	err = unix.Bind(socket, &unix.SockaddrInet4{Port: port})
	if err != nil && !sendToErrorCenter(err) {
		return ExecCode(ERROR)
	}
	err = unix.Listen(socket, 2058)
	if err != nil && !sendToErrorCenter(err) {
		return ExecCode(ERROR)
	}
	return ExecCode(ASYNC | NORMAL)
}

func initKqueue() error {
	fd, err := unix.Kqueue()
	events := []unix.Kevent_t{}
	timeout := unix.Timespec{}
	if err != nil && !sendToErrorCenter(err) {

		//	return ExecCode(ERROR)
	}
	//todo tomorrow
	unix.Kevent(fd, events, events, &timeout)
	return err
}

func sendToErrorCenter(e error) bool {
	println("ERROR:", e)
	return false
}
func NewBaseEventLoop() BaseEventLoop {
	return &AbstractEventLoop{port: 8679}
}

func (a *AbstractEventLoop) StartLoop() ExecCode {
	code := init0(a.port)
	return code
}
func (a *AbstractEventLoop) haltLoop() ExecCode {
	return 0
}
func (a *AbstractEventLoop) resumeLoop() ExecCode {
	return 0
}
func (a *AbstractEventLoop) stopLoop() ExecCode {
	return 0
}
