workpool
====

### 关于如何使用
```
请结合 proxypool 一起使用，

// 设置gin 的全局变量
export GIN_MODU=release

// 编译执行该程序（该程序依赖 redis 服务。请先启动proxypool 中的redis服务）
```


### 访问方式
```
curl -X POST http://127.0.0.1:8080/v1/index?proxy=true -d '{"action":"https://jimqaweb.mlytics.ai/cache.txt"}'
```

