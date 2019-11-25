package models

import (
	"fmt"
)

// IP struct
type IP struct {
	ID    uint   `gorm:"primary_key" json:"-"`
	Data  string `gorm:"not null" json:"ip"`
	Type1 string `gorm:"not null" json:"type1"`
	Type2 string `gorm:"null" json:"type2,omitempty"`
	Speed int64  `gorm:"not null" json:"speed,omitempty"`
}

func (i *IP) String() string  {
	return fmt.Sprintf("Data:%s Type1:%s Type2:%s Speed:%d",i.Data,i.Type1,i.Type2,i.Speed)
}

// NewIP .
func NewIP() *IP {
	//init the speed to 100 Sec
	return &IP{Speed: 100}
}
