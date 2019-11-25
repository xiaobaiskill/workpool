package goubanjia

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGouBanJia_Execute(t *testing.T) {
	g := New()
	ips := g.Execute()
	assert.Equal(t,true,len(ips)>0)
}
