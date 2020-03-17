package threading

type PoolStatusTransfer struct {
	s PoolState
}

func (p *PoolStatusTransfer) init() {
	p.s = STOPPED
}

func (p *PoolStatusTransfer) whenThreadBooting() {
	if p.s != STOPPED {
		panic("status is error")
	}
	p.s = BOOTING
}

func (p *PoolStatusTransfer) whenThreadBooted() {
	if p.s != BOOTING {
		panic("status is error")
	}
	p.s = RUNNING
}

func (p *PoolStatusTransfer) whenThreadStopping() {
	if p.s != RUNNING {
		panic("status is error")
	}
	p.s = STOPPING
}

func (p *PoolStatusTransfer) whenThreadStopped() {
	if p.s != STOPPING {
		panic("status is error")
	}
	p.s = STOPPED
}

func (p *PoolStatusTransfer) get() PoolState {
	return p.s
}
