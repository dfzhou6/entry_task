package response

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// RpcRspToHttpRsp rpc响应转换为http响应
func RpcRspToHttpRsp(err error, ctx *gin.Context) {
	if err == nil {
		return
	}

	message := ""
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			message = "账号或密码错误"
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": message,
			})
		case codes.Internal:
			message = "服务器内部错误"
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": message,
			})
		case codes.PermissionDenied:
			message = "账号或密码错误"
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": message,
			})
		case codes.Unauthenticated:
			message = "未授权"
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": message,
			})
		case codes.AlreadyExists:
			message = "用户已存在"
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": message,
			})
		default:
			message = "服务器其他错误"
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": message,
			})
		}
	} else {
		message = "非法请求"
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
	}
}

// BadRequestRsp 非法请求响应
func BadRequestRsp(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "非法请求",
		"error":   err.Error(),
	})
}

// ValidErrorRsp 参数错误响应
func ValidErrorRsp(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "参数解析失败",
		"error":   err.Error(),
	})
}

// ErrorRsp 内部错误响应
func ErrorRsp(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "请求处理失败",
		"error":   err.Error(),
	})
}

// SuccessDataRsp 成功响应(带上数据)
func SuccessDataRsp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// SuccessRsp 成功响应
func SuccessRsp(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
