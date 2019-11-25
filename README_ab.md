* 请求：

`ab -c 10 -n 100 -p post.txt http://127.0.0.1:8080/v1/index?proxy=true`

* 结果
```
Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8080

Document Path:          /v1/index?proxy=true
Document Length:        24 bytes

Concurrency Level:      10
Time taken for tests:   28.853 seconds
Complete requests:      100
Failed requests:        40
   (Connect: 0, Receive: 0, Length: 40, Exceptions: 0)
Total transferred:      15460 bytes
Total body sent:        20000
HTML transferred:       3160 bytes
Requests per second:    3.47 [#/sec] (mean)
Time per request:       2885.327 [ms] (mean)
Time per request:       288.533 [ms] (mean, across all concurrent requests)
Transfer rate:          0.52 [Kbytes/sec] received
                        0.68 kb/s sent
                        1.20 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0       2
Processing:     1 2307 1761.8   3629    4008
Waiting:        1 2307 1761.8   3629    4008
Total:          1 2307 1761.7   3629    4008

Percentage of the requests served within a certain time (ms)
  50%   3629
  66%   4004
  75%   4005
  80%   4006
  90%   4007
  95%   4007
  98%   4008
  99%   4008
 100%   4008 (longest request)
```