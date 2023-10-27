package app

import (
	"context"
	ws2 "github.com/networm6/PoliteCat/app/ws"
	"github.com/networm6/PoliteCat/common/encrypt"
	"github.com/networm6/PoliteCat/common/tools"
	"github.com/networm6/PoliteCat/tunnel"
	"log"
	"os"
)

// Cat 结构体
type Cat struct {
	Version    string
	LifeCtx    *context.Context
	LifeCancel *context.CancelFunc

	TotalReadBytes    uint64
	TotalWrittenBytes uint64

	_tunnelDev *tunnel.Tunnel
	_wsConf    *ws2.Config
	_tunConf   *tunnel.Config
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
func (cat *Cat) InitApp(conf *Config) {
	tunConf := &tunnel.Config{
		DeviceName:     "simonTunnel",
		MTU:            1500,
		ServerAddr:     conf.ServerAddr,
		BufferSize:     64 * 1024,
		MixinFunc:      conf.MixinFunc,
		LocalGateway:   tools.DiscoverGateway(true),
		LocalGatewayv6: tools.DiscoverGateway(false),
	}
	address, err := tools.InitAddress(conf.ServerAddr+"/address", conf.Key)
	if err != nil {
		log.Fatalf("App error :%v", err)
		os.Exit(-1)
	}
	tunConf.Address = *address

	cat._tunConf = tunConf

	wsConf := &ws2.Config{
		ServerAddr: conf.ServerAddr,
		WSPath:     conf.WSPath,
		Timeout:    conf.Timeout,
		Key:        conf.Key,
	}

	cat._wsConf = wsConf

	encrypt.SetMixinKey(conf.Key)
}

// Start 开始。
func (cat *Cat) Start() {
	cat._tunnelDev.SetConf(cat._tunConf, &cat.TotalReadBytes, &cat.TotalWrittenBytes)
	ws2.StartClient(cat._wsConf, cat._tunnelDev)
}

// Destroy 结束。
func (cat *Cat) Destroy() {
	cat._tunnelDev.Destroy()
}
