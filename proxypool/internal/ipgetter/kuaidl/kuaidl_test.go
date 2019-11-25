package kuaidl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKuaidl_Execute(t *testing.T) {
	k := New()
	ips := k.Execute()
	assert.Equal(t,true,len(ips)>0)
}
