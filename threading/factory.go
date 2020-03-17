package threading

import "time"

func CreateNewThreadPool(core, ext int, span time.Duration, w int, strategy func(interface{})) ThreadPool {
	tp := threadPoolImpl{}
	tp.Init(core, ext, span, w, strategy)
	tp.Boot()
	return &tp
}

func CreateNewStealPool(core, w int, strategy func(interface{})) StealPool {
	tp := stealPoolImpl{}
	tp.Init(core, w, strategy)
	tp.Boot()
	return &tp
}
