package config

import (
	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/util"
	"strings"
)

func StandDecoder(buffer buffer.ByteBuffer, config *Config) error {
	res, ok := buffer.Read(buffer.Size())
	if ok != nil {
		return ok
	}

	for _, v := range strings.Split(string(res), "\n") {
		name := util.GetFieldsFromTag(config, "stand", strings.TrimSpace(strings.Split(v, "=")[0]))
		util.GetFieldsFromNameAndSet(config, name, util.GetFieldsFromTag(config, "stand", strings.TrimSpace(strings.Split(v, "=")[0])))
	}
	return nil
}
