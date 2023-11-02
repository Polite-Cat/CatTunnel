package server

import (
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/networm6/CatTunnel/common/cache"
	"github.com/networm6/CatTunnel/common/tools"
	ws2 "github.com/networm6/CatTunnel/protocol/ws"
	"github.com/networm6/CatTunnel/tunnel"
	"log"
	"net"
	"net/http"
	"time"
)

func StartServer(conf *ws2.WSConfig, tun *tunnel.Tunnel) {
	go tunToWs(tun.OutputStream, *tun.LifeCtx)
	http.HandleFunc(conf.WSPath, func(w http.ResponseWriter, r *http.Request) {
		if !ws2.CheckPermission(w, r, conf) {
			return
		}
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}

		wsToTun(conn, tun.InputStream, conf.Timeout)
	})

	err := http.ListenAndServe(conf.ServerAddr, nil)
	if err != nil {
		log.Fatalf("server error %v", err)
		return
	}
}

// tun --> outputStream --> ws
func tunToWs(outputStream <-chan []byte, _ctx context.Context) {
	for tools.ContextOpened(_ctx) {
		bytes := <-outputStream
		if clientIP := tools.GetDstKey(bytes); clientIP != "" {
			if conn, ok := cache.GetCache().Get(clientIP); ok {
				err := wsutil.WriteServerBinary(conn.(net.Conn), bytes)
				if err != nil {
					cache.GetCache().Delete(clientIP)
					continue
				}
			}
		}
	}
}

// tun <-- inputStream <-- ws
func wsToTun(wsconn net.Conn, inputStream chan<- []byte, timeout int) {
	defer wsconn.Close()
	for {
		_ = wsconn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		bytes, op, err := wsutil.ReadClientData(wsconn)
		if err != nil {
			break
		}
		if op == ws.OpText {
			if clientIP := tools.GetSrcKey(bytes); clientIP != "" {
				log.Printf("recv ping from %s %s\n", wsconn.RemoteAddr(), clientIP)
			}
			_ = wsutil.WriteServerMessage(wsconn, op, bytes)
		} else if op == ws.OpBinary {
			if clientIP := tools.GetSrcKey(bytes); clientIP != "" {
				cache.GetCache().Set(clientIP, wsconn, 24*time.Hour)
				inputStream <- bytes
			}
		}
	}
}
