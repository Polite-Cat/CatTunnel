package server

import (
	"github.com/networm6/CatTunnel/common/data"
	"github.com/networm6/CatTunnel/protocol/dhcp"
	"net/rpc"
)

var _conf dhcp.Config

type DHCP struct{}

func (server DHCP) PickIP(requestKey string, reply *data.Address) {
	if requestKey != _conf.Key {
		*reply = data.Address{}
		return
	}
	key, clientIP := PickIP(_conf.CIDR)
	*reply = data.Address{
		ServerTunnelIP:   _conf.CIDR,
		ServerTunnelIPv6: _conf.CIDRv6,
		CIDR:             clientIP,
		CIDRv6:           "",
		Key:              key,
	}
}

func StartDHCPServer(config dhcp.Config) (*rpc.Server, error) {
	_conf = config
	d := new(DHCP)
	server := rpc.NewServer()
	err := server.RegisterName("PickIP", d)
	if err != nil {
		return nil, err
	}
	return server, nil
}
