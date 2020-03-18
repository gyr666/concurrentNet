package threading

import "time"

func CreateNewThreadPool(core, ext int, span time.Duration, w int, strategy func(interface{})) ThreadPool {
	return new(threadPoolImpl).Init(core, ext, span, w, strategy).Boot().(ThreadPool)
}

func CreateNewStealPool(core, w int, strategy func(interface{})) ThreadPool {
	return new(stealPoolImpl).Init(core, w, strategy).Boot().(ThreadPool)
}
