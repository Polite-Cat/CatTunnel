package main

import (
	"flag"
	"github.com/networm6/PoliteCat/app"
)

var cfg = app.Config{}

func init() {
	flag.StringVar(&cfg.ServerAddr, "s", app.DefaultConfig.ServerAddr, "server address")
	flag.StringVar(&cfg.Key, "k", app.DefaultConfig.Key, "key")

	flag.StringVar(&cfg.WSPath, "path", app.DefaultConfig.WSPath, "ws path")
	flag.IntVar(&cfg.Timeout, "t", app.DefaultConfig.Timeout, "dial timeout in seconds")
	flag.StringVar(&cfg.MixinFunc, "f", app.DefaultConfig.MixinFunc, "mixin function xor/none")

	flag.Parse()
}

func main() {
	runClient()
}
