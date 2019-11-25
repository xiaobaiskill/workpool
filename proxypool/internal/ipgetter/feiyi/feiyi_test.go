package feiyi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFeiyi_Execute(t *testing.T) {
	f := New()
	ips := f.Execute()
	assert.Equal(t,true,len(ips)>0)
}
