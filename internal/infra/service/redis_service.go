package service

import sfredis "git.snappfood.ir/backend/go/packages/sf-redis"

type Redis struct {
	conn *sfredis.SfRedis
}

func (r Redis) Set(key string, value interface{}) {

}

func NewRedis(redis *sfredis.SfRedis) *Redis {
	return &Redis{
		conn: redis,
	}
}
