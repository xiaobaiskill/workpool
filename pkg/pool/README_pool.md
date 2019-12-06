使用 workpool
===

#### 一、如何使用
```
// 1、开启workpool
StartDispathcher(worknum int)gin.HandlerFunc
// 参数： worknum    启动多少个work
// 返回： 用于监控workpool 的job 使用情况。

// 2、关闭workpool
StopDispathcher()


// 3、如何传入job
WorkQueue<-job


// 4、job 实现
job 是一个接口 ，只需要实现job 的方法即可
```