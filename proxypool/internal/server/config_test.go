package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig_Load(t *testing.T) {

	cfg := NewConfig()
	assert.Equal(t,time.Hour*168,cfg.Server.SessionExpires)
	assert.Equal(t,true,cfg.LogXorm.Rotate)
	assert.Equal(t,100,cfg.Log.BufferLen)

	cfg.Load("app_test.ini")

	assert.Equal(t,time.Hour*72,cfg.Server.SessionExpires)
	assert.Equal(t,false,cfg.LogXorm.Rotate)
	assert.Equal(t,200,cfg.Log.BufferLen)
}
