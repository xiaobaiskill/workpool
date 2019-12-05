package pool

import "github.com/gin-gonic/gin"

//type DispathCher struct{
//	Metrics *metrics
//}

func StartDispathcher(nworks int) gin.HandlerFunc{
	// 初始化 prometheus
	initMetrics("workpool")

	// 开启工作
	workerPool = make(workerPoolType, nworks)
	for i:=1;i<nworks;i++{
		newWorker(i, workerPool).start()
	}

	// 用于接收任务 并将任务 分发给 工作池中的一个work
	go func() {
		for {
			select {
			case work := <-WorkQueue:
				workerPool := <-workerPool
				workerPool <- work
			}
		}
	}()

	return metric.ginHandler()
}

func StopDispathcher() {
	for _, work := range works {
		work.stop()
	}
}
