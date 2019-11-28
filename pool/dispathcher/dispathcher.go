package dispathcher

import (
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pkg/models"
	"github.com/xiaobaiskill/workpool/pool/queue"
	"github.com/xiaobaiskill/workpool/pool/worker"
	"sync"
)

var (
	lock sync.Mutex
	done = make(chan bool)  // 用于是否获取代理ip的chan
)

func StartDispathcher(nworks int) {
	// 初始化 并 获取代理ip
	queue.InitQueue()
	HttpClientMapAdd(conf.Conf.Pool.ProxyIpSize)

	// 开启工作
	worker.WorkerPool = make(worker.WorkerPoolType, nworks)
	for i:=1;i<nworks;i++{
		worker.NewWorker(i, worker.WorkerPool).Start()
	}

	// 用于获取ip, 获取方式：1 定时30秒获取一次 ， 2 当工作量 > 代理ip量时 获取一次
	go func() {
		//t := time.NewTicker(30 * time.Second)
		//defer t.Stop()
		for {
			select {
			/*case <-t.C:
				if queue.HCMap.Len() > 2 * conf.Conf.Pool.ProxyIpSize {
					continue
				}
				log.Logger.Info("定时获取一次 代理ip")
				HttpClientMapAdd(conf.Conf.Pool.ProxyIpSize)*/
			case <-queue.GetIpChan:
				log.Logger.Info("代理ip 数量减少，获取代理ip")
				HttpClientMapAdd(conf.Conf.Pool.ProxyIpSize)
			case <-done:
				return
			}

		}
	}()

	// 用于接收任务 并将任务 分发给 工作池中的一个work
	go func() {
		for {
			select {
			case work := <-worker.WorkQueue:
				go func() {
					workerPool := <-worker.WorkerPool
					workerPool <- work
				}()
			}
		}
	}()
	log.Logger.Info("dispathcher 程序启动完成")
}

func StopDispathcher() {
	done <- true  // 接收获取ip的程序
	for _, work := range worker.Works {
		work.Stop()
	}
}

// 用于获取 num 个代理ip， 添加至代理ip池中 和 代理ip 管道中
func HttpClientMapAdd(num int) {
	lock.Lock()
	ips, err := models.Conn.GetNumIPWithType(num)
	if err != nil || len(ips) == 0 {
		log.Logger.Error("没有获取到可用的ip")
		panic("没有获取可用的IP ，程序终止")
	}

	for _, ip := range ips {
		queue.HCMap.Add(ip)
	}
	lock.Unlock()
}
