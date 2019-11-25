package models

import (
	"github.com/ruoklive/proxypool/pkg/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisDB_InsertIP(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	err := d.InsertIP(&IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	})
	assert.NoError(t,err)
}

func TestRedisDB_CountIP(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	err := d.InsertIP(&IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	})
	assert.NoError(t,err)
	count,err := d.CountIP()
	assert.NoError(t,err)
	assert.Equal(t,true,count==1)
}

func TestRedisDB_DeleteIP(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	ip := &IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	}
	err := d.InsertIP(ip)
	assert.NoError(t,err)
	count,err := d.CountIP()
	assert.Equal(t,true,count==1)

	err = d.DeleteIP(ip)
	assert.NoError(t,err)

	count,err = d.CountIP()
	assert.NoError(t,err)
	assert.Equal(t,true,count==0)
}

func TestRedisDB_GetAllIP(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	ip := &IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	}
	err := d.InsertIP(ip)
	assert.NoError(t,err)

	resultIP,err := d.GetAllIP()
	assert.NoError(t,err)
	assert.Equal(t,true,len(resultIP)==1)
}

func TestRedisDB_FindIPWithType(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	ip := &IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	}
	err := d.InsertIP(ip)
	assert.NoError(t,err)
	count,err := d.CountIP()
	assert.Equal(t,true,count==1)

	resultIP,err := d.FindIPWithType("http")
	assert.NoError(t,err)
	assert.Equal(t,true,len(resultIP)==1)
}

func TestRedisDB_GetOneIP(t *testing.T) {
	d := NewRedisDB(redis.New("127.0.0.1:6379"))

	ip := &IP{
		Data: "10.10.10.0",
		Type1: "http",
		Speed: 1000,
	}
	err := d.InsertIP(ip)
	assert.NoError(t,err)
	count,err := d.CountIP()
	assert.Equal(t,true,count==1)

	resultIP,err := d.GetOneIP("10.10.10.0")
	assert.NoError(t,err)

	assert.Equal(t,"http",resultIP.Type1)
	assert.Equal(t,"10.10.10.0",resultIP.Data)
	assert.Equal(t,int64(1000),resultIP.Speed)
}