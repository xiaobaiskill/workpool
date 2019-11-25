package ip181

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIp181_Execute(t *testing.T) {
	p := New()
	ips := p.Execute()
	assert.Equal(t,true,len(ips)>0)
}