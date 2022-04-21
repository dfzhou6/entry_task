package handler

import (
	"github.com/gin-gonic/gin"
	"user_web/client"
	"user_web/pkg/captcha"
	"user_web/pkg/logger"
	"user_web/pkg/util"
	"user_web/proto"
	"user_web/request"
	"user_web/response"
)

// BaseController 登录前控制器
type BaseController struct {
}

// RegisterHandler 注册接口
func (ctrl *BaseController) RegisterHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")

	// 验证请求参数
	req := request.RegisterRequest{}
	if err := request.Validate(ctx, &req, request.RegisterRequestValid); err != nil {
		logger.Warn("register", "requestID:", requestID, "request valid error", err.Error())
		return
	}

	// rpc调用
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	rsp, err := client.RpcClient.CreateUserProfile(rpcCtx, &proto.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		logger.Error("register", "requestID:", requestID, "call CreateUserProfile error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 返回结果
	data := gin.H{
		"id":       rsp.Id,
		"username": rsp.Username,
		"nickname": rsp.Nickname,
		"pic_path": rsp.PicPath,
	}
	logger.Debug("register", "requestID:", requestID, "register success", data)

	response.SuccessDataRsp(ctx, data)
}

// LoginHandler 登录接口
func (ctrl *BaseController) LoginHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")

	// 验证请求参数
	req := request.LoginRequest{}
	if err := request.Validate(ctx, &req, request.LoginRequestValid); err != nil {
		logger.Warn("login", "requestID:", requestID, "request valid error", err.Error())
		return
	}

	// rpc调用
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	rsp, err := client.RpcClient.Login(rpcCtx, &proto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("login", "requestID:", requestID, "call login error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 返回结果
	data := gin.H{
		"id":       rsp.Id,
		"username": rsp.Username,
		"nickname": rsp.Nickname,
		"pic_path": rsp.PicPath,
		"token":    rsp.Token,
	}
	logger.Debug("login", "requestID:", requestID, "login success", data)

	response.SuccessDataRsp(ctx, data)
}

// LogoutHandler 登出接口
func (ctrl *BaseController) LogoutHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	token := ctx.GetHeader("Authorization")

	// rpc调用
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	_, err := client.RpcClient.Logout(rpcCtx, &proto.AuthRequest{
		Token: token,
	})
	if err != nil {
		logger.Error("logout", "requestID:", requestID,
			"token:", token, "call logout error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 返回结果
	logger.Debug("logout", "requestID:", requestID, "token:", token, "logout success")

	response.SuccessRsp(ctx)
}

// ShowCaptcha 获取验证码接口
func (ctrl *BaseController) ShowCaptcha(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")

	// 生成验证码
	id, b64s, err := captcha.NewCaptcha().GenerateCaptcha()
	if err != nil {
		logger.Error("showCaptcha", "requestID:", requestID, "generateCaptcha error", err.Error())
		response.ErrorRsp(ctx, err)
		return
	}

	// 返回结果
	data := gin.H{
		"captcha_id":  id,
		"captcha_img": b64s,
	}
	logger.Debug("showCaptcha", "requestID:", requestID, "generateCaptcha success", data)

	response.SuccessDataRsp(ctx, data)
}
