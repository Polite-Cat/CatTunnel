package ws

import (
	"context"
	"github.com/gobwas/ws/wsutil"
	"github.com/networm6/PoliteCat/common/cache"
	"github.com/networm6/PoliteCat/common/tools"
	"net"
)

// wsToTun Client《--请求结果《--TUN《--inputStream--ws
func wsToTun(conn net.Conn, inputStream chan<- []byte, _cancel context.CancelFunc, _ctx context.Context) {
	for tools.ContextOpened(_ctx) {
		packet, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			break
		}
		inputStream <- packet[:]
	}
	_cancel()
}

// tunToWs Client--网络请求--》TUN--outputStream--》ws
func tunToWs(outputStream <-chan []byte, _ctx context.Context) {
	for tools.ContextOpened(_ctx) {
		bytes := <-outputStream
		if v, ok := cache.GetCache().Get(ConnTag); ok {
			conn := v.(net.Conn)
			if err := wsutil.WriteClientBinary(conn, bytes); err != nil {
				continue
			}
		}
	}
}
