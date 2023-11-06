package app

import (
	"context"
	"github.com/networm6/CatTunnel/protocol/dhcp"
	"github.com/networm6/CatTunnel/protocol/dhcp/server"
	"github.com/networm6/gopherBox/lifecycle"
	"github.com/networm6/gopherBox/tunnel"
	"io"
	"net/http"
	"runtime"
)

// Cat 结构体
type Cat struct {
	lifecycle.LifeInterface
	Version    string
	LifeCtx    *context.Context
	LifeCancel *context.CancelFunc

	TotalReadBytes    uint64
	TotalWrittenBytes uint64

	_tunnelDev *tunnel.Tunnel
	_appConf   *AppConfig
}

// NewCat 创建。
func NewCat() *Cat {
	ctx, cancel := context.WithCancel(context.Background())
	app := &Cat{
		Version:           "beta4",
		LifeCtx:           &ctx,
		LifeCancel:        &cancel,
		TotalReadBytes:    0,
		TotalWrittenBytes: 0,
	}
	//TODO: 内存映射到某个文件，这样可以断网保存
	app._tunnelDev = tunnel.NewTunnel(ctx, &app.TotalReadBytes, &app.TotalWrittenBytes)
	return app
}

// InitApp 初始化
func (cat *Cat) InitApp(conf *AppConfig) {
	cat._appConf = conf
	http.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, cat.PrintBytes(true))
		runtime.Gosched()
	})
}

// Start 运行。
func (cat *Cat) Start() {
	if cat._appConf.ServerMode {
		cat._startServer()
	} else {
		cat._startClient()
	}
}

// Destroy 结束。
func (cat *Cat) Destroy() {
	cat._tunnelDev.Destroy()
}

// _startClient
func (cat *Cat) _startClient() {

}

// _startServer 负责启动DHCP端，http端，ws端
func (cat *Cat) _startServer() {
	key := cat._appConf.Key
	server.StartDHCPServer(dhcp.Config{
		CIDR:   "",
		CIDRv6: "",
		Key:    key,
	})
}
