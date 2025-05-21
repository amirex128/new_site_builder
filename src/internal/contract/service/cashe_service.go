package service

type ICacheService interface {
	Set(key string, value interface{})
}
