package worker

import (
	"github.com/xiaobaiskill/workpool/pool/job"
)

var (
	WorkQueue  = make(chan job.Job, 100)
	WorkerPool WorkerPoolType
	Works = make([]*Worker,0)
)


type WorkerPoolType chan chan job.Job

// 工作结构体
type Worker struct {
	ID       int
	Work     chan job.Job      // 工作管道
	WorkPool chan chan job.Job // 全局 工作管道 池
	End      chan bool
}

func (w Worker) Start() {
	go func() {
		for {
			// 将工作加入至工作池全局队列中
			w.WorkPool <- w.Work
			select {
			// 从工作管道中获取任务
			case work := <-w.Work:
				work.Execute()
			case <-w.End:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.End <- true
	}()
}

/*
	// id
    // IP 生成代理httpclient
    // workerQueue 全局 工作管道 队列
	// timeOut http 请求超时时间
*/
func NewWorker(id int, workerQueue chan chan job.Job) *Worker {
	work := &Worker{
		ID:       id,
		Work:     make(chan job.Job),
		WorkPool: workerQueue,
		End:      make(chan bool),
	}

	Works = append(Works,work)

	return work

}

