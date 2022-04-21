package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"user_web/pkg/logger"
	"user_web/pkg/util"
)

// GenerateGlobalRequestID 根据时间戳生成请求ID
func GenerateGlobalRequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		randStr := fmt.Sprintf("%d-%s", time.Now().UnixNano(), util.GenerateRandomStr())
		requestID := util.Md5Data(randStr)
		ctx.Set("request_id", requestID)
		ctx.Next()
	}
}

// CheckTokenExist 检查Authorization参数
func CheckTokenExist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if len(token) == 0 {
			logger.Warn("middleware", "token not exist", token)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token 必传",
			})
			return
		}
		ctx.Next()
	}
}
