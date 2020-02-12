package core

type Data interface{
	Transfer(interface{})
	To() interface{}
}

type dataImpl struct{
}

func (d *dataImpl)Transfer(interface{}){
}

func (d *dataImpl)To() interface{}{
	return nil
}
