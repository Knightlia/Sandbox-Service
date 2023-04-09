package repository

import "github.com/jellydator/ttlcache/v3"

type UserRepository struct {
	cache *ttlcache.Cache[string, string]
}

func NewUserRepository(cache *ttlcache.Cache[string, string]) UserRepository {
	return UserRepository{cache}
}

func (u UserRepository) Remove(token string) {
	u.cache.Delete(token)
}

func (u UserRepository) HasValue(value string) bool {
	for _, item := range u.cache.Items() {
		if item.Value() == value {
			return true
		}
	}
	return false
}

func (u UserRepository) Store(key string, value string) {
	u.cache.Set(key, value, ttlcache.DefaultTTL)
}

func (u UserRepository) Values() []string {
	v := make([]string, 0)
	for _, i := range u.cache.Items() {
		v = append(v, i.Value())
	}
	return v
}

func (u UserRepository) Get(key string) string {
	item := u.cache.Get(key)
	if item != nil {
		return item.Value()
	}
	return ""
}
