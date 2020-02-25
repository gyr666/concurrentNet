package config

import (
	"strconv"
	"strings"
	"sync"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/util"
)

const SpiltChar = ":"
const LineSpiltChar = "\n"
const DecoderTagName = "line"

func LineDecoder(buffer buffer.ByteBuffer, config *Config) error {
	res := buffer.FastMoveOut()
	w := sync.WaitGroup{}
	for _, v := range strings.Split(string(res), LineSpiltChar) {
		w.Add(-1)
		go parallel(v, config, &w)
	}
	w.Wait()
	return nil
}

func parallel(list string, config *Config, w *sync.WaitGroup) {
	l := strings.Split(list, SpiltChar)
	k := strings.TrimSpace(l[0])
	v := strings.TrimSpace(l[1])
	name := util.GetFieldFromTag(config, DecoderTagName, k)
	if util.GetFieldTag(config, name, TypeName) == Number {
		v, _ := strconv.Atoi(v)
		util.GetFieldsFromNameAndSet(config, name, v)
	} else if util.GetFieldTag(config, name, TypeName) == String {
		util.GetFieldsFromNameAndSet(config, name, v)
	}
	w.Done()
}
