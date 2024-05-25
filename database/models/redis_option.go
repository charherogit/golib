package models

import (
	"github.com/go-redis/redis"
)

type RDBOption func(opt *redis.Options)

func WithOptions(f func(opt *redis.Options)) RDBOption {
	return func(opt *redis.Options) {
		f(opt)
	}
}

func Use(addr string, db int) RDBOption {
	return func(opt *redis.Options) {
		opt.Addr, opt.DB = addr, db
	}
}
