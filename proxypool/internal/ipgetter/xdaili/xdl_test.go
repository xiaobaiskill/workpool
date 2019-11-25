package xdaili

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXdaili_Execute(t *testing.T) {
	x := New()
	ips := x.Execute()
	assert.Equal(t,true,len(ips)>0)
}
