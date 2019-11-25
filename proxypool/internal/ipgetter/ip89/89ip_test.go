package ip89

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIp89_Execute(t *testing.T) {
	e := New()
	ips := e.Execute()
	assert.Equal(t,true,len(ips)>0)
}
