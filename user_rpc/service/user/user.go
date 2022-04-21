package user

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	userDao "user_rpc/dao/user"
	userModel "user_rpc/model/user"
	"user_rpc/pkg/logger"
	"user_rpc/pkg/util"
	"user_rpc/proto"
)

// 检查用户是否登录
func checkAuth(moduleName string, token string, requestID string) (string, error) {
	username, err := userDao.GetTokenCache(token)
	if err == nil && len(username) == 0 {
		logger.Warn(moduleName, "requestID:", requestID, "not token cache, token:", token)
		return "", status.Error(codes.Unauthenticated, "unauthorized")
	}
	if err != nil {
		logger.Error(moduleName, "requestID:", requestID, "get token cache, err:", err)
		return "", status.Error(codes.Internal, err.Error())
	}
	return username, nil
}

// 根据用户名查询用户
func getUserByUsername(ctx context.Context, moduleName string, username string) (userModel.User, error) {
	requestID := util.GetRequestIDFromContext(ctx)
	user, err := userDao.GetByUsername(username)
	if err == nil && user.ID == 0 {
		logger.Warn(moduleName, "requestID:", requestID,
			"get user db not exist, username:", username)
		return user, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		logger.Error(moduleName, "requestID:", requestID, "get user db, err:", err)
		return user, status.Error(codes.Internal, err.Error())
	}
	return user, nil
}

// Login 登录逻辑
func Login(ctx context.Context, req *proto.LoginRequest) (user userModel.User, token string, err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 根据用户名查询用户
	user, err = userDao.GetByUsername(req.Username)
	if err == nil && user.ID == 0 {
		logger.Warn("login", "requestID:", requestID, "username not exist, username:", req.Username)
		return user, "", status.Error(codes.NotFound, "username not exist")
	}
	if err != nil {
		logger.Error("login", "requestID:", requestID, fmt.Sprintf("get user error: %s, username: %s", err.Error(), req.Username))
		return user, "", status.Error(codes.Internal, err.Error())
	}

	// 校验密码
	if ok := util.ComparePwdHash(req.Password, user.Password, user.Salt); !ok {
		logger.Warn("login", "requestID:", requestID, fmt.Sprintf("password incorrect, password: %s, salt:%s",
			req.Password, user.Salt))
		return user, "", status.Error(codes.PermissionDenied, "password incorrect")
	}

	// 生成并写入token
	token = util.GenerateToken(user.Username)
	if err := userDao.SetTokenCache(token, user.Username); err != nil {
		logger.Error("login", "requestID:", requestID, fmt.Sprintf(
			"set token cache error: %s, username: %s", err.Error(), req.Username))
		return user, "", status.Error(codes.Internal, err.Error())
	}

	// 写入user cache
	if err = userDao.SetUserCache(user.Username, user); err != nil {
		logger.Error("login", "requestID:", requestID, fmt.Sprintf(
			"set user cache error: %s, username: %s", err.Error(), req.Username))
		return user, "", status.Error(codes.Internal, err.Error())
	}

	return user, token, nil
}

// Logout 登出逻辑
func Logout(ctx context.Context, req *proto.AuthRequest) (err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 检查用户是否登录
	token := req.Token
	if _, err = checkAuth("logout", token, requestID); err != nil {
		return err
	}

	// 删除用户token
	if err := userDao.DelTokenCache(token); err != nil {
		logger.Error("logout", "requestID:", requestID, "delete token cache, err:", err)
		return status.Error(codes.Internal, err.Error())
	}

	return
}

// GetUserProfile 获取用户信息逻辑
func GetUserProfile(ctx context.Context, req *proto.AuthRequest) (user userModel.User, err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 检查用户是否登录
	token := req.Token
	username, err := checkAuth("getUserProfile", token, requestID)
	if err != nil {
		return user, err
	}

	// 查询user cache
	user, err = userDao.GetUserCache(username)
	if err != nil {
		logger.Error("getUserProfile", "requestID:", requestID, "get user cache, err:", err)
		return user, status.Error(codes.Internal, err.Error())
	}

	// user cache不存在
	if err == nil && user.ID == 0 {
		logger.Debug("getUserProfile", "requestID:", requestID, "not user cache, username:", username)

		// 查询user db
		user, err = getUserByUsername(ctx, "getUserProfile", username)
		if err != nil {
			return user, err
		}

		// 写入user cache
		if err = userDao.SetUserCache(username, user); err != nil {
			logger.Error("getUserProfile", "requestID:", requestID, "set user cache, err:", err)
			return user, status.Error(codes.Internal, err.Error())
		}
	}

	return user, nil
}

// EditUserProfile 编辑用户信息逻辑
func EditUserProfile(ctx context.Context, req *proto.EditUserRequest) (username string, err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 检查用户是否登录
	token := req.Token
	username, err = checkAuth("editUserProfile", token, requestID)
	if err != nil {
		return "", err
	}

	// 更新db user
	if len(req.Nickname) > 0 {
		err = userDao.UpdateNickNameByUsername(username, req.Nickname)
	} else {
		err = userDao.UpdatePicPathByUsername(username, req.PicPath)
	}
	if err != nil {
		logger.Error("editUserProfile", "requestID:", requestID, "update db user, err:", err)
		return "", status.Error(codes.Internal, err.Error())
	}

	// 删除db cache
	if err = userDao.DelUserCache(username); err != nil {
		logger.Error("editUserProfile", "requestID:", requestID, "delete user cache, err:", err)
		return "", status.Error(codes.Internal, err.Error())
	}

	return username, nil
}

// CreateUserProfile 创建用户信息
func CreateUserProfile(ctx context.Context, req *proto.CreateUserRequest) (user userModel.User, err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	// 检查用户是否已存在
	user, err = userDao.GetByUsername(req.Username)
	if err != nil {
		logger.Error("createUserProfile", "requestID:", requestID, fmt.Sprintf(
			"get user error: %s, username: %s", err.Error(), req.Username))
		return user, status.Error(codes.Internal, err.Error())
	}
	if user.ID != 0 {
		logger.Warn("createUserProfile", "requestID:", requestID, fmt.Sprintf(
			"username already exist, username: %s", req.Username))
		return user, status.Error(codes.AlreadyExists, "username already exist")
	}

	// 写入user db
	user, err = userDao.CreateOne(req.Username, req.Password, req.Nickname)
	if err != nil {
		logger.Error("createUserProfile", "requestID:", requestID, fmt.Sprintf(
			"create user error: %s, username: %s", err.Error(), req.Username))
		return user, status.Error(codes.Internal, err.Error())
	}

	return user, nil
}

// Auth 认证逻辑
func Auth(ctx context.Context, req *proto.AuthRequest) (err error) {
	requestID := util.GetRequestIDFromContext(ctx)

	token := req.Token
	if _, err = checkAuth("auth", token, requestID); err != nil {
		return err
	}

	return nil
}
