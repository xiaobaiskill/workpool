package server

import (
	"fmt"
	"github.com/chenjiandongx/ginprom"
	"github.com/go-clog/clog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/natefinch/lumberjack"
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/log"
	"github.com/xiaobaiskill/workpool/pkg/models"
	"github.com/xiaobaiskill/workpool/pkg/redis"
	"github.com/xiaobaiskill/workpool/routes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	mysql_log "log"
	"os"
	"path"
)

type Server struct {
	cfg *conf.Config
	db  models.DB
}

func (s *Server) Run() {
	router := routes.NewRouter()

	// router.use   //可加中间件
	router.Use(ginprom.PromMiddleware(nil))

	router.Run(fmt.Sprintf("%s:%v", s.cfg.Server.HttpAddr, s.cfg.Server.HttpPort))
}

// 日志配置
func (s *Server) configureLog() {
	levelNames := map[string]zapcore.Level{
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		level, ok := levelNames[s.cfg.LogFile.Level]
		if ok {
			return lvl >= level
		}
		return lvl >= zapcore.ErrorLevel
	})
	hook := lumberjack.Logger{
		Filename:   path.Join(s.cfg.Log.RootPath, "log.txt"),
		MaxSize:    1 << s.cfg.LogFile.MaxSizeShift, // megabytes
		MaxBackups: s.cfg.LogFile.MaxBackups,
		MaxAge:     s.cfg.LogFile.MaxDays,  //days
		Compress:   s.cfg.LogFile.Compress, // disabled by default
	}
	fileWriter := zapcore.AddSync(&hook)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	var allCore []zapcore.Core
	if s.cfg.Debug {
		consoleDebugging := zapcore.Lock(os.Stdout)
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			level, ok := levelNames[s.cfg.LogConsole.Level]
			if ok {
				return lvl >= level
			}
			return lvl >= zapcore.DebugLevel
		})
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	}
	allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, highPriority))

	core := zapcore.NewTee(allCore...)
	log.InitLog(core, s.cfg.AppName)
	log.Logger.Info("程序日志启动。。。")
}

func (s *Server) newDB() (models.DB, error) {
	if s.cfg.Redis.On {
		conn := redis.New(s.cfg.Redis.Addr)
		if _, err := conn.Ping(); err != nil {
			panic("redis 链接失败")
		}
		return models.NewRedisDB(conn), nil
	} else {
		g, err := gorm.Open(s.cfg.Database.Driver, s.cfg.Database.DataSource)

		if err != nil {
			return nil, err
		}

		// configure xorm
		err = s.configureGORM(g)
		if err != nil {
			return nil, err
		}
		return models.NewDefaultDB(g), nil
	}
}

// db 的配置
func (s *Server) configureGORM(g *gorm.DB) error {
	logger, err := clog.NewFileWriter(path.Join(s.cfg.Log.RootPath, "gorm.log"),
		clog.FileRotationConfig{
			Rotate:  s.cfg.LogXorm.Rotate,
			Daily:   s.cfg.LogXorm.RotateDaily,
			MaxSize: s.cfg.LogXorm.MaxSize * 1024 * 1024,
			MaxDays: s.cfg.LogXorm.MaxDays,
		})
	if err != nil {
		return fmt.Errorf("Fail to create 'xorm.log': %s", err)
	}

	g.SetLogger(mysql_log.New(logger, "\r\n", 0))

	g.AutoMigrate(models.NewIP())
	return nil
}

func New(cfg *conf.Config) *Server {
	s := &Server{cfg: cfg}
	s.configureLog()
	var err error
	models.Conn, err = s.newDB()

	if err != nil {
		panic(err)
	}
	return s
}
