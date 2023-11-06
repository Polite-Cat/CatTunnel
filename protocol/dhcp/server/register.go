package server

import (
	"crypto/rand"
	"encoding/hex"
	strbox "github.com/networm6/gopherBox/strings"
	"log"
	"net"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// The global cache for register
var _ipPool *cache.Cache

func init() {
	_ipPool = cache.New(30*time.Minute, 3*time.Minute)
}

// generateKey 生成随机key
func generateKey() string {
	number := make([]byte, 32)
	_, _ = rand.Read(number)
	return hex.EncodeToString(number)
}

// addIP 向ip池内加入ip
func addIP(key, ip string) {
	_ipPool.Set(key, ip, cache.DefaultExpiration)
}

// existIP 判断ip是否存在
func existIP(key, ip string) bool {
	if get, found := _ipPool.Get(key); found {
		foundIP := get.(string)
		if foundIP == ip {
			return true
		}
	}
	return false
}

// KeepAliveIP 保持ip连接
func KeepAliveIP(key, ip string) bool {
	if existIP(key, ip) {
		addIP(key, ip)
		return true
	}
	return false
}

// ListIP 获取所有ip
func ListIP() []string {
	var result []string
	for _, value := range _ipPool.Items() {
		result = append(result, value.Object.(string))
	}
	return result
}

// PickIP 生成ip
func PickIP(cidr string) (key string, clientIP string) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}
	total := addressCount(ipNet) - 3
	index := uint64(0)
	//skip first ip
	ip = incr(ipNet.IP.To4())
	ipList := ListIP()
	for {
		ip = incr(ip)
		index++
		if index >= total {
			break
		}
		genIP := ip.String()
		if !strbox.StrIN(genIP, ipList) {
			genKey := generateKey()
			addIP(genKey, genIP)
			return genKey, genIP + "/" + strings.Split(cidr, "/")[1]
		}
	}
	return "", ""
}

// addressCount returns the number of addresses in a CIDR network.
func addressCount(network *net.IPNet) uint64 {
	prefixLen, bits := network.Mask.Size()
	return 1 << (uint64(bits) - uint64(prefixLen))
}

// incr increments the ip by 1
func incr(IP net.IP) net.IP {
	IP = checkIPv4(IP)
	incIP := make([]byte, len(IP))
	copy(incIP, IP)
	for j := len(incIP) - 1; j >= 0; j-- {
		incIP[j]++
		if incIP[j] > 0 {
			break
		}
	}
	return incIP
}

// checkIPv4 checks if the ip is IPv4
func checkIPv4(ip net.IP) net.IP {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip
}
