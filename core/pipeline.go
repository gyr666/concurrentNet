// pipeline.go
package core

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/coder"
)

type Pipeline interface {
	AddEnCoder(coder.Encode)
	AddDeCoder(coder.Decode)
	AddLast(func(Data) Data)
	AddFirst(func(Data) Data)
	AddExceptionHandler(func(Throwable, Channel))
	doPipeline(buffer buffer.ByteBuffer) (buffer.ByteBuffer, error)
	doException(throwable Throwable, channel Channel)
}

type pipelineImpl struct {
	encoder coder.Encode
	decoder coder.Decode
	pipe    []func(Data)
	t       func(Throwable, Channel)
}

func (s *pipelineImpl) addEnCoder(e coder.Encode) {
	s.encoder = e
}

func (s *pipelineImpl) addDeCoder(d coder.Decode) {
	s.decoder = d
}

func (s *pipelineImpl) AddLast(f func(Data)) {
	s.pipe = append(s.pipe, f)
}

func (s *pipelineImpl) AddFirst(f func(Data)) {
	s.pipe = append(s.pipe, f)
}

func (s *pipelineImpl) doPipeline(inBuffer buffer.ByteBuffer, outBuffer buffer.ByteBuffer) error {
	//TODO
	//for{
	//	tran,err:= s.decoder(inBuffer)
	//	if err!=nil{
	//		return err
	//	}
	//	d := &dataImpl{data: tran}
	//	for i := range s.pipe {
	//		s.pipe[i](d)
	//	}
	//	s.encoder
	//}
	//
	//return buffer, s.encoder(d.data, buffer)
	return nil
}

func (s *pipelineImpl) AddExceptionHandler(f func(Throwable, Channel)) {
	s.t = f
}

func (s *pipelineImpl) doException(throwable Throwable, channel Channel) {
	s.t(throwable, channel)
}
