package coder

import "gunplan.top/concurrentNet/buffer"
type StringTransfer struct {

}
func (s *StringTransfer) encode(i interface{},b *buffer.ByteBuffer) error  {
	(*b).Write([]byte(i.(string)))
	return nil
}

func (s *StringTransfer) decode(b *buffer.ByteBuffer) interface{} {
	data,_ := (*b).Read(-1)
	return string(data)
}