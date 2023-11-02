package tunnel

import (
	"github.com/networm6/CatTunnel/common/encrypt"
	"github.com/networm6/CatTunnel/common/tools"
)

type TunConfig struct {
	DeviceName string `json:"device_name"`
	MTU        int    `json:"mtu"`

	ServerMode bool `json:"server_mode"`

	ServerAddr     string        `json:"server_addr"`
	Address        tools.Address `json:"address"`
	LocalGateway   string        `json:"local_gateway"`
	LocalGatewayv6 string        `json:"local_gateway_ipv6"`

	BufferSize int    `json:"buffer_size"`
	MixinFunc  string `json:"mixin_func"`
}

func getFunc(mixin string) func([]byte) []byte {
	switch mixin {
	case "none":
		return encrypt.None
	case "reverse":
		return encrypt.Reverse
	case "xor":
	default:
		return encrypt.Xor
	}
	return encrypt.None
}
