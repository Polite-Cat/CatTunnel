package tunnel

import (
	"context"
	"github.com/net-byte/water"
)

// Tunnel 结构体
type Tunnel struct {
	OutputStream chan []byte
	InputStream  chan []byte
	LifeCtx      *context.Context
	LifeCancel   *context.CancelFunc

	_conf         *TunConfig
	_tunInterface *water.Interface
	_mixinFunc    func([]byte) []byte
	_bufferSize   int

	_totalReadBytes    *uint64
	_totalWrittenBytes *uint64
}

// NewTunnel 创建。
func NewTunnel(parentCtx context.Context) *Tunnel {
	_ctx, _cancel := context.WithCancel(parentCtx)
	tunnel := &Tunnel{
		OutputStream: make(chan []byte),
		InputStream:  make(chan []byte),
		LifeCtx:      &_ctx,
		LifeCancel:   &_cancel,
	}
	return tunnel
}

func (tun *Tunnel) SetConf(conf *TunConfig, readBytes, writtenBytes *uint64) {
	tun._conf = conf
	tun._tunInterface = CreateTunnelInterface(*tun._conf)
	tun._mixinFunc = getFunc(tun._conf.MixinFunc)
	tun._bufferSize = tun._conf.BufferSize
	tun._totalReadBytes = readBytes
	tun._totalWrittenBytes = writtenBytes
}

func (tun *Tunnel) Start() {
	SetRoute(*tun._conf, tun._tunInterface)
	go tun.readFromTunnel()
	go tun.writeToTunnel()
}

func (tun *Tunnel) Destroy() {
	(*tun.LifeCancel)()
	ResetRoute(*tun._conf)
	_ = tun._tunInterface.Close()
	close(tun.OutputStream)
	close(tun.InputStream)
}
