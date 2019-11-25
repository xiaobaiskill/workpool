package plp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPLP_Execute(t *testing.T) {
	p := New()
	ips := p.Execute()
	assert.Equal(t,true,len(ips)>0)
}
