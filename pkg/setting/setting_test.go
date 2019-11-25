package setting

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContext(t *testing.T) {
	ini,err := NewContext("../../conf/app.ini")
	assert.NoError(t,err)
	assert.Equal(t,true,len(ini.GetString("","APP_NAME",""))> 0)
}
