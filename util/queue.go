package util

import "sync"

type Task func() error

type Queue struct {
	lock sync.Locker
	tasks []Task
}
func NewQueue()*Queue{
	return &Queue{
		lock: NewSpinlock(),
	}
}

func (q *Queue)ForEach()error{
	q.lock.Lock()
	tasks := q.tasks
	q.lock.Unlock()
	for _,task := range tasks{
		if err:=task();err!=nil{
			return err
		}
	}
	return nil
}

func (q *Queue)Push(task Task)(num int){
	q.lock.Lock()
	q.tasks=append(q.tasks,task)
	num = len(q.tasks)
	q.lock.Unlock()
	return
}