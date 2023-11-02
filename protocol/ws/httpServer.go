package ws

import (
	"fmt"
	"github.com/networm6/CatTunnel/protocol/ws/dhcp"
	"github.com/networm6/CatTunnel/tunnel"
	"io"
	"net/http"
	"runtime"
	"strings"
)

func CheckPermission(w http.ResponseWriter, req *http.Request, config *WSConfig) bool {
	key := req.Header.Get("key")
	if key != config.Key {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("No permission"))
		return false
	}
	return true
}

func StartHttpServer(config *WSConfig, tunConfig *tunnel.TunConfig) {
	http.HandleFunc("/ip", func(w http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = strings.Split(req.RemoteAddr, ":")[0]
		}
		resp := fmt.Sprintf("%v", ip)
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	http.HandleFunc("/register/list/ip", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		_, _ = io.WriteString(w, strings.Join(dhcp.ListIP(), "\r\n"))
		runtime.Gosched()
	})

	http.HandleFunc("/register/prefix/ipv4", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		resp := tunConfig.Address.CIDR
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	http.HandleFunc("/register/prefix/ipv6", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		resp := tunConfig.Address.CIDRv6
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})
}
