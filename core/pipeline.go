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
	doException(throwable Throwable, channel ChildChannel)
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

func (s *pipelineImpl) doPipeline(buffer buffer.ByteBuffer, a buffer.Allocator) (buffer.ByteBuffer, error) {
	tran := s.decoder(buffer)
	d := &dataImpl{data: tran}
	for i := range s.pipe {
		s.pipe[i](d)
	}
	return buffer, s.encoder(d.data, buffer)
}

func (s *pipelineImpl) AddExceptionHandler(f func(Throwable, Channel)) {
	s.t = f
}

func (s *pipelineImpl) doException(throwable Throwable, channel ChildChannel) {
	s.t(throwable, channel)
}
