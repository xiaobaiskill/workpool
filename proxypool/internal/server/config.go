package server

import (
	"github.com/go-ini/ini"
	"time"
)

type Config struct {
	AppName string
	AppVer string
	Debug bool
	Pool struct{
		WorkSize int64
	}
	Server  struct {
		HttpAddr       string
		HttpPort       int
		SessionExpires time.Duration
	}
	Redis struct {
		On   bool
		Addr string
	}
	Database struct {
		Driver,DataSource string
	}
	Log struct {
		Mode, Level, RootPath string
		BufferLen             int64
	}
	LogConsole struct {
		Level string
	}
	LogFile struct {
		Level        string
		LogRotate    bool
		DailyRotate  bool
		MaxSizeShift int
		MaxLines     int64
		MaxDays      int64
	}
	LogXorm struct {
		Rotate      bool
		RotateDaily bool
		MaxSize     int64 // M
		MaxDays     int64
	}
	Security struct {
		InstallLock bool
	}
}

func NewConfig() *Config {
	c := &Config{}
	c.AppName = "ProxyPool"
	c.AppVer = "1.0"
	c.Debug = false
	// pool
	c.Pool.WorkSize = 100
	// server
	c.Server.HttpAddr = "0.0.0.0"
	c.Server.HttpPort = 3000
	c.Server.SessionExpires = time.Hour * 168
	// redis
	c.Redis.On = true
	c.Redis.Addr = "127.0.0.1:6379"
	// database
	c.Database.Driver = "postgres"
	c.Database.DataSource = "postgres://postgres:example@127.0.0.1:5432/public?sslmode=disable"
	// log
	c.Log.Mode = "file"
	c.Log.BufferLen = 100
	c.Log.Level = "info"
	c.Log.RootPath = ""
	// log console
	c.LogConsole.Level = "Trace"
	// log file
	c.LogFile.Level = "Info"
	c.LogFile.LogRotate = true
	c.LogFile.DailyRotate = true
	c.LogFile.MaxSizeShift = 28
	c.LogFile.MaxLines = 1000000
	c.LogFile.MaxDays = 7
	// log xorm
	c.LogXorm.Rotate = true
	c.LogXorm.RotateDaily = true
	c.LogXorm.MaxSize = 100
	c.LogXorm.MaxDays = 3
	// security
	c.Security.InstallLock = true
	return c
}

// Load 加载配置
func (c *Config) Load(confFile string) {
	cfg, err := ini.Load(confFile)
	if err != nil {
		panic(err)
	}
	c.AppName = cfg.Section("").Key("APP_NAME").MustString(c.AppName)
	c.AppVer = cfg.Section("").Key("APP_VER").MustString(c.AppVer)
	c.Debug = cfg.Section("").Key("DEBUG").MustBool(c.Debug)
	// pool
	c.Pool.WorkSize = cfg.Section("pool").Key("WORK_SIZE").MustInt64(c.Pool.WorkSize)
	// server
	c.Server.HttpAddr = cfg.Section("server").Key("HTTP_ADDR").MustString(c.Server.HttpAddr)
	c.Server.HttpPort = cfg.Section("server").Key("HTTP_PORT").MustInt(c.Server.HttpPort)
	c.Server.SessionExpires = cfg.Section("server").Key("SESSION_EXPIRES").MustDuration(c.Server.SessionExpires)
	// redis
	c.Redis.On = cfg.Section("redis").Key("ON").MustBool(c.Redis.On)
	c.Redis.Addr = cfg.Section("redis").Key("ADDR").MustString(c.Redis.Addr)
	// database
	c.Database.Driver = cfg.Section("database").Key("DRIVER").MustString(c.Database.Driver)
	c.Database.DataSource = cfg.Section("database").Key("DATA_SOURCE").MustString(c.Database.DataSource)
	// log
	c.Log.Mode = cfg.Section("log").Key("MODE").MustString(c.Log.Mode)
	c.Log.BufferLen = cfg.Section("log").Key("BUFFER_LEN").MustInt64(c.Log.BufferLen)
	c.Log.Level = cfg.Section("log").Key("LEVEL").MustString(c.Log.Level)
	c.Log.RootPath = cfg.Section("log").Key("ROOT_PATH").MustString(c.Log.RootPath)
	// log.console
	c.LogConsole.Level = cfg.Section("log.console").Key("LEVEL").MustString(c.LogConsole.Level)
	// log.file
	c.LogFile.Level = cfg.Section("log.file").Key("LEVEL").MustString(c.LogFile.Level)
	c.LogFile.LogRotate = cfg.Section("log.file").Key("LOG_ROTATE").MustBool(c.LogFile.LogRotate)
	c.LogFile.DailyRotate = cfg.Section("log.file").Key("DAILY_ROTATE").MustBool(c.LogFile.DailyRotate)
	c.LogFile.MaxSizeShift = cfg.Section("log.file").Key("MAX_SIZE_SHIFT").MustInt(c.LogFile.MaxSizeShift)
	c.LogFile.MaxLines = cfg.Section("log.file").Key("MAX_LINES").MustInt64(c.LogFile.MaxLines)
	c.LogFile.MaxDays = cfg.Section("log.file").Key("MAX_DAYS").MustInt64(c.LogFile.MaxDays)
	// log.xorm
	c.LogXorm.Rotate = cfg.Section("log.xorm").Key("ROTATE").MustBool(c.LogXorm.Rotate)
	c.LogXorm.RotateDaily = cfg.Section("log.xorm").Key("ROTATE_DAILY").MustBool(c.LogXorm.RotateDaily)
	c.LogXorm.MaxSize = cfg.Section("log.xorm").Key("MAX_SIZE").MustInt64(c.LogXorm.MaxSize)
	c.LogXorm.MaxDays = cfg.Section("log.xorm").Key("MAX_DAYS").MustInt64(c.LogXorm.MaxDays)

	//security
	c.Security.InstallLock = cfg.Section("security").Key("INSTALL_LOCK").MustBool(c.Security.InstallLock)

}
