package captcha

import (
	"fmt"
	"time"
	"user_web/pkg/config"
	"user_web/pkg/logger"
	"user_web/pkg/redis"
)

// StoreRedis 验证码的redis缓存实现
type StoreRedis struct {
	Client *redis.RedisClient
}

// Set 保存验证码到缓存
func (sr *StoreRedis) Set(key string, val string) error {
	newKey := fmt.Sprintf("%s:%s", config.GetString("CAPTCHA_PREFIX"), key)
	expire := time.Duration(config.GetInt("CAPTCHA_KEY_EXPIRE")) * time.Second
	return sr.Client.Set(newKey, val, expire)
}

// Get 从缓存中获取验证码
func (sr *StoreRedis) Get(key string, clear bool) string {
	newKey := fmt.Sprintf("%s:%s", config.GetString("CAPTCHA_PREFIX"), key)
	val, err := sr.Client.Get(newKey)
	if err != nil {
		logger.Error("store_redis", "get key error, key:", newKey)
		return ""
	}
	if clear {
		if err := sr.Client.Del(newKey); err != nil {
			logger.Error("store_redis", "del key error, key:", newKey)
			return ""
		}
	}
	return val
}

// Verify 验证验证码是否正确
func (sr *StoreRedis) Verify(key, ans string, clear bool) bool {
	return sr.Get(key, clear) == ans
}
