package backupproxy

import (
	"bufio"
	"fmt"
	. "github.com/xiaobaiskill/workpool/internal/proxypool"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

type backupproxy struct {
	maxSize          int    // 设置最多放多少个备用ip
	backupproxyfile  string // 备用代理池ip 的文件
	httpClients      map[string]HTTPClientMap
	httpClientsQueue chan HTTPClientMap
	sync.Mutex
}

func (b *backupproxy) init() {
	b.httpClients = make(map[string]HTTPClientMap)
	b.httpClientsQueue = make(chan HTTPClientMap, b.maxSize)

	b.getFileipAdd()

	go func(fileTime int64) {
		tc := time.NewTicker(2 * time.Minute) // 5秒检测一次文件是否修改
		defer tc.Stop()
		for {
			select {
			case <-tc.C:
				if fileNewTime := b.getFileModTime(); fileNewTime != fileTime {
					fileTime = fileNewTime // 更新文件修改时间
					b.getFileipAdd()       // 添加新的备用代理
				}
			}
		}
	}(b.getFileModTime())

}

func (b *backupproxy) Pop() (client HTTPClientMap, ok bool) {
	b.Lock()
	defer b.Unlock()
	if len(b.httpClientsQueue) > 0{
		ok = true
		client = <-b.httpClientsQueue
	}
	return
}

func (b *backupproxy) Push(httpclientip HTTPClientMap) {
	b.httpClientsQueue <- httpclientip
}

func (b *backupproxy) Del(ip string) {
	b.Lock()
	delete(b.httpClients,ip)
	b.Unlock()
}
func (b *backupproxy) Len()int{
	return len(b.httpClients)
}

func (b *backupproxy) add(ip string) {
	b.Lock()
	defer b.Unlock()
	if len(b.httpClients) >= b.maxSize {
		return
	}

	if _, ok := b.httpClients[ip]; ok {
		return
	}

	hcip := HTTPClientMap{ip, b.createHttpClient(ip)}
	b.httpClients[ip] = hcip
	b.httpClientsQueue <- hcip
}

func (b *backupproxy) getFileipAdd() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	f, err := os.Open(b.backupproxyfile)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	bf := bufio.NewReader(f)
	reg := regexp.MustCompile(`.*?((http|https):\/\/\d{1,3}\.\d{1,3}.\d{1,3}\.\d{1,3}:\d+).*?`)

	for {
		s, err := bf.ReadString('\n')

		isbool := reg.MatchString(s)
		if isbool {
			ip := reg.FindStringSubmatch(s)[1]
			b.add(ip)
		}

		if err == io.EOF || err != nil {
			break
		}
	}
}

func (b *backupproxy) createHttpClient(ip string) *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(ip)
	}

	transport := &http.Transport{
		Proxy:               proxy,
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 5,
	}
	client := &http.Client{
		Transport: transport,
		Timeout: 4 * time.Second,
	}
	return client
}

func (b *backupproxy) getFileModTime() int64 {
	f, err := os.Open(b.backupproxyfile)
	if err != nil {
		return 0
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0
	}

	return fi.ModTime().Unix()
}

func NewBackupProxy(file string, maxsize int) *backupproxy {
	b := new(backupproxy)
	b.backupproxyfile = file
	b.maxSize = maxsize
	b.init()
	return b
}
