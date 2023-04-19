package repository

import (
	"github.com/jellydator/ttlcache/v3"
	"nhooyr.io/websocket"
)

type SessionRepository struct {
	cache *ttlcache.Cache[*websocket.Conn, string]
}

func NewSessionRepository(cache *ttlcache.Cache[*websocket.Conn, string]) SessionRepository {
	return SessionRepository{cache}
}

func (s SessionRepository) HasValue(token string) bool {
	for _, k := range s.cache.Items() {
		if k.Value() == token {
			return true
		}
	}
	return false
}

func (s SessionRepository) Store(key *websocket.Conn, value string) {
	s.cache.Set(key, value, ttlcache.DefaultTTL)
}

func (s SessionRepository) Remove(key *websocket.Conn) string {
	item := s.cache.Get(key)
	s.cache.Delete(key)
	return item.Value()
}
