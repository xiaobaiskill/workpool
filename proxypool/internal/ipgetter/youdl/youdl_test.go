package youdl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYoudl_Execute(t *testing.T) {
	y := New()
	ips := y.Execute()
	assert.Equal(t,true,len(ips)>0)
}
