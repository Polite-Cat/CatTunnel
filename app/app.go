package app

import (
	"context"
	"github.com/networm6/PoliteCat/common/encrypt"
	"github.com/networm6/PoliteCat/common/tools"
	"github.com/networm6/PoliteCat/protocol/ws"
	"github.com/networm6/PoliteCat/protocol/ws/client"
	"github.com/networm6/PoliteCat/protocol/ws/server"
	"github.com/networm6/PoliteCat/tunnel"
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
	_wsConf    *ws.WSConfig
	_tunConf   *tunnel.TunConfig
}

// NewCat 创建。
func NewCat() *Cat {
	ctx, cancel := context.WithCancel(context.Background())
	app := &Cat{
		Version:           "beta2",
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
			ServerTunnelIP:   "172.16.0.1",
			ServerTunnelIPv6: "fced:9999::1",

			CIDR:   "172.16.0.1/24",
			CIDRv6: "fced:9999::9999/64",
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
}

// StartClient 开始。
func (cat *Cat) StartClient() {
	cat._tunnelDev.SetConf(cat._tunConf, &cat.TotalReadBytes, &cat.TotalWrittenBytes)
	cat._tunnelDev.Start()
	client.StartClient(cat._wsConf, cat._tunnelDev)
}

// StartServer 开始。
func (cat *Cat) StartServer() {
	cat._tunnelDev.SetConf(cat._tunConf, &cat.TotalReadBytes, &cat.TotalWrittenBytes)
	cat._tunnelDev.Start()
	ws.StartHttpServer(cat._wsConf, cat._tunConf)
	server.StartServer(cat._wsConf, cat._tunnelDev)
}

// Destroy 结束。
func (cat *Cat) Destroy() {
	cat._tunnelDev.Destroy()
}
