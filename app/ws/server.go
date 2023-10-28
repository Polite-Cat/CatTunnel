package ws

import (
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/networm6/PoliteCat/common/cache"
	"github.com/networm6/PoliteCat/common/tools"
	"github.com/networm6/PoliteCat/tunnel"
	"log"
	"net"
	"net/http"
	"time"
)

func StartServer(conf *WSConfig, tun *tunnel.Tunnel) {
	go serverTunToWs(tun.OutputStream, *tun.LifeCtx)
	http.HandleFunc(conf.WSPath, func(w http.ResponseWriter, r *http.Request) {
		if !checkPermission(w, r, conf) {
			return
		}
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}

		serverWsToTun(conn, tun.InputStream, conf.Timeout)
	})

	err := http.ListenAndServe(conf.ServerAddr, nil)
	if err != nil {
		log.Fatalf("server error %v", err)
		return
	}
}

func checkPermission(w http.ResponseWriter, req *http.Request, config *WSConfig) bool {
	key := req.Header.Get("key")
	if key != config.Key {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("No permission"))
		return false
	}
	return true
}

// tun --> outputStream --> ws
func serverTunToWs(outputStream <-chan []byte, _ctx context.Context) {
	for tools.ContextOpened(_ctx) {
		bytes := <-outputStream
		if key := tools.GetDstKey(bytes); key != "" {
			if conn, ok := cache.GetCache().Get(key); ok {
				err := wsutil.WriteServerBinary(conn.(net.Conn), bytes)
				if err != nil {
					cache.GetCache().Delete(key)
					continue
				}
			}
		}
	}
}

// tun <-- inputStream <-- ws
func serverWsToTun(wsconn net.Conn, inputStream chan<- []byte, timeout int) {
	defer wsconn.Close()
	for {
		_ = wsconn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		bytes, op, err := wsutil.ReadClientData(wsconn)
		if err != nil {
			break
		}
		if op == ws.OpText {
			_ = wsutil.WriteServerMessage(wsconn, op, bytes)
		} else if op == ws.OpBinary {
			if key := tools.GetSrcKey(bytes); key != "" {
				cache.GetCache().Set(key, wsconn, 24*time.Hour)
				inputStream <- bytes
			}
		}
	}
}
