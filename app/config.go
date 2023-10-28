package app

import (
	"encoding/json"
	"os"
)

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

func (c *AppConfig) UnmarshalJSON(data []byte) error {
	_ = json.Unmarshal(data, &DefaultConfig)
	*c = AppConfig(DefaultConfig)
	return nil
}

func (c *AppConfig) LoadConfig(configFile string) (err error) {
	file, err := os.Open(configFile)
	if err != nil {
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return
	}
	return
}
