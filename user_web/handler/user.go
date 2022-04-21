package handler

import (
	"github.com/gin-gonic/gin"
	"user_web/client"
	"user_web/pkg/logger"
	"user_web/pkg/util"
	"user_web/proto"
	"user_web/request"
	"user_web/response"
)

// UserController 用户控制器
type UserController struct {
}

// GetUserProfileHandler 获取用户信息接口
func (ctrl *UserController) GetUserProfileHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	token := ctx.GetHeader("Authorization")

	// rpc调用
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	rsp, err := client.RpcClient.GetUserProfile(rpcCtx, &proto.AuthRequest{
		Token: token,
	})
	if err != nil {
		logger.Error("getUserProfile", "requestID:", requestID,
			"token:", token, "call getUserProfile error", err.Error())
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
	logger.Debug("getUserProfile", "requestID:", requestID, "data:", data)

	response.SuccessDataRsp(ctx, data)
}

// EditUserProfileHandler 编辑用户信息接口
func (ctrl *UserController) EditUserProfileHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	token := ctx.GetHeader("Authorization")

	// 验证请求参数
	req := request.EditUserProfileRequest{}
	if err := request.Validate(ctx, &req, request.EditUserProfileRequestValid); err != nil {
		logger.Warn("editUserProfile", "requestID:", requestID, "request valid error", err.Error())
		return
	}

	// rpc调用
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	rsp, err := client.RpcClient.EditUserProfile(rpcCtx, &proto.EditUserRequest{
		Token:    token,
		Nickname: req.Nickname,
	})
	if err != nil {
		logger.Error("editUserProfile", "requestID:", requestID,
			"token:", token, "call editUserProfile error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 返回结果
	data := gin.H{
		"username": rsp.Username,
		"nickname": rsp.Nickname,
		"pic_path": rsp.PicPath,
	}
	logger.Debug("editUserProfile", "requestID:", requestID, "data:", data)

	response.SuccessDataRsp(ctx, data)
}

// UploadPicHandler 上传图片接口
func (ctrl *UserController) UploadPicHandler(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	token := ctx.GetHeader("Authorization")

	// rcp调用：认证
	rpcCtx := util.GenRpcCtxWithRequestID(ctx)
	_, err := client.RpcClient.Auth(rpcCtx, &proto.AuthRequest{
		Token: token,
	})
	if err != nil {
		logger.Error("uploadPic", "requestID:", requestID,
			"token:", token, "call auth error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 验证图片格式
	req := request.UploadPicRequest{}
	if err := request.Validate(ctx, &req, request.UploadPicValid); err != nil {
		logger.Warn("uploadPic", "requestID:", requestID, "request valid error", err.Error())
		return
	}

	// 保存图片
	avatar, err := util.SaveUploadPic(ctx, req.Avatar)
	if err != nil {
		logger.Error("uploadPic", "requestID:", requestID,
			"save upload pic error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// rpc调用：更新用户信息
	_, err = client.RpcClient.EditUserProfile(rpcCtx, &proto.EditUserRequest{
		Token:   token,
		PicPath: avatar,
	})
	if err != nil {
		logger.Error("uploadPic", "requestID:", requestID,
			"call editUserProfile error", err.Error())
		response.RpcRspToHttpRsp(err, ctx)
		return
	}

	// 返回结果
	data := gin.H{
		"avatar": util.WrapUploadPic(avatar),
	}
	logger.Debug("uploadPic", "requestID:", requestID, "upload success, data", data)

	response.SuccessDataRsp(ctx, data)
}
