package sfredis

import (
	"time"
)

// Options represents SfRedis connection options
type Options struct {
	PoolSize       int
	MinIdleConns   int
	MaxRetries     int
	DialTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxConnAge     time.Duration
	PoolTimeout    time.Duration
	ClusterOptions struct {
		MaxRedirects   int
		RouteRandomly  bool
		RouteByLatency bool
	}
}

// Z represents a member and its score in a sorted set
type Z struct {
	Score  float64
	Member interface{}
}

// ZRangeBy represents a range query for sorted sets
type ZRangeBy struct {
	Min    string
	Max    string
	Offset int64
	Count  int64
}
