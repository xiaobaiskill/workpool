package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultDB(t *testing.T) {
	x,err := gorm.Open("postgres","host=127.0.0.1 port=5432 user=postgres dbname=postgres sslmode=disable password=example")
	assert.NoError(t,err)
	x.AutoMigrate(NewIP())
	db  := NewDefaultDB(x)
	db.InsertIP(&IP{21,"127.0.0.1","http","",10})
	db.InsertIP(&IP{32,"127.0.0.1","http","",20})

	ips,err := db.GetNumIPWithType(50)
	assert.NoError(t,err)
	assert.Equal(t,true,len(ips)>0)
	count,err := db.CountIP()
	assert.NoError(t,err)
	assert.Equal(t,true,count>0)

}
