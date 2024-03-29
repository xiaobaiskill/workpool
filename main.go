package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xiaobaiskill/workpool/internal/server"
	"github.com/xiaobaiskill/workpool/pkg/conf"
	"github.com/xiaobaiskill/workpool/pkg/pool"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := conf.NewConfig()
	fmt.Println(os.Getpid())
	if os.Getenv("GIN_MODE") == gin.ReleaseMode {
		cfg.Load("conf/app.ini")
	} else {
		cfg.Load("conf/app_dev.ini")
	}
	s := server.New(cfg)
	go s.Run()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println(sig)

		pool.StopDispathcher()
		time.Sleep(2 * time.Second)

		done <- true
	}()

	fmt.Println("workpool 程序启动")
	<-done
	fmt.Println("workpool  exiting")
}
