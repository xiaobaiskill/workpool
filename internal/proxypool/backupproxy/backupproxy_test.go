package backupproxy

import (
	"github.com/stretchr/testify/assert"
	"github.com/xiaobaiskill/workpool/internal/proxypool"
	"net/http"
	"testing"
)

func TestNewBackupProxy(t *testing.T) {
	//assert.NoError()
	backup := NewBackupProxy("../../../conf/backupproxy.test.conf", 10)
	backup.AddMetric(proxypool.NewMetrics("aaa"))

	assert.NotNil(t, backup.httpClients["http://127.0.0.1:5839"].Client)
	assert.NotNil(t, backup.httpClients["https://39.193.12.7:8889"].Client)
	assert.NotNil(t, backup.httpClients["https://38.23.214.22:7912"].Client)
	assert.NotNil(t, backup.httpClients["http://23.32.231.4:19283"].Client)
	assert.NotNil(t, backup.httpClients["https://192.32.32.23:11113"].Client)
	backup.Del("http://127.0.0.1:5839")
	backup.Del("https://39.193.12.7:8889")
	backup.Del("https://38.23.214.22:7912")
	backup.Del("http://23.32.231.4:19283")
	backup.Del("https://192.32.32.23:11113")
	assert.Equal(t,0,len(backup.httpClients))

	backup.Pop()
	backup.Pop()
	backup.Pop()
	backup.Pop()
	backup.Pop()

	backup.Push(proxypool.HTTPClientMap{"http://127.0.0.1:1111",&http.Client{}})
	client,ok := backup.Pop()
	assert.Equal(t,true,ok)
	assert.Equal(t,"http://127.0.0.1:1111",client.Ip)
}
