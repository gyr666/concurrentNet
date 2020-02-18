package threading

func CreateNewThreadPool() ThreadPool {
	return &threadPoolImpl{}
}
