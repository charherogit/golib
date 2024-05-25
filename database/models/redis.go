package models

import (
	log "golib/logutil"
	"time"

	"github.com/go-redis/redis"
)

var (
	RedisPoolSize     = 10
	RedisMinIdleConns = 2

	ExpireTime1Day     = 24 * time.Hour
	ExpireTime7Day     = 7 * ExpireTime1Day
	ExpireTime15Day    = 15 * ExpireTime1Day
	ExpireTime30Day    = 30 * ExpireTime1Day
	TokenExpireTime    = ExpireTime7Day
	GameDataExpireTime = 14 * 24 * time.Hour // 暂定为 14 天
)

func CreateRedis(addr string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     RedisPoolSize,
		PoolTimeout:  10 * time.Second,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		DB:           db,
		MinIdleConns: RedisMinIdleConns, // 最小连接数
	})
	log.Infof(addr)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedis(opts ...RDBOption) (*redis.Client, error) {
	cfg := &redis.Options{
		PoolSize:     RedisPoolSize,
		PoolTimeout:  10 * time.Second,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		MinIdleConns: RedisMinIdleConns,
	}
	for _, v := range opts {
		v(cfg)
	}
	c := redis.NewClient(cfg)
	if err := c.Ping().Err(); err != nil {
		return nil, err
	}
	return c, nil
}
