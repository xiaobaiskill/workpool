# goProxyPool

爬取公開代理池，提供給爬蟲做IP偽裝的一個專案。

# 支持网站

[支持网站](./internal/ipgetter/README.md)

# 快速啟用
1. docker-compose up -d
2. 檢查端口`http://localhost:8080/v2/ip`
3. 监控数据`http://localhost:8080/v2/monitor`

# refer:
參考專案go代理池
- https://github.com/henson/proxypool


# 備註:
預設app.ini使用db為postgres帳密(postgres/example)
