package service

import (
	"encoding/json"
	"log"

	"github.com/networm6/PoliteCat/common/cipher"
)

type Config struct {
	LocalAddr  string
	ServerAddr string
	CIDR       string // TUN网卡的IP
	Key        string // 密钥
	DNS        string // TUN网卡的DNS
	ServerMode bool   // 标记位，是否是服务端
	GlobalMode bool   // 是否全局流量转发
	Enc        bool   // 是否开启加密
	MTU        int    // 一个网络包最大的值
	Timeout    int    // 每个TCP的超时时间
}

func (config *Config) Init() {
	cipher.GenerateKey(config.Key)
	json, _ := json.Marshal(config)
	log.Printf("init config:%s", string(json))
}
