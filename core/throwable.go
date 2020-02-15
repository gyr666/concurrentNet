package core

type Throwable uint8

var errorTable = []string{"NO_ACK","PEER RESET","KEEP COUNT IS MAX","PEER CLOSED"}

func (t *Throwable) Reason() string {
	if *t > 16 {
	}
	return errorTable[*t]
}

func (t *Throwable) isUserDefine() bool{
	return *t > 16
}