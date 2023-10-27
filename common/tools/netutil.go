package tools

import (
	"encoding/json"
	"github.com/net-byte/go-gateway"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

type Address struct {
	ServerTunnelIP   string `json:"server_ip"`
	ServerTunnelIPv6 string `json:"server_ipv6"`
	CIDR             string `json:"cidr"`
	CIDRv6           string `json:"cidr_ipv6"`
}

func InitAddress(url string, key string) (*Address, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	req.Header.Add("key", key)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ddr Address

	err = json.Unmarshal(b, &ddr)
	if err != nil {
		return nil, err
	}
	return &ddr, nil
}

// GetInterface returns the name of interface
func GetInterface() (name string) {
	ifaces := getAllInterfaces()
	if len(ifaces) == 0 {
		return ""
	}
	netAddrs, _ := ifaces[0].Addrs()
	for _, addr := range netAddrs {
		ip, ok := addr.(*net.IPNet)
		if ok && ip.IP.To4() != nil && !ip.IP.IsLoopback() {
			name = ifaces[0].Name
			break
		}
	}
	return name
}

// getAllInterfaces returns all interfaces
func getAllInterfaces() []net.Interface {
	iFaceList, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return nil
	}

	var outInterfaces []net.Interface
	for _, iFace := range iFaceList {
		if iFace.Flags&net.FlagLoopback == 0 && iFace.Flags&net.FlagUp == 1 && isPhysicalInterface(iFace.Name) {
			netAddrList, _ := iFace.Addrs()
			if len(netAddrList) > 0 {
				outInterfaces = append(outInterfaces, iFace)
			}
		}
	}
	return outInterfaces
}

// isPhysicalInterface returns true if the interface is physical
func isPhysicalInterface(addr string) bool {
	prefixArray := []string{"ens", "enp", "enx", "eno", "eth", "en0", "wlan", "wlp", "wlo", "wlx", "wifi0", "lan0"}
	for _, pref := range prefixArray {
		if strings.HasPrefix(strings.ToLower(addr), pref) {
			return true
		}
	}
	return false
}

// LookupIP Lookup IP address of the given hostname
func LookupIP(domain string) net.IP {
	ips, err := net.LookupIP(domain)
	if err != nil || len(ips) == 0 {
		log.Println(err)
		return nil
	}
	return ips[0]
}

// ExecCmd executes the given command
func ExecCmd(c string, args ...string) string {
	//log.Printf("exec %v %v", c, args)
	cmd := exec.Command(c, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Println("failed to exec cmd:", err)
	}
	if len(out) == 0 {
		return ""
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", "")
}

// DiscoverGateway returns the local gateway IP address
func DiscoverGateway(ipv4 bool) string {
	var ip net.IP
	var err error
	if ipv4 {
		ip, err = gateway.DiscoverGatewayIPv4()
	} else {
		ip, err = gateway.DiscoverGatewayIPv6()
	}
	if err != nil {
		log.Println(err)
		return ""
	}
	return ip.String()
}

// LookupServerAddrIP returns the IP of server address
func LookupServerAddrIP(serverAddr string) net.IP {
	host, _, err := net.SplitHostPort(serverAddr)
	if err != nil {
		log.Panic("error server address")
		return nil
	}
	ip := LookupIP(host)
	return ip
}
