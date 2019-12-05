package conf

import (
	"github.com/xiaobaiskill/workpool/pkg/setting"
	"time"
)
var Conf *Config
type Config struct {
	AppName string
	AppVer string
	Debug bool
	Pool struct{
		WorkSize int
		TimeOut int
		ProxyIpSize int
		RetryMax int
	}
	Redis struct {
		On   bool
		Addr string
	}
	Server  struct {
		HttpAddr       string
		HttpPort       int
		SessionExpires time.Duration
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
		MaxDays      int
		Compress     bool
		MaxBackups   int
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
	Proxypool struct{
		PublicproxyMinsize int
		PublicproxyRetryNum int
		BackupproxyConf string
		BackupproxyRetryNum int
		SelfUrl string
		SelfRetryNum int
	}
}

func NewConfig() *Config {
	c := &Config{}
	c.AppName = "WorkPool"
	c.AppVer = "1.0"
	c.Debug = true

	c.Pool.WorkSize = 50
	c.Pool.TimeOut = 4
	c.Pool.ProxyIpSize = 50
	// redis
	c.Redis.On = false
	c.Redis.Addr = "127.0.0.1:6379"

	// server
	c.Server.HttpAddr = "0.0.0.0"
	c.Server.HttpPort = 3000
	c.Server.SessionExpires = time.Hour * 168

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
	c.LogFile.Compress = false
	c.LogFile.MaxBackups = 2
	// log xorm
	c.LogXorm.Rotate = true
	c.LogXorm.RotateDaily = true
	c.LogXorm.MaxSize = 100
	c.LogXorm.MaxDays = 3
	// security
	c.Security.InstallLock = true
	// proxypool
	c.Proxypool.PublicproxyMinsize = 50
	c.Proxypool.PublicproxyRetryNum = 5
	c.Proxypool.BackupproxyConf = "conf/backupproxy.conf"
	c.Proxypool.BackupproxyRetryNum = -1
	c.Proxypool.SelfUrl = "https://hooks.slack.com/services/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"
	c.Proxypool.SelfRetryNum = 1

	Conf = c
	return c
}

// Load 加载配置
func (c *Config) Load(confFile string) {
	ini_file := confFile
	ini, err := setting.NewContext(ini_file)
	if err != nil {
		panic("文件有误:" + err.Error())
	}

	c.AppName = ini.GetString("", "APP_NAME", c.AppName)
	c.AppVer = ini.GetString("", "APP_VER", c.AppVer)
	c.Debug = ini.GetBool("", "APP_DEBUG", c.Debug)

	// workpool
	c.Pool.WorkSize = ini.GetInt("workpool","WORK_SIZE",c.Pool.WorkSize )
	c.Pool.TimeOut = ini.GetInt("workpool","TIME_OUT",c.Pool.TimeOut)
	c.Pool.ProxyIpSize = ini.GetInt("workpool","PROXY_IP_SIZE",c.Pool.ProxyIpSize)
	c.Pool.RetryMax = ini.GetInt("workpool","RETRY_MAX",c.Pool.RetryMax)
	// redis
	c.Redis.On = ini.GetBool("redis","ON",c.Redis.On)
	c.Redis.Addr = ini.GetString("redis","ADDR",c.Redis.Addr)
	// server
	c.Server.HttpAddr = ini.GetString("server","HTTP_ADDR",c.Server.HttpAddr)
	c.Server.HttpPort = ini.GetInt("server","HTTP_PORT",c.Server.HttpPort)
	c.Server.SessionExpires = ini.GetDuration("server","SESSION_EXPIRES",c.Server.SessionExpires)
	// database
	c.Database.Driver = ini.GetString("database","DRIVER",c.Database.Driver)
	c.Database.DataSource = ini.GetString("database","DATA_SOURCE",c.Database.DataSource)
	// log
	c.Log.Mode = ini.GetString("log","MODE",c.Log.Mode)
	c.Log.BufferLen = ini.GetInt64("log","BUFFER_LEN",c.Log.BufferLen)
	c.Log.Level = ini.GetString("log","LEVEL",c.Log.Level)
	c.Log.RootPath = ini.GetString("log","ROOT_PATH",c.Log.RootPath)
	// log.console
	c.LogConsole.Level = ini.GetString("log.console","LEVEL",c.LogConsole.Level)
	// log.file
	c.LogFile.Level = ini.GetString("log.file","LEVEL",c.LogFile.Level)
	c.LogFile.LogRotate = ini.GetBool("log.file","LOG_ROTATE",c.LogFile.LogRotate)
	c.LogFile.DailyRotate = ini.GetBool("log.file","DAILY_ROTATE",c.LogFile.DailyRotate)
	c.LogFile.MaxSizeShift = ini.GetInt("log.file","MAX_SIZE_SHIFT",c.LogFile.MaxSizeShift)
	c.LogFile.MaxLines = ini.GetInt64("log.file","MAX_LINES",c.LogFile.MaxLines)
	c.LogFile.MaxDays = ini.GetInt("log.file","MAX_DAYS",c.LogFile.MaxDays)
	c.LogFile.Compress = ini.GetBool("log.file","COMPRESS",c.LogFile.Compress)
	c.LogFile.MaxBackups = ini.GetInt("log.file","MAX_BACKUPS",c.LogFile.MaxBackups)
	// log.xorm
	c.LogXorm.Rotate = ini.GetBool("log.xorm","ROTATE",c.LogXorm.Rotate)
	c.LogXorm.RotateDaily = ini.GetBool("log.xorm","ROTATE_DAILY",c.LogXorm.RotateDaily)
	c.LogXorm.MaxSize = ini.GetInt64("log.xorm","MAX_SIZE",c.LogXorm.MaxSize)
	c.LogXorm.MaxDays = ini.GetInt64("log.xorm","MAX_DAYS",c.LogXorm.MaxDays)

	//security
	c.Security.InstallLock = ini.GetBool("security","INSTALL_LOCK",c.Security.InstallLock)

	// proxypool
	c.Proxypool.PublicproxyMinsize = ini.GetInt("proxypool","PUBLICPROXYMINSIZE",c.Proxypool.PublicproxyMinsize)
	c.Proxypool.PublicproxyRetryNum = ini.GetInt("proxypool","PUBLICPROXYRETRYNUM",c.Proxypool.PublicproxyRetryNum)
	c.Proxypool.BackupproxyConf = ini.GetString("proxypool","BACKUPPROXYCONF",c.Proxypool.BackupproxyConf)
	c.Proxypool.BackupproxyRetryNum = ini.GetInt("proxypool","BACKUPPROXYRETRYNUM",c.Proxypool.BackupproxyRetryNum)
	c.Proxypool.SelfUrl = ini.GetString("proxypool","SELFURL",c.Proxypool.SelfUrl)
	c.Proxypool.SelfRetryNum = ini.GetInt("proxypool","SELFRETRYNUM",c.Proxypool.SelfRetryNum)
}
