package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
	"nhooyr.io/websocket"
)

var (
	SessionCache *ttlcache.Cache[*websocket.Conn, string]
	UserCache    *ttlcache.Cache[string, string]
)

func InitCaches() {
	initSessionCache()
	initUserCache()
}

func initSessionCache() {
	SessionCache = ttlcache.New[*websocket.Conn, string](
		// TODO: Update/remove capacity
		ttlcache.WithCapacity[*websocket.Conn, string](1000),
		ttlcache.WithTTL[*websocket.Conn, string](24*time.Hour),
	)
	go SessionCache.Start()
}

func initUserCache() {
	UserCache = ttlcache.New[string, string](
		// TODO: Update/remove capacity
		ttlcache.WithCapacity[string, string](1000),
		ttlcache.WithTTL[string, string](24*time.Hour),
	)
	go UserCache.Start()
}
