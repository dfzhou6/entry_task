package route

import (
	"github.com/gin-gonic/gin"
	"user_web/handler"
	"user_web/middleware"
)

// RegisterHandler 注册路由接口
func RegisterHandler(r *gin.Engine) {
	// 注册recovery、全局请求ID中间件
	r.Use(gin.Recovery(), middleware.GenerateGlobalRequestID())

	baseCtrl := new(handler.BaseController)
	r.GET("/captcha", baseCtrl.ShowCaptcha)
	r.POST("/register", baseCtrl.RegisterHandler)
	r.POST("/login", baseCtrl.LoginHandler)
	r.POST("/logout", middleware.CheckTokenExist(), baseCtrl.LogoutHandler)

	// 设置user分组，并注册"检查Authorization参数"中间件
	userGroup := r.Group("user", middleware.CheckTokenExist())
	{
		userCtrl := new(handler.UserController)
		userGroup.GET("/", userCtrl.GetUserProfileHandler)
		userGroup.PUT("/", userCtrl.EditUserProfileHandler)
		userGroup.POST("/avatar", userCtrl.UploadPicHandler)
	}
}
