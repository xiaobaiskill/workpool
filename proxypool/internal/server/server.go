package server

import (
	"fmt"
	sj "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/go-clog/clog"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/parnurzeal/gorequest"
	"github.com/ruoklive/proxypool/pkg/models"
	"github.com/ruoklive/proxypool/pkg/pool"
	"github.com/ruoklive/proxypool/pkg/redis"
	"github.com/ruoklive/proxypool/pkg/register"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

const VERSION = "/v2"

type Server struct {
	cfg           *Config
	db            models.DB
	checkChanSize int // check ip channel size

	workPool pool.Collector
}

func New(cfg *Config) *Server {
	s := &Server{
		cfg:           cfg,
		checkChanSize: 50,
		workPool: pool.StartDispatcher(cfg.Pool.WorkSize),
	}
	s.configureLog()
	// init db
	var err error
	s.db, err = s.newDB()
	if err != nil {
		panic(err)
	}
	// check proxy in db
	go s.CheckProxyDB()
	// start ip getter
	go s.startIPGetter()

	return s
}

func (s *Server) Run() error {
	r := gin.Default()
	// random ip
	r.GET(VERSION+"/ip", s.ProxyRandom)
	r.GET(VERSION+"/https", s.ProxyFind)
	r.GET(VERSION+"/total", s.HealthCheck)
	r.GET(VERSION+"/health", s.HealthCheck)
	r.GET(VERSION+"/monitor",s.Monitor)
	return r.Run(fmt.Sprintf("%s:%d", s.cfg.Server.HttpAddr, s.cfg.Server.HttpPort))
}


// start ip getter
func (s *Server) startIPGetter() {
	go func() {
		for {
			n, err := s.db.CountIP()
			if err != nil {
				panic(err)
			}
			log.Printf("IP: %v\n", n)
			go s.executeIPGetter()
			time.Sleep(10 * time.Minute)
		}
	}()
}

func (s *Server) executeIPGetter() {
	var wg sync.WaitGroup
	executors := register.GetExecutors()
	wg.Add(len(executors))
	for _, executor := range executors {
		go func(executor register.Executor) {
			defer func() {
				if err := recover(); err != nil {
					clog.Warn("执行失败！", err)
				}
				wg.Done()
			}()
			ips := executor().Execute()
			for _, ip := range ips {
				s.workPool.Work <- &pool.Job{
					Data: ip,
					JobFunc: func(id int64, data interface{}) {
						s.CheckProxy(data.(*models.IP))
					},
				}
			}
		}(executor)
	}
	wg.Wait()
	log.Println("All getters finished.")
}

func (s *Server) newDB() (models.DB, error) {

	if s.cfg.Redis.On {
		conn := redis.New(s.cfg.Redis.Addr)
		return models.NewRedisDB(conn), nil
	} else {
		x, err := xorm.NewEngine(s.cfg.Database.Driver, s.cfg.Database.DataSource)
		if err != nil {
			return nil, err
		}
		// configure xorm
		err = s.configureXORM(x)
		if err != nil {
			return nil, err
		}
		return models.NewDefaultDB(x), nil
	}
}

// CheckProxyDB to check the ip in DB
func (s *Server) CheckProxyDB() {
	count, err := s.db.CountIP()
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	clog.Info("Before check, DB has: %d records.", count)
	ips, err := s.db.GetAllIP()
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			s.workPool.Work <- &pool.Job{
				Data: v,
				JobFunc: func(id int64, data interface{}) {
					ip :=data.(*models.IP)
					if !s.CheckIP(ip) {
						err = s.db.DeleteIP(ip)
						if err != nil {
							clog.Warn(err.Error())
						}
					}
				},
			}

			wg.Done()
		}(v)
	}
	wg.Wait()
	count, err = s.db.CountIP()
	if err != nil {
		clog.Warn(err.Error())
		return
	}
	clog.Info("After check, DB has: %d records.", count)
}

