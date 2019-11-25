package register

import "github.com/ruoklive/proxypool/pkg/models"

type IPGetter interface {
	// Execute 执行获取IP操作
	Execute() []*models.IP
}

type Executor func() IPGetter

var executors = make([]Executor, 0)
// Add 添加IP getter
func Add(executor Executor) {
	executors = append(executors, executor)
}

func GetExecutors()  []Executor {
	return executors
}