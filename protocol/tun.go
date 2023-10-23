package protocol

import (
	"github.com/networm6/PoliteCat/daemon/service"
	"log"
	"strconv"
	"strings"

	"github.com/networm6/PoliteCat/common/std"
	"github.com/songgao/water"
)

func CreateTun(config service.Config) *water.Interface {
	c := water.Config{DeviceType: water.TUN}
	device, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to create tun interface:", err)
	}
	log.Println("interface created:", device.Name())
	configTun(config, device)
	return device
}

func configTun(config service.Config, device *water.Interface) {
	std.ExecCmd("/sbin/ip", "link", "set", "dev", device.Name(), "mtu", strconv.Itoa(config.MTU))
	std.ExecCmd("/sbin/ip", "addr", "add", config.CIDR, "dev", device.Name())
	std.ExecCmd("/sbin/ip", "link", "set", "dev", device.Name(), "up")
	if config.GlobalMode {
		physicalInterface, gateway, _ := std.GetPhysicalInterface()
		serverIP := std.LookupIP(strings.Split(config.ServerAddr, ":")[0])
		if physicalInterface != "" && serverIP != "" {
			std.ExecCmd("/sbin/ip", "route", "add", "0.0.0.0/1", "dev", device.Name())
			std.ExecCmd("/sbin/ip", "route", "add", "128.0.0.0/1", "dev", device.Name())
			std.ExecCmd("/sbin/ip", "route", "add", strings.Join([]string{serverIP, "32"}, "/"), "via", gateway, "dev", physicalInterface)
			std.ExecCmd("/sbin/ip", "route", "add", strings.Join([]string{strings.Split(config.DNS, ":")[0], "32"}, "/"), "via", gateway, "dev", physicalInterface)
		}
	}
}
