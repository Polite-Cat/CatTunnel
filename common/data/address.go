package data

type Address struct {
	ServerTunnelIP   string `json:"server_ip"`
	ServerTunnelIPv6 string `json:"server_ipv6"`
	CIDR             string `json:"cidr"`
	CIDRv6           string `json:"cidr_ipv6"`
	Key              string
}
