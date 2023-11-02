package app

import (
	"context"
	"github.com/networm6/CatTunnel/common/encrypt"
	"github.com/networm6/CatTunnel/common/tools"
	"github.com/networm6/CatTunnel/protocol/ws"
	"github.com/networm6/CatTunnel/protocol/ws/client"
	"github.com/networm6/CatTunnel/protocol/ws/dhcp"
	"github.com/networm6/CatTunnel/protocol/ws/server"
	"github.com/networm6/CatTunnel/tunnel"
	"io"
	"net/http"
	"runtime"
)

// Cat 结构体
type Cat struct {
	Version    string
	LifeCtx    *context.Context
	LifeCancel *context.CancelFunc

	TotalReadBytes    uint64
	TotalWrittenBytes uint64

	_tunnelDev *tunnel.Tunnel
	_appConf   *AppConfig
	_wsConf    *ws.WSConfig
	_tunConf   *tunnel.TunConfig
}

// NewCat 创建。
func NewCat() *Cat {
	ctx, cancel := context.WithCancel(context.Background())
	app := &Cat{
		Version:           "beta3",
		LifeCtx:           &ctx,
		LifeCancel:        &cancel,
		TotalReadBytes:    0,
		TotalWrittenBytes: 0,
		_tunnelDev:        tunnel.NewTunnel(ctx),
	}
	return app
}

// InitApp 初始化
func (cat *Cat) InitApp(conf *AppConfig) {
	tunConf := &tunnel.TunConfig{
		DeviceName: "simonTunnel",
		MTU:        1500,
		ServerMode: conf.ServerMode,
		ServerAddr: conf.ServerAddr,
		BufferSize: 64 * 1024,
		MixinFunc:  conf.MixinFunc,
	}
	if conf.ServerMode {
		serverAddress := tools.Address{
			CIDR:   "172.16.0.1/24",   //必须
			CIDRv6: "fced:9999::1/64", //必须
		}
		tunConf.Address = serverAddress
	} else {
		serverAddress := tools.Address{
			ServerTunnelIP:   "172.16.0.1",
			ServerTunnelIPv6: "fced:9999::1",

			CIDR:   "172.16.0.14/24",
			CIDRv6: "fced:9999::9999/64",
		}
		tunConf.Address = serverAddress
		tunConf.LocalGateway = tools.DiscoverGateway(true)
		tunConf.LocalGatewayv6 = tools.DiscoverGateway(false)
	}

	cat._tunConf = tunConf

	wsConf := &ws.WSConfig{
		ServerAddr: conf.ServerAddr,
		WSPath:     conf.WSPath,
		Timeout:    conf.Timeout,
		Key:        conf.Key,
	}

	cat._wsConf = wsConf

	encrypt.SetMixinKey(conf.Key)
	http.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, cat.PrintBytes(true))
		runtime.Gosched()
	})
	cat._appConf = conf
}

func (cat *Cat) Start() {
	if cat._appConf.ServerMode {
		cat.startServer()
	} else {
		cat.startClient()
	}
}

// StartClient 开始。
func (cat *Cat) startClient() {
	cat._tunnelDev.SetConf(cat._tunConf, &cat.TotalReadBytes, &cat.TotalWrittenBytes)
	cat._tunnelDev.Start()
	client.StartClient(cat._wsConf, cat._tunnelDev)
}

// StartServer 开始。
func (cat *Cat) startServer() {
	cat._tunnelDev.SetConf(cat._tunConf, &cat.TotalReadBytes, &cat.TotalWrittenBytes)
	cat._tunnelDev.Start()
	ws.StartHttpServer(cat._wsConf, cat._tunConf)
	dhcp.StartDHCPServer(dhcp.Config{
		CIDR:   cat._tunConf.Address.CIDR,
		CIDRv6: cat._tunConf.Address.CIDRv6,
		Key:    cat._appConf.Key,
	})
	server.StartServer(cat._wsConf, cat._tunnelDev)
}

// Destroy 结束。
func (cat *Cat) Destroy() {
	cat._tunnelDev.Destroy()
}
