package dhcp

type Config struct {
	CIDR   string `json:"cidr"`
	CIDRv6 string `json:"cidr_ipv6"`
	Key    string `json:"key"`
}
