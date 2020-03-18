package core

type listener struct {
}

func (ln *listener) close() {

}

func (ln *listener) accept() Channel {
	return nil
}
