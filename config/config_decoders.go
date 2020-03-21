package config

import (
	"strconv"
	"strings"
	"sync"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/util"
)

const SpiltChar = ": "
const LineSpiltChar = "\n"
const DecoderTagName = "line"

func LineDecoder(buffer buffer.ByteBuffer, config *Config) error {
	res := buffer.FastMoveOut()
	w := sync.WaitGroup{}
	for _, v := range strings.Split(string(*res), LineSpiltChar) {
		w.Add(1)
		go parallel(v, config, &w)
	}
	w.Wait()
	return nil
}

func parallel(list string, config *Config, w *sync.WaitGroup) {
	defer w.Done()
	if strings.HasPrefix(list, "#") {
		return
	}
	l := strings.Split(list, SpiltChar)
	k := strings.TrimSpace(l[0])
	v := strings.TrimSpace(l[1])
	name := util.GetFieldFromTag(config, DecoderTagName, k)
	tag := util.GetFieldTag(config, name, TypeName)
	if tag == Number {
		v, _ := strconv.Atoi(v)
		util.GetFieldsFromNameAndSetInt(config, name, v)
	} else if tag == String {
		util.GetFieldsFromNameAndSetString(config, name, v)
	} else if tag == Map {
		util.InvokeMapMethod(config, "Set"+name, v)
	}
}
