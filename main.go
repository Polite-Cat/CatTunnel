package main

import (
	"flag"
	"github.com/networm6/PoliteCat/app"
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
	flag.StringVar(&cfg.MixinFunc, "f", app.DefaultConfig.MixinFunc, "mixin function xor/none")

	flag.Parse()
}

func main() {
	runClient()
}
func runClient() {
	// 创建App和TUN类
	var cat = app.NewCat()
	// 加载TUN和websocket配置
	cat.InitApp(&cfg)
	// 开启
	go cat.StartClient()

	runtime.Gosched()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cat.Destroy()
}
func runServer() {
	// 创建App和TUN类
	var cat = app.NewCat()
	// 加载TUN和websocket配置
	cat.InitApp(&cfg)
	// 开启
	go cat.StartServer()

	runtime.Gosched()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cat.Destroy()
}

/**
未知适配器 simonTunnel:

   连接特定的 DNS 后缀 . . . . . . . : wintun.dns
   IPv6 地址 . . . . . . . . . . . . : fced:9999::9999
   IPv4 地址 . . . . . . . . . . . . : 172.16.0.1
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . :
*/
/**
 -s="8.219.91.90:3001" -c="172.16.0.11/24"
未知适配器 vtun:

   连接特定的 DNS 后缀 . . . . . . . : wintun.dns
   IPv6 地址 . . . . . . . . . . . . : fced:9999::9999
   IPv4 地址 . . . . . . . . . . . . : 172.16.0.11
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . :

*/
/**
-s="8.219.91.90:3001" -c="172.16.0.13/24" -sip="172.16.0.1" -g
未知适配器 vtun:

   连接特定的 DNS 后缀 . . . . . . . : wintun.dns
   IPv6 地址 . . . . . . . . . . . . : fced:9999::9999
   IPv4 地址 . . . . . . . . . . . . : 172.16.0.13
   子网掩码  . . . . . . . . . . . . : 255.255.255.0
   默认网关. . . . . . . . . . . . . : 172.16.0.1
*/
