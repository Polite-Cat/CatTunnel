package dhcp

import "github.com/networm6/CatTunnel/common/tools"

var _conf Config

type DHCP struct{}

func (server DHCP) PickIP(requestKey string, reply *tools.Address) {
	if requestKey != _conf.Key {
		*reply = tools.Address{}
		return
	}
	key, clientIP := PickIP(_conf.CIDR)
	*reply = tools.Address{
		ServerTunnelIP:   _conf.CIDR,
		ServerTunnelIPv6: _conf.CIDRv6,
		CIDR:             clientIP,
		CIDRv6:           "",
		Key:              key,
	}
}

func StartDHCPServer(config Config) {
	_conf = config

}
