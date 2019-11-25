package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xiaobaiskill/workpool/pkg/redis"
	"strconv"
	"strings"
)

var Conn DB

type DB interface {
	// InsertIp 存储IP
	InsertIP(ip *IP) error
	// CountIPs 统计数据库里IP总数量
	CountIP() (int64, error)
	// DeleteIP 删除指定IP
	DeleteIP(ip *IP) error
	// GetOneIP 获取指定的一个IP
	GetOneIP(ip string) (*IP, error)
	// GetAllIP 获取所有IP
	GetAllIP() ([]*IP, error)
	// 获取指定数量的ip
	GetNumIPWithType(int)([]*IP,error)
	// FindIPWithType 通过类型获取IP
	FindIPWithType(typ string) ([]*IP, error)
	// UpdateToFirstIP 将指定IP更新到第一个位置的IP数据上
	//UpdateToFirstIP(ip *IP) error
	// 是否存在https的IP
	ExistHttps() (bool,error)
}

type DefaultDB struct {
	x *gorm.DB
}

func NewDefaultDB(x *gorm.DB) *DefaultDB {
	return &DefaultDB{
		x: x,
	}
}

func (d *DefaultDB) InsertIP(ip *IP) error {

	tx := d.x.Begin()
	if err := tx.Create(ip).Error;err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (d *DefaultDB) CountIP() (int64, error)  {
	var count int64
	if err := d.x.Model(NewIP()).Where("id >= ?",0).Count(&count).Error;err !=nil {
		return count,err
	}
	return count,nil
}

func (d *DefaultDB) DeleteIP(ip *IP) (err error ) {
	err = d.x.Where("id = ?",ip.ID).Error
	return
}

func (d *DefaultDB) GetOneIP(ip_data string) (ip *IP,err error)  {
	ip = NewIP()
	if err := d.x.Where("data=?",ip_data).Find(ip).Error; err !=nil{
		return nil,err
	}
	return
}

func (d *DefaultDB) GetAllIP() ([]*IP, error)  {


	tmpIp := make([]*IP, 0)

	err := d.x.Where("speed <= 1000").Find(&tmpIp).Error
	if err != nil {
		return nil, err
	}
	return tmpIp, nil
}

func (d *DefaultDB) GetNumIPWithType(num int)([]*IP,error){
	temIp := make([]*IP,num)
	err := d.x.Order("speed asc").Limit(num).Find(&temIp).Error
	return temIp,err
}

func (d *DefaultDB) FindIPWithType(typ string) ([]*IP, error)  {
	tmpIp := make([]*IP, 0)
	switch typ {
	case "http":
		err := d.x.Where("speed <= 1000 and type1=?", "http").Find(&tmpIp).Error
		if err != nil {
			return tmpIp, err
		}
	case "https":
		//test has https proxy on databases or not
		hasHttps,err := d.ExistHttps()
		if err!=nil {
			return nil,err
		}
		if hasHttps == false {
			return tmpIp, nil
		}
		err = d.x.Where("speed <= 1000 and type2=?", "https").Find(&tmpIp).Error
		if err != nil {
			return tmpIp, err
		}
	default:
		return tmpIp, nil
	}

	return tmpIp, nil
}

func (d *DefaultDB) UpdateToFirstIP(ip *IP) error  {
	err := d.x.Where("id=?",1).Update(ip).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DefaultDB) ExistHttps() (bool,error){
	ip := NewIP()
	err := d.x.Where("type2=?","https").Find(ip).Error
	if err != nil {
		return false,err
	}

	return true,nil
}

type RedisDB struct {
	conn *redis.Conn
	ipsKey string
	ipPrefixKey string
}

func NewRedisDB(conn *redis.Conn) *RedisDB  {
	return &RedisDB{
		conn:conn,
		ipsKey: "ips",
		ipPrefixKey: "ip:",
	}
}

func (r *RedisDB) InsertIP(ip *IP) error {
	err := r.conn.Hmset(r.ipPrefixKey + ip.Data,"type1",ip.Type1,"type2",ip.Type2,"speed",fmt.Sprintf("%d",ip.Speed))
	if err!=nil {
		return err
	}
	err = r.conn.ZAdd(r.ipsKey,float64(ip.Speed),ip.Data)
	if err!=nil {
		return err
	}
	return nil
}

func (r *RedisDB) CountIP() (int64, error) {
	return r.conn.ZCard(r.ipsKey)
}

func (r *RedisDB) DeleteIP(ip *IP) error  {
	err := r.conn.Del(r.ipPrefixKey+ip.Data)
	if err!=nil {
		return err
	}
	_,err = r.conn.ZRem(r.ipsKey,ip.Data)
	if err!=nil {
		return err
	}
	return nil
}

func (r *RedisDB) GetOneIP(ip string) (*IP, error)   {
	fields,err := r.conn.Hgetall(r.ipPrefixKey+ip)
	if err!=nil {
		return nil,err
	}
	ipModel := NewIP()
	ipModel.Data = ip
	if fields!=nil && len(fields)>0 {
		for i:=0;i< len(fields);i++ {
			field := fields[i]
			switch field.Field {
			case "type1":
				ipModel.Type1 = field.Value
				break
			case "type2":
				ipModel.Type2 = field.Value
				break
			case "speed":
				ipModel.Speed,_ = strconv.ParseInt(field.Value,10,64)
				break


			}
		}
	}
	return ipModel,nil
}

func (r *RedisDB) GetAllIP() ([]*IP, error) {
	fieldValues,err := r.conn.ZRange(r.ipsKey,0,-1)
	if err!=nil {
		return nil,err
	}
	ipModels := make([]*IP,0)
	if len(fieldValues)>0 {
		for i:=0;i< len(fieldValues);i=i+2 {
			ipModel,err := r.GetOneIP(fieldValues[i])
			if err!=nil {
				continue
			}
			ipModels = append(ipModels,ipModel)
		}
	}
	return ipModels,nil
}

func (r *RedisDB) GetNumIPWithType(num int)([]*IP,error){
	fieldValues,err := r.conn.ZRange(r.ipsKey,0,int64(num))
	if err!=nil {
		return nil,err
	}
	ipModels := make([]*IP,0)
	if len(fieldValues) >0{
		for _,v := range fieldValues{
			ipModel,err := r.GetOneIP(v)
			if err!=nil {
				continue
			}
			ipModels = append(ipModels,ipModel)
		}
	}
	return ipModels,nil
}


func (r *RedisDB) FindIPWithType(typ string) ([]*IP, error) {
	ips,err := r.GetAllIP()
	if err!=nil {
		return nil,err
	}
	resultIps := make([]*IP,0)
	if len(ips)>0 {
		for _,ip :=range ips {
			if strings.ToLower(typ) == "http" && strings.ToLower(ip.Type1) !="http" {
				continue
			}
			if strings.ToLower(typ) == "https" && strings.ToLower(ip.Type2) !="https" {
				continue
			}
			resultIps = append(resultIps,ip)
		}
	}
	return resultIps,nil
}

func (r *RedisDB) ExistHttps() (bool,error) {
	ips,err := r.FindIPWithType("https")
	if err!=nil {
		return false,err
	}
	return len(ips)>0,nil
}