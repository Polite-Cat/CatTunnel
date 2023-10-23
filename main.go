package main

import (
	"flag"
	"github.com/networm6/PoliteCat/daemon/service"
	"github.com/networm6/PoliteCat/protocol/tcp"
)

func main() {
	appConfig := service.Config{}
	flag.StringVar(&appConfig.CIDR, "c", "172.16.0.10/24", "tun interface cidr")
	flag.IntVar(&appConfig.MTU, "mtu", 1500, "tun mtu")
	flag.StringVar(&appConfig.LocalAddr, "l", ":3000", "local address")
	flag.StringVar(&appConfig.ServerAddr, "s", ":3001", "server address")
	flag.StringVar(&appConfig.Key, "k", "123456", "key")
	flag.StringVar(&appConfig.DNS, "d", "8.8.8.8:53", "dns address")
	flag.BoolVar(&appConfig.ServerMode, "S", false, "server mode")
	flag.BoolVar(&appConfig.GlobalMode, "g", false, "client global mode")
	flag.BoolVar(&appConfig.Enc, "enc", false, "enable data encry")
	flag.IntVar(&appConfig.Timeout, "t", 30, "dial timeout in seconds")
	flag.Parse()
	appConfig.Init()

	if appConfig.ServerMode {
		tcp.StartServer(appConfig)
	} else {
		tcp.StartClient(appConfig)
	}
}
