package pool


var (
	WorkQueue  = make(chan Job, 100)
	workerPool workerPoolType
	works = make([]*worker,0)
)


type workerPoolType chan chan Job

// 工作结构体
type worker struct {
	id       int
	work     chan Job      // 工作管道
	workPool chan chan Job // 全局 工作管道 池
	end      chan bool
}

func (w worker) start() {
	go func() {
		for {
			// 将工作加入至工作池全局队列中
			w.workPool <- w.work
			select {
			// 从工作管道中获取任务
			case work := <-w.work:
				metric.jobTotalInc()
				metric.jobexectingInc()
				Execute()
				metric.jobexectingDec()
			case <-w.end:
				return
			}
		}
	}()
}

func (w worker) stop() {
	go func() {
		w.end <- true
	}()
}

/*
	// id
    // IP 生成代理httpclient
    // workerQueue 全局 工作管道 队列
	// timeOut http 请求超时时间
*/
func newWorker(id int, workerQueue chan chan Job) *worker {
	work := &worker{
		id:       id,
		work:     make(chan Job),
		workPool: workerQueue,
		end:      make(chan bool),
	}

	works = append(works,work)

	return work

}
