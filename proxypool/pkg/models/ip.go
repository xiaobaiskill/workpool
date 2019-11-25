package models

import (
	"fmt"
)

// IP struct
type IP struct {
	ID    int64  `xorm:"pk autoincr" json:"-"`
	Data  string `xorm:"NOT NULL" json:"ip"`
	Type1 string `xorm:"NOT NULL" json:"type1"`
	Type2 string `xorm:"NULL" json:"type2,omitempty"`
	Speed int64  `xorm:"NOT NULL" json:"speed,omitempty"`
}

func (i *IP) String() string  {
	return fmt.Sprintf("Data:%s Type1:%s Type2:%s Speed:%d",i.Data,i.Type1,i.Type2,i.Speed)
}

// NewIP .
func NewIP() *IP {
	//init the speed to 100 Sec
	return &IP{Speed: 100}
}
