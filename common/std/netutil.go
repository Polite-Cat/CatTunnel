package std

import (
	"log"
	"net"
	"strings"
)

func GetPhysicalInterface() (name string, gateway string, network string) {
	interfaces := getAllPhysicalInterfaces()
	if len(interfaces) == 0 {
		return "", "", ""
	}
	addresses, _ := interfaces[0].Addrs()
	for _, addr := range addresses {
		ip, ok := addr.(*net.IPNet)
		if ok && ip.IP.To4() != nil && !ip.IP.IsLoopback() {
			ipNet := ip.IP.To4().Mask(ip.IP.DefaultMask()).To4()
			network = strings.Join([]string{ipNet.String(), strings.Split(ip.String(), "/")[1]}, "/")
			ipNet[3]++
			gateway = ipNet.String()
			name = interfaces[0].Name
			log.Printf("physical interface %v gateway %v network %v", name, gateway, network)
			break
		}
	}
	return name, gateway, network
}

func getAllPhysicalInterfaces() []net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return nil
	}

	var outInterfaces []net.Interface
	for _, elem := range interfaces {
		if elem.Flags&net.FlagLoopback == 0 && elem.Flags&net.FlagUp == 1 && isPhysicalInterface(elem.Name) {
			addresses, _ := elem.Addrs()
			if len(addresses) > 0 {
				outInterfaces = append(outInterfaces, elem)
			}
		}
	}
	return outInterfaces
}

func isPhysicalInterface(addr string) bool {
	prefixArray := []string{"ens", "enp", "enx", "eno", "eth", "en0", "wlan", "wlp", "wlo", "wlx", "wifi0", "lan0"}
	for _, pref := range prefixArray {
		if strings.HasPrefix(strings.ToLower(addr), pref) {
			return true
		}
	}
	return false
}

func LookupIP(domain string) string {
	ips, err := net.LookupIP(domain)
	if err != nil {
		log.Println(err)
		return ""
	}
	for _, ip := range ips {
		return ip.To4().String()
	}
	return ""
}

func IsIPv4(packet []byte) bool {
	return 4 == (packet[0] >> 4)
}

func IsIPv6(packet []byte) bool {
	return 6 == (packet[0] >> 4)
}

func GetIPv4Source(packet []byte) net.IP {
	return net.IPv4(packet[12], packet[13], packet[14], packet[15])
}

func GetIPv4Destination(packet []byte) net.IP {
	return net.IPv4(packet[16], packet[17], packet[18], packet[19])
}

func GetIPv6Source(packet []byte) net.IP {
	return packet[8:24]
}

func GetIPv6Destination(packet []byte) net.IP {
	return packet[24:40]
}

func GetSourceKey(packet []byte) string {
	key := ""
	if IsIPv4(packet) && len(packet) >= 20 {
		key = GetIPv4Source(packet).To4().String()
	} else if IsIPv6(packet) && len(packet) >= 40 {
		key = GetIPv6Source(packet).To16().String()
	}
	return key
}

func GetDestinationKey(packet []byte) string {
	key := ""
	if IsIPv4(packet) && len(packet) >= 20 {
		key = GetIPv4Destination(packet).To4().String()
	} else if IsIPv6(packet) && len(packet) >= 40 {
		key = GetIPv6Destination(packet).To16().String()
	}
	return key
}
