package setting

import (
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

var (
	// 文件
	Settings map[string]*IniParse
)

func init() {
	Settings = make(map[string]*IniParse)
}

type IniParse struct {
	ConfReader *ini.File
}

func (this *IniParse) GetString(section, key, other string) string {
	if this.ConfReader == nil {
		return other
	}

	s := this.ConfReader.Section(section)
	if s == nil {
		return other
	}

	return s.Key(key).MustString(other)
}

func (this *IniParse) GetInt(section, key string, other int) int {
	if this.ConfReader == nil {
		return other
	}

	s := this.ConfReader.Section(section)
	if s == nil {
		return other
	}

	return s.Key(key).MustInt(other)
}

func (this *IniParse) GetInt64(section, key string, other int64) int64 {
	if this.ConfReader == nil {
		return other
	}

	s := this.ConfReader.Section(section)
	if s == nil {
		return other
	}

	return s.Key(key).MustInt64(other)
}

func (this *IniParse) GetBool(section, key string, other bool) bool {
	if this.ConfReader == nil {
		return other
	}

	s := this.ConfReader.Section(section)
	if s == nil {
		return other
	}

	return s.Key(key).MustBool(other)
}

func (this *IniParse) GetDuration(section, key string, other time.Duration) time.Duration {
	if this.ConfReader == nil {
		return other
	}

	s := this.ConfReader.Section(section)
	if s == nil {
		return other
	}

	return s.Key(key).MustDuration(other)
}

func NewContext(file string) (conf *IniParse, err error) {
	if Settings[file] != nil {
		return Settings[file], nil
	}

	f, err := ini.Load(file)
	if err != nil {
		fmt.Printf("文件：%s 加载失败\n", file)
		return
	}
	conf = &IniParse{f}
	Settings[file] = conf
	return
}
