package utils

import (
	"fmt"
)

type Task struct {
	f func() error // 一个无参的函数类型
}

// 通过NewTask来创建一个Task
func NewTask(f func() error) *Task {
	t := Task{
		f: f,
	}
	return &t
}

// 执行Task任务的方法
func (t *Task) Execute() {
	t.f()
}

type Pool struct {
	// 对外接收Task的入口
	EntryChannel chan *Task

	// 携程吃最大worker数量，限定Goroutine的个数
	WorkerNum int

	//协程池内部的任务就绪队列
	JobsChannel chan *Task
}

//创建一个携程池
func NewPool(cap int) *Pool {
	p := Pool{
		EntryChannel: make(chan *Task),
		WorkerNum:    cap,
		JobsChannel:  make(chan *Task),
	}

	return &p
}

// 协程池创建一个worker并且开始工作
func (p *Pool) worker(workId int) {
	// work不断的从JobsChannel内部任务队列中那任务
	for task := range p.JobsChannel {
		// 如果拿到任务，则执行task任务
		task.Execute()
		fmt.Println("worker id ", workId, " 执行完毕任务")
	}

}

// 协程池Pool开始工作
func (p *Pool) Run() {
	// 1.首先根据协程池的worker数量限定，开启固定数量的worker,每一个worker用一个Goroutine承载
	for i := 0; i < p.WorkerNum; i++ {
		go p.worker(i)
	}

	// 2.从EntryChannel协程池入口取外界传递过来的任务，并且将任务送进JobsChannel
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	// 3.执行完毕需要关闭JobsChannel
	close(p.JobsChannel)

	// 4.执行完毕需要关闭EntryChannel
	close(p.EntryChannel)
}

//func test(test *testing.T) {
//
//	// 创建一个Task
//	t := NewTask(func() error {
//		fmt.Println(time.Now())
//		return nil
//	})
//
//	// 创建一个协程池，最大开启3个协程worker
//	p := NewPool(3)
//
//	// 开一个协程 不断的向pool输送打印一条时间的task任务
//	go func() {
//		for {
//			p.EntryChannel <- t
//		}
//	}()
//
//	p.Run()
//
//}
