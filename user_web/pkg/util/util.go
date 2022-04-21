package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"user_web/pkg/config"
)

// Md5Data md5加密
func Md5Data(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	return hex.EncodeToString(md.Sum(nil))
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

// TimeNowInTimeZone 获取时间区
func TimeNowInTimeZone() time.Time {
	timeZone, _ := time.LoadLocation(config.GetString("APP_TIMEZONE"))
	return time.Now().In(timeZone)
}

// 生成随机文件名
func randomFileName(file *multipart.FileHeader) string {
	return fmt.Sprintf("%s%s", GenerateRandomStr(), filepath.Ext(file.Filename))
}

// SaveUploadPic 保存上传的图片
func SaveUploadPic(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	publicPath := config.GetString("APP_UPLOAD_DIR")
	avatarDir := fmt.Sprintf("/uploads/avatars/%s", TimeNowInTimeZone().Format("2006/01/02"))
	uploadDir := fmt.Sprintf("%s%s", publicPath, avatarDir)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	avatar := fmt.Sprintf("%s/%s", avatarDir, randomFileName(file))
	dstPath := fmt.Sprintf("%s%s", publicPath, avatar)
	if err := ctx.SaveUploadedFile(file, dstPath); err != nil {
		return "", err
	}
	return avatar, nil
}

// WrapUploadPic 拼接图片路径
func WrapUploadPic(avatar string) string {
	return fmt.Sprintf("%s%s", config.GetString("APP_UPLOAD_DIR"), avatar)
}

// GenRpcCtxWithRequestID 生成带请求ID的context
func GenRpcCtxWithRequestID(ctx *gin.Context) context.Context {
	requestID := ctx.GetString("request_id")
	return metadata.NewOutgoingContext(context.Background(), metadata.Pairs("request_id", requestID))
}
