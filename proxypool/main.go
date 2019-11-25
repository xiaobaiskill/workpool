package main

import (
	_ "github.com/ruoklive/proxypool/internal/ipgetter/all"
	"github.com/ruoklive/proxypool/internal/server"
	"os"
)

const (
	Debug = "debug"
	Release = "release"
)

func main() {
	// init config
	cfg := server.NewConfig()
	// load config
	if os.Getenv("MODE") == Release {
		cfg.Load("conf/app.ini")
	}else {
		cfg.Load("conf/app_dev.ini")
	}

	// new server
	s := server.New(cfg)
	// run
	err := s.Run()
	if err!=nil {
		panic(err)
	}
}