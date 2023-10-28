package tunnel

import (
	"github.com/networm6/PoliteCat/common/tools"
	"log"
	"net"
	"runtime"
	"strconv"

	"github.com/net-byte/water"
)

// CreateTunnelInterface creates a tunnel interface
func CreateTunnelInterface(config TunConfig) (iFace *water.Interface) {
	CIDR := config.Address.CIDR
	CIDRv6 := config.Address.CIDRv6

	c := water.Config{DeviceType: water.TUN}
	c.PlatformSpecificParams = water.PlatformSpecificParams{}
	os := runtime.GOOS
	if os == "windows" {
		c.PlatformSpecificParams.Name = "vtun"
		c.PlatformSpecificParams.Network = []string{CIDR, CIDRv6}
	}
	if config.DeviceName != "" {
		c.PlatformSpecificParams.Name = config.DeviceName
	}
	iFace, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to create tunnel interface:", err)
	}
	return iFace
}

// SetRoute sets the system routes
func SetRoute(config TunConfig, iFace *water.Interface) {
	CIDR := config.Address.CIDR
	CIDRv6 := config.Address.CIDRv6

	ServerTunnelIP := config.Address.ServerTunnelIP
	ServerTunnelIPv6 := config.Address.ServerTunnelIPv6

	LocalGateway := config.LocalGateway
	LocalGatewayv6 := config.LocalGatewayv6

	MTU := config.MTU
	ServerAddr := config.ServerAddr

	ip, _, err := net.ParseCIDR(CIDR)
	if err != nil {
		log.Panicf("error cidr %v", CIDR)
	}
	ipv6, _, err := net.ParseCIDR(CIDRv6)
	if err != nil {
		log.Panicf("error ipv6 cidr %v", CIDRv6)
	}
	os := runtime.GOOS
	if os == "linux" {
		tools.ExecCmd("/sbin/ip", "link", "set", "dev", iFace.Name(), "mtu", strconv.Itoa(MTU))
		tools.ExecCmd("/sbin/ip", "addr", "add", CIDR, "dev", iFace.Name())
		tools.ExecCmd("/sbin/ip", "-6", "addr", "add", CIDRv6, "dev", iFace.Name())
		tools.ExecCmd("/sbin/ip", "link", "set", "dev", iFace.Name(), "up")
		if !config.ServerMode {
			physicaliFace := tools.GetInterface()
			serverAddrIP := tools.LookupServerAddrIP(ServerAddr)
			if physicaliFace != "" && serverAddrIP != nil {
				if LocalGateway != "" {
					tools.ExecCmd("/sbin/ip", "route", "add", "0.0.0.0/1", "dev", iFace.Name())
					tools.ExecCmd("/sbin/ip", "route", "add", "128.0.0.0/1", "dev", iFace.Name())
					if serverAddrIP.To4() != nil {
						tools.ExecCmd("/sbin/ip", "route", "add", serverAddrIP.To4().String()+"/32", "via", LocalGateway, "dev", physicaliFace)
					}
				}
				if LocalGatewayv6 != "" {
					tools.ExecCmd("/sbin/ip", "-6", "route", "add", "::/1", "dev", iFace.Name())
					if serverAddrIP.To16() != nil {
						tools.ExecCmd("/sbin/ip", "-6", "route", "add", serverAddrIP.To16().String()+"/128", "via", LocalGatewayv6, "dev", physicaliFace)
					}
				}
			}
		}
	} else if os == "darwin" {
		tools.ExecCmd("ifconfig", iFace.Name(), "inet", ip.String(), ServerTunnelIP, "up")
		tools.ExecCmd("ifconfig", iFace.Name(), "inet6", ipv6.String(), ServerTunnelIPv6, "up")
		if !config.ServerMode {
			physicaliFace := tools.GetInterface()
			serverAddrIP := tools.LookupServerAddrIP(ServerAddr)
			if physicaliFace != "" && serverAddrIP != nil {
				if LocalGateway != "" {
					tools.ExecCmd("route", "add", "default", ServerTunnelIP)
					tools.ExecCmd("route", "change", "default", ServerTunnelIP)
					tools.ExecCmd("route", "add", "0.0.0.0/1", "-interface", iFace.Name())
					tools.ExecCmd("route", "add", "128.0.0.0/1", "-interface", iFace.Name())
					if serverAddrIP.To4() != nil {
						tools.ExecCmd("route", "add", serverAddrIP.To4().String(), LocalGateway)
					}
				}
				if LocalGatewayv6 != "" {
					tools.ExecCmd("route", "add", "-inet6", "default", ServerTunnelIPv6)
					tools.ExecCmd("route", "change", "-inet6", "default", ServerTunnelIPv6)
					tools.ExecCmd("route", "add", "-inet6", "::/1", "-interface", iFace.Name())
					if serverAddrIP.To16() != nil {
						tools.ExecCmd("route", "add", "-inet6", serverAddrIP.To16().String(), LocalGatewayv6)
					}
				}
			}
		}
	} else if os == "windows" {
		if !config.ServerMode {
			serverAddrIP := tools.LookupServerAddrIP(ServerAddr)
			if serverAddrIP != nil {
				if LocalGateway != "" {
					tools.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
					tools.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", ServerTunnelIP, "metric", "6")
					if serverAddrIP.To4() != nil {
						tools.ExecCmd("cmd", "/C", "route", "add", serverAddrIP.To4().String()+"/32", config.LocalGateway, "metric", "5")
					}
				}
				if LocalGatewayv6 != "" {
					tools.ExecCmd("cmd", "/C", "route", "-6", "delete", "::/0", "mask", "::/0")
					tools.ExecCmd("cmd", "/C", "route", "-6", "add", "::/0", "mask", "::/0", ServerTunnelIPv6, "metric", "6")
					if serverAddrIP.To16() != nil {
						tools.ExecCmd("cmd", "/C", "route", "-6", "add", serverAddrIP.To16().String()+"/128", LocalGatewayv6, "metric", "5")
					}
				}
			}
		}
	} else {
		log.Printf("not support os %v", os)
	}
	log.Printf("interface configured %v", iFace.Name())
}

// ResetRoute resets the system routes
func ResetRoute(config TunConfig) {
	LocalGateway := config.LocalGateway
	LocalGatewayv6 := config.LocalGatewayv6
	ServerAddr := config.ServerAddr

	if config.ServerMode {
		return
	}

	os := runtime.GOOS
	if os == "darwin" {
		if config.LocalGateway != "" {
			tools.ExecCmd("route", "add", "default", LocalGateway)
			tools.ExecCmd("route", "change", "default", LocalGateway)
		}
		if config.LocalGatewayv6 != "" {
			tools.ExecCmd("route", "add", "-inet6", "default", LocalGatewayv6)
			tools.ExecCmd("route", "change", "-inet6", "default", LocalGatewayv6)
		}
	} else if os == "windows" {
		serverAddrIP := tools.LookupServerAddrIP(ServerAddr)
		if serverAddrIP != nil {
			if LocalGateway != "" {
				tools.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
				tools.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", LocalGateway, "metric", "6")
			}
			if LocalGatewayv6 != "" {
				tools.ExecCmd("cmd", "/C", "route", "-6", "delete", "::/0", "mask", "::/0")
				tools.ExecCmd("cmd", "/C", "route", "-6", "add", "::/0", "mask", "::/0", LocalGatewayv6, "metric", "6")
			}
		}
	}
}
