workpool
====

### 关于proxypool
```
请将proxypool 的程序放入 proxypool 的目录下

目前使用的是汤涛的 proxypool 程序 
``` 

### 快速启动
`docker-compose up -d`

### api
 ```
 1、 /v1/catproxyip     查看当前 代理ip的 个数
 2、 /v1/index          代理请求
 3、 /metrics           prometheus 监控   
```

### 访问方式
```
curl -X POST http://127.0.0.1:8080/v1/index?proxy=true -d '{"action":"https://jimqaweb.mlytics.ai/cache.txt"}'
```

