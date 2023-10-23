package tcp

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(30*time.Minute, 10*time.Minute)

const Key = "simon_tcp"

func getCache() *cache.Cache {
	return c
}
