package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/util"
	"user_rpc/proto"
	userService "user_rpc/service/user"
)

// UserHandler 用户服务
type UserHandler struct {
}

// 用户模型转换为登录接口的返回格式
func userModelToLoginRsp(user userModel.User, token string) proto.LoginResponse {
	return proto.LoginResponse{
		Id:         user.ID,
		Username:   user.Username,
		Nickname:   user.Nickname,
		PicPath:    user.PicPath,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
		Token:      token,
	}
}

// 用户模型转换为接口的用户信息返回格式
func userModelToUserRsp(user userModel.User) proto.UserResponse {
	return proto.UserResponse{
		Id:         user.ID,
		Username:   user.Username,
		Nickname:   user.Nickname,
		PicPath:    user.PicPath,
		CreateTime: user.CreateTime.Unix(),
		UpdateTime: user.UpdateTime.Unix(),
	}
}

// Login 登录接口
func (u *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用login逻辑
	user, token, err := userService.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	rsp := userModelToLoginRsp(user, token)
	logger.Debug("login", "requestID:", requestID, "login success, rsp:", &rsp)

	return &rsp, nil
}

// Logout 登出接口
func (u *UserHandler) Logout(ctx context.Context, req *proto.AuthRequest) (*emptypb.Empty, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用登出逻辑
	token := req.Token
	if err := userService.Logout(ctx, req); err != nil {
		return nil, err
	}

	logger.Debug("logout", "requestID:", requestID, "logout success, token:", token)
	return &emptypb.Empty{}, nil
}

// GetUserProfile 获取用户信息接口
func (u *UserHandler) GetUserProfile(ctx context.Context, req *proto.AuthRequest) (*proto.UserResponse, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用获取用户信息逻辑
	user, err := userService.GetUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	rsp := userModelToUserRsp(user)
	logger.Debug("getUserProfile", "requestID:", requestID, "getUserProfile data:", &rsp)

	return &rsp, nil
}

// EditUserProfile 编辑用户信息接口
func (u *UserHandler) EditUserProfile(ctx context.Context, req *proto.EditUserRequest) (*proto.UserResponse, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用编辑用户逻辑
	username, err := userService.EditUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	rsp := userModelToUserRsp(userModel.User{
		Username: username,
		Nickname: req.Nickname,
		PicPath:  req.PicPath,
	})
	logger.Debug("editUserProfile", "requestID:", requestID, "data:", &rsp)

	return &rsp, nil
}

// CreateUserProfile 创建用户信息接口
func (u *UserHandler) CreateUserProfile(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用创建用户逻辑
	user, err := userService.CreateUserProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	rsp := userModelToUserRsp(user)
	logger.Debug("createUserProfile", "requestID:", requestID, "data:", &rsp)

	return &rsp, nil
}

// Auth 认证接口
func (u *UserHandler) Auth(ctx context.Context, req *proto.AuthRequest) (*emptypb.Empty, error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 调用认证逻辑
	token := req.Token
	if err := userService.Auth(ctx, req); err != nil {
		return nil, err
	}

	logger.Debug("auth", "requestID:", requestID, "auth success, token:", token)
	return &emptypb.Empty{}, nil
}
