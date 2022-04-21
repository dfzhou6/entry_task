package user

import (
	"encoding/json"
	"fmt"
	"time"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/config"
	"user_rpc/pkg/redis"
)

// SetTokenCache 写入token缓存
func SetTokenCache(token string, value string) error {
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_TOKEN_PREFIX"), token)
	expire := time.Duration(config.GetInt("CACHE_TOKEN_EXPIRE")) * time.Second
	return redis.Redis.Set(key, value, expire)
}

// GetTokenCache 获取token缓存
func GetTokenCache(token string) (string, error) {
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_TOKEN_PREFIX"), token)
	val, err := redis.Redis.Get(key)
	return val, err
}

// DelTokenCache 删除token缓存
func DelTokenCache(token string) (err error) {
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_TOKEN_PREFIX"), token)
	err = redis.Redis.Del(key)
	return
}

// GetUserCache 获取用户信息缓存
func GetUserCache(username string) (userModel.User, error) {
	var user userModel.User
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_USER_PREFIX"), username)
	val, err := redis.Redis.Get(key)
	if err != nil || len(val) == 0 {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

// SetUserCache 写入用户信息缓存
func SetUserCache(username string, user userModel.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_USER_PREFIX"), username)
	expire := time.Duration(config.GetInt("CACHE_USER_EXPIRE")) * time.Second
	return redis.Redis.Set(key, string(data), expire)
}

// DelUserCache 删除用户信息缓存
func DelUserCache(username string) (err error) {
	key := fmt.Sprintf("%s:%s", config.GetString("CACHE_USER_PREFIX"), username)
	err = redis.Redis.Del(key)
	return
}
