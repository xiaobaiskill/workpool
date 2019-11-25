package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

var levelNames = map[string]zapcore.Level{
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func TestInitLog(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("日志有误")
		} else {
			t.Log("日志记录成功")
		}
	}()

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	var allCore []zapcore.Core

	consoleDebugging := zapcore.Lock(os.Stdout)
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		level, ok := levelNames["Debug"]
		if ok {
			return lvl >= level
		}
		return lvl >= zapcore.DebugLevel
	})
	allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))

	core := zapcore.NewTee(allCore...)
	InitLog(core, "test")

	Logger.Info("Info 测试")
}

func TestFileLog(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Fatal("日志有误")
		} else {
			t.Log("日志记录成功")
		}
	}()
	hook := lumberjack.Logger{
		Filename:   "../../logs/log.txt",
		MaxSize:    1 << 28, // megabytes
		MaxBackups: 2,
		MaxAge:     7,  //days
		Compress:   false, // disabled by default
	}
	fileWriter := zapcore.AddSync(&hook)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	var allCore []zapcore.Core

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		level, ok := levelNames["info"]
		if ok {
			return lvl >= level
		}
		return lvl >= zapcore.DebugLevel
	})
	allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, highPriority))

	core := zapcore.NewTee(allCore...)
	InitLog(core, "test")

	Logger.Info("info 测试")

}
