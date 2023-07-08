package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type UserCache struct {
	cache *ttlcache.Cache[string, string]
}

func NewUserCache() UserCache {
	cache := ttlcache.New(
		// TODO: Update capacity
		ttlcache.WithCapacity[string, string](1000),
		ttlcache.WithTTL[string, string](24*time.Hour),
	)

	return UserCache{cache}
}

func (u UserCache) HasValue(value string) bool {
	for _, v := range u.cache.Items() {
		if v.Value() == value {
			return true
		}
	}

	return false
}

func (u UserCache) Store(k string, v string) {
	u.cache.Set(k, v, ttlcache.DefaultTTL)
}

func (u UserCache) Values() []string {
	userList := make([]string, 0, u.cache.Len())

	for _, v := range u.cache.Items() {
		if v.Value() == "" {
			continue
		}
		userList = append(userList, v.Value())
	}

	return userList
}

func (u UserCache) Remove(k string) {
	u.cache.Delete(k)
}

func (u UserCache) HasKey(k string) bool {
	return u.cache.Get(k) != nil
}

func (u UserCache) Get(k string) string {
	if i := u.cache.Get(k); i != nil {
		return i.Value()
	}
	return ""
}
