package main

import (
	"flag"
	"github.com/networm6/CatTunnel/app"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var cfg = app.AppConfig{}

func init() {
	flag.BoolVar(&cfg.ServerMode, "S", app.DefaultConfig.ServerMode, "server mode")
	flag.StringVar(&cfg.ServerAddr, "s", app.DefaultConfig.ServerAddr, "server address")
	flag.StringVar(&cfg.Key, "k", app.DefaultConfig.Key, "key")

	flag.StringVar(&cfg.WSPath, "path", app.DefaultConfig.WSPath, "ws path")
	flag.IntVar(&cfg.Timeout, "t", app.DefaultConfig.Timeout, "dial timeout in seconds")
	flag.StringVar(&cfg.MixinFunc, "f", app.DefaultConfig.MixinFunc, "mixin function xor/none/reverse")

	flag.Parse()
}

func main() {
	// 创建App和TUN类
	var cat = app.NewCat()
	// 加载TUN和websocket配置
	cat.InitApp(&cfg)
	// 开启
	go cat.Start()

	runtime.Gosched()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cat.Destroy()
}
