package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	r := New("127.0.0.1:6379")
	_,err := r.Ping()
	assert.NoError(t,err)
}

func TestConn_Set(t *testing.T) {
	r := New("127.0.0.1:6379")
	err := r.Set("name","jmz")
	assert.NoError(t,err)
	s,err := r.GetString("name")
	assert.NoError(t,err)
	assert.Equal(t,"jmz",s)
}


