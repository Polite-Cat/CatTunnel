package client

import (
	"context"
	ws2 "github.com/networm6/CatTunnel/protocol/ws"
	"github.com/networm6/gopherBox/ctxbox"
	"github.com/networm6/gopherBox/tunnel"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/networm6/CatTunnel/common/cache"
)

const ConnTag = "conn"

/*
Client发出的所有网络请求包都会走tun网卡

在StartClient中，Client被抽象为一个双向流，outputStream是用户的请求，inputStream是请求的结果。
在mapStreamsToWebSocket中，用这两股流与ws交互。
在tunToWs中，不断的从outputStream中读取数据，并检测是否存在连接，如果存在则发送到ws。
在wsToTun中，不断的从ws中读取数据，并发送到inputStream。

*/

// UserApp --> Kernel --> UserApp(TUN) --> ReadFromTun --> tunToWs --> ws
func mapStreamsToServer(config *ws2.WSConfig, outputStream <-chan []byte, inputStream chan<- []byte, tunCtx context.Context) {
	go tunToWs(outputStream, tunCtx)
	for ctxbox.Opened(tunCtx) {
		log.Println("new ws")
		// 为每个ws链接建立新的ctx
		connCtx, connCancel := context.WithCancel(tunCtx)
		conn := connectServer(config, connCtx)
		if conn == nil {
			connCancel()
			time.Sleep(3 * time.Second)
			continue
		}
		// 设置一个链接的有效时长为24小时
		cache.GetCache().Set(ConnTag, conn, 24*time.Hour)
		go wsToTun(conn, inputStream, connCancel, connCtx)
		// 建立连接后，每3秒发送一次ping，检测是否断开。
		alive(conn, connCtx, connCancel)
		cache.GetCache().Delete(ConnTag)
		_ = conn.Close()
	}
}

// StartClient 启动Client端。
func StartClient(conf *ws2.WSConfig, tun *tunnel.Tunnel) {
	mapStreamsToServer(conf, tun.OutputStream, tun.InputStream, *tun.LifeCtx)
}

func alive(conn net.Conn, _ctx context.Context, _cancel context.CancelFunc) {
	ticker := time.NewTicker(5 * time.Second)
	for ctxbox.Opened(_ctx) {
		err := wsutil.WriteClientMessage(conn, ws.OpText, []byte("ping"))
		log.Printf("send ping to %s\n", conn.RemoteAddr())
		if err != nil {
			log.Printf("alive error %v\n", err)
			break
		}
		<-ticker.C
	}
	ticker.Stop()
	_cancel()
}

// ConnectServer connects to the server with the given address.
func connectServer(config *ws2.WSConfig, tunCtx context.Context) net.Conn {
	scheme := "ws"

	u := url.URL{Scheme: scheme, Host: config.ServerAddr, Path: config.WSPath}
	header := make(http.Header)
	header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	if config.Key != "" {
		header.Set("key", config.Key)
	}
	dialer := ws.Dialer{
		Header:  ws.HandshakeHeaderHTTP(header),
		Timeout: time.Duration(config.Timeout) * time.Second,
		NetDial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, config.ServerAddr)
		},
	}
	c, _, _, err := dialer.Dial(tunCtx, u.String())
	if err != nil {
		log.Printf("[client] failed to dial websocket %s %v", u.String(), err)
		return nil
	}
	return c
}

// wsToTun Client《--请求结果《--TUN《--inputStream--ws
func wsToTun(conn net.Conn, inputStream chan<- []byte, _cancel context.CancelFunc, _ctx context.Context) {
	for ctxbox.Opened(_ctx) {
		packet, err := wsutil.ReadServerBinary(conn)
		if err != nil {
			log.Printf("packet error %v\n", err)
			break
		}
		inputStream <- packet[:]
	}
	_cancel()
}

// tunToWs Client--网络请求--》TUN--outputStream--》ws
func tunToWs(outputStream <-chan []byte, _ctx context.Context) {
	for ctxbox.Opened(_ctx) {
		bytes := <-outputStream
		if v, ok := cache.GetCache().Get(ConnTag); ok {
			conn := v.(net.Conn)
			if err := wsutil.WriteClientBinary(conn, bytes); err != nil {
				log.Printf("tunToWs error %v\n", err)
				continue
			}
		}
	}
}
