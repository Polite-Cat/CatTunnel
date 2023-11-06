package app

// AppConfig The config struct
type AppConfig struct {
	ServerMode bool   `json:"server_mode"`
	ServerAddr string `json:"server_addr"`
	Key        string `json:"key"`

	WSPath    string `json:"path"`
	Timeout   int    `json:"timeout"`
	MixinFunc string `json:"mixin_func"`
}

type nativeConfig AppConfig

var DefaultConfig = nativeConfig{
	ServerMode: false,
	ServerAddr: ":3001",
	Key:        "fuck_key",
	WSPath:     "/freedom",
	Timeout:    30,
	MixinFunc:  "xor",
}
