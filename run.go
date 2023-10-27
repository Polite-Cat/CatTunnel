package main

import (
	"github.com/networm6/PoliteCat/app"
	"os"
	"os/signal"
	"syscall"
)

func runClient() {
	// 创建App和TUN类
	var cat = app.NewCat()
	// 加载TUN和websocket配置
	cat.InitApp(&cfg)
	// 开启
	go cat.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cat.Destroy()
}
