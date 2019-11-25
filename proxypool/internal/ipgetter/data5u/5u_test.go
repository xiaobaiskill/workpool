package data5u

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestData5u_Execute(t *testing.T) {
	u5 := New()
	ips := u5.Execute()
	assert.Equal(t,true,len(ips)>0)
}