// CheckIP is to check the ip work or not
func (s *Server) CheckIP(ip *models.IP) bool {
	var pollURL string
	var testIP string
	if ip.Type2 == "https" {
		testIP = "https://" + ip.Data
		pollURL = "https://httpbin.org/get"
	} else {
		testIP = "http://" + ip.Data
		pollURL = "http://httpbin.org/get"
	}
	clog.Info(testIP)
	begin := time.Now()
	resp, _, errs := gorequest.New().Timeout(time.Second * 20).Proxy(testIP).Get(pollURL).End()
	if errs != nil {
		clog.Warn("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, errs)
		return false
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		//harrybi 20180815 判断返回的数据格式合法性
		_, err := sj.NewFromReader(resp.Body)
		if err != nil {
			clog.Warn("[CheckIP] testIP = %s, pollURL = %s: Error = %v", testIP, pollURL, err)
			return false
		}
		//harrybi 计算该代理的速度，单位毫秒
		ip.Speed = time.Now().Sub(begin).Nanoseconds() / 1000 / 1000 //ms
		//TODO 这里为什么更新第一条IP？这样不是把第一条IP的数据覆盖了？ 没太明白 先注释掉
		//if err = s.db.UpdateToFirstIP(ip); err != nil {
		//	clog.Warn("[CheckIP] Update IP = %v Error = %v", *ip, err)
		//}
		clog.Info("IP[%s] is OK", ip)

		return true
	}
	return false
}

// CheckProxy .
func (s *Server) CheckProxy(ip *models.IP) {
	if s.CheckIP(ip) {
		err := s.db.InsertIP(ip)
		if err != nil {
			clog.Info("IP[%s] check error = %v", ip, err)
		}
	}
}


func (s *Server) configureXORM(x *xorm.Engine) error {
	x.SetMapper(core.GonicMapper{})
	logger, err := clog.NewFileWriter(path.Join(s.cfg.Log.RootPath, "xorm.log"),
		clog.FileRotationConfig{
			Rotate:  s.cfg.LogXorm.Rotate,
			Daily:   s.cfg.LogXorm.RotateDaily,
			MaxSize: s.cfg.LogXorm.MaxSize * 1024 * 1024,
			MaxDays: s.cfg.LogXorm.MaxDays,
		})
	if err != nil {
		return fmt.Errorf("Fail to create 'xorm.log': %s", err)
	}

	x.SetLogger(xorm.NewSimpleLogger(logger))

	if err := x.StoreEngine("InnoDB").Sync2(models.NewIP()); err != nil {
		return fmt.Errorf("sync database struct error: %v", err)
	}
	x.ShowSQL(true)
	return nil
}


func (s *Server) configureLog()  {
	// Because we always create a console logger as primary logger before all settings are loaded,
	// thus if user doesn't set console logger, we should remove it after other loggers are created.
	hasConsole := false

	// Get the log mode.
	var logModes []string
	if s.cfg.Debug {
		logModes = strings.Split("console", ",")
	} else {
		logModes = strings.Split(s.cfg.Log.Mode, ",")
	}
	logConfigs := make([]interface{}, len(logModes))
	levelNames := map[string]clog.LEVEL{
		"trace": clog.TRACE,
		"info":  clog.INFO,
		"warn":  clog.WARN,
		"error": clog.ERROR,
		"fatal": clog.FATAL,
	}
	for i, mode := range logModes {
		mode = strings.ToLower(strings.TrimSpace(mode))
		//sec, err := Cfg.GetSection("log." + mode)
		//if err != nil {
		//	clog.Fatal(2, "Unknown logger mode: %s", mode)
		//}
		var err error
		var levelName string
		// Generate log configuration.
		switch clog.MODE(mode) {
		case clog.CONSOLE:
			hasConsole = true
			levelName = s.cfg.LogConsole.Level
			level := levelNames[levelName]
			logConfigs[i] = clog.ConsoleConfig{
				Level:      level,
				BufferSize:s.cfg.Log.BufferLen,
			}

		case clog.FILE:
			levelName = s.cfg.LogFile.Level
			level := levelNames[levelName]
			logPath := path.Join(s.cfg.Log.RootPath, "ProxyPool.log")
			if err = os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				clog.Fatal(2, "Fail to create log directory '%s': %v", path.Dir(logPath), err)
			}

			logConfigs[i] = clog.FileConfig{
				Level:      level,
				BufferSize: s.cfg.Log.BufferLen,
				Filename:   logPath,
				FileRotationConfig: clog.FileRotationConfig{
					Rotate:   s.cfg.LogFile.LogRotate,
					Daily:    s.cfg.LogFile.DailyRotate,
					MaxSize:  1 << uint(s.cfg.LogFile.MaxSizeShift),
					MaxLines: s.cfg.LogFile.MaxLines,
					MaxDays:  s.cfg.LogFile.MaxDays,
				},
			}
		}

		clog.New(clog.MODE(mode), logConfigs[i])
		clog.Trace("Log Mode: %s (%s)", strings.Title(mode), strings.Title(levelName))
	}

	// Make sure everyone gets version info printed.
	clog.Info("%s %s", s.cfg.AppName, s.cfg.AppVer)
	if !hasConsole {
		clog.Delete(clog.CONSOLE)
	}
}
