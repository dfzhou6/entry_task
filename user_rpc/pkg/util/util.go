package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"time"
	"user_rpc/pkg/config"
)

// Md5Data md5加密
func Md5Data(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	return hex.EncodeToString(md.Sum(nil))
}

// GenerateToken 根据用户名生成token
func GenerateToken(username string) string {
	key := fmt.Sprintf("%s-%d-%s", username, time.Now().UnixNano(), GenerateRandomStr())
	return Md5Data(key)
}

// GeneratePwdHash 根据密码和盐值生成哈希密码
func GeneratePwdHash(password string, salt string) string {
	key := fmt.Sprintf("%s-%s", password, salt)
	return Md5Data(key)
}

// ComparePwdHash 对比密码是否正确
func ComparePwdHash(password string, passwordHash string, salt string) bool {
	return GeneratePwdHash(password, salt) == passwordHash
}

// GenerateRandomStr 生成随机字符串
func GenerateRandomStr() string {
	rand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 14)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetRequestIDFromContext 从context中获取请求ID
func GetRequestIDFromContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	return md["request_id"][0]
}

// GetTableByUsername 根据用户名获取表名称
func GetTableByUsername(username string) string {
	var cSum int32
	for _, c := range username {
		cSum += c
	}
	return fmt.Sprintf("users_%d", cSum%config.GetInt32("DB_USER_TABLE_COUNT"))
}
