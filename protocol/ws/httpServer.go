package ws

import (
	"fmt"
	"github.com/networm6/PoliteCat/protocol/ws/register"
	"github.com/networm6/PoliteCat/tunnel"
	"io"
	"net"
	"net/http"
	"runtime"
	"strings"
)

func CheckPermission(w http.ResponseWriter, req *http.Request, config *WSConfig) bool {
	key := req.Header.Get("key")
	if key != config.Key {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("No permission"))
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

	http.HandleFunc("/register/pick/ip", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		ip, pl := register.PickClientIP(tunConfig.Address.CIDR)
		resp := fmt.Sprintf("%v/%v", ip, pl)
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	http.HandleFunc("/register/delete/ip", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		ip := r.URL.Query().Get("ip")
		if ip != "" {
			register.DeleteClientIP(ip)
		}
		_, _ = io.WriteString(w, "OK")
		runtime.Gosched()
	})

	http.HandleFunc("/register/keepalive/ip", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		ip := r.URL.Query().Get("ip")
		if ip != "" {
			register.KeepAliveClientIP(ip)
		}
		_, _ = io.WriteString(w, "OK")
		runtime.Gosched()
	})

	http.HandleFunc("/register/list/ip", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		_, _ = io.WriteString(w, strings.Join(register.ListClientIPs(), "\r\n"))
		runtime.Gosched()
	})

	http.HandleFunc("/register/prefix/ipv4", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		_, ipv4Net, err := net.ParseCIDR(tunConfig.Address.CIDR)
		var resp string
		if err != nil {
			resp = "error"
		} else {
			resp = ipv4Net.String()
		}
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

	http.HandleFunc("/register/prefix/ipv6", func(w http.ResponseWriter, r *http.Request) {
		if !CheckPermission(w, r, config) {
			return
		}
		_, ipv6Net, err := net.ParseCIDR(tunConfig.Address.CIDRv6)
		var resp string
		if err != nil {
			resp = "error"
		} else {
			resp = ipv6Net.String()
		}
		_, _ = io.WriteString(w, resp)
		runtime.Gosched()
	})

}
