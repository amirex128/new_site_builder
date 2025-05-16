package cache

type ICacheService interface {
	Set(key string, value interface{})
}
