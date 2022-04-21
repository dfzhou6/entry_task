package redis

import (
	"context"
	"fmt"
	redisLib "github.com/go-redis/redis/v8"
	"time"
	"user_rpc/pkg/config"
	"user_rpc/pkg/logger"
)

type RedisClient struct {
	Client *redisLib.Client
	Ctx    context.Context
}

var Redis *RedisClient

// SetupRedis 初始化redis
func SetupRedis() {
	Redis = &RedisClient{
		Client: redisLib.NewClient(&redisLib.Options{
			Addr: fmt.Sprintf("%s:%d", config.GetString("REDIS_HOST"),
				config.GetInt("REDIS_PORT")),
			Password: config.GetString("REDIS_PASSWORD"),
			DB:       config.GetInt("REDIS_DB"),
		}),
		Ctx: context.Background(),
	}

	if err := Redis.Ping(); err != nil {
		logger.Error("redis", "ping redis error", err)
		panic(err)
	}

	logger.Debug("redis", "conn success")
}

func (rs *RedisClient) Ping() error {
	return rs.Client.Ping(rs.Ctx).Err()
}

func (rs *RedisClient) Set(key string, value interface{}, expire time.Duration) error {
	return rs.Client.Set(rs.Ctx, key, value, expire).Err()
}

func (rs *RedisClient) Get(key string) (string, error) {
	val, err := rs.Client.Get(rs.Ctx, key).Result()
	if err != nil && err == redisLib.Nil {
		return "", nil
	}
	return val, err
}

func (rs *RedisClient) Del(key string) error {
	return rs.Client.Del(rs.Ctx, key).Err()
}

func Close() {
	if Redis == nil {
		return
	}
	if err := Redis.Client.Close(); err != nil {
		logger.Error("redis", "close redis conn error", err)
		panic(err)
	}
	logger.Debug("redis", "close success")
}
