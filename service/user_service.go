package service

import (
	"context"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"library-system/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService{
	return &UserService{userRepo: repo}
}

func (s *UserService) Register(ctx context.Context, req request.UserRegisterRequest) (response.UserRegisterResponse, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email: req.Email,
		Phone: req.Phone,
	}

	// 判断是否存在相同用户名
	existingUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return response.UserRegisterResponse{}, err
	}
	if existingUser != nil {
		return response.UserRegisterResponse{}, common.ErrUsernameExist
	}

	// 判断是否存在相同邮箱
	existingUser, err = s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return response.UserRegisterResponse{}, err
	}
	if existingUser != nil {
	return response.UserRegisterResponse{}, common.ErrEmailExist
	}

	// 加密密码
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return response.UserRegisterResponse{}, err
	}
	user.Password = hashedPwd

	// 调用数据库函数
	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return response.UserRegisterResponse{}, err
	}

	// 构建返回值
	data := response.UserRegisterResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
	}

	return data, nil
}

func (s *UserService) Login(ctx context.Context, req request.UserLoginRequest) (response.UserLoginResponse, error) {
	// 调用数据库函数
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return response.UserLoginResponse{}, err
	}
	if user == nil {
		return response.UserLoginResponse{}, common.ErrInvalidAuth
	}

	if user.Status == "disabled" {
		return response.UserLoginResponse{}, common.ErrUserDisabled
	}

	// 验证密码
	err = utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		return response.UserLoginResponse{}, common.ErrInvalidAuth
	}

	accessToken, refreshToken, tokenID, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return response.UserLoginResponse{}, common.ErrInternalServer
	}	

	err = repository.Rdb.StoreRefreshToken(ctx, user.ID, tokenID, refreshToken)
	if err != nil {
		return response.UserLoginResponse{}, err
	}

	userResponse := response.UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	}
	
	// 构建返回值
	data := response.UserLoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		TokenType: "Bearer",
		ExpiresIn: 86400,
		User: userResponse,
	}

	return data, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req request.UserRefreshTokenRequest) (response.UserTokenRefreshResponse, error) {
	// 验证 Refresh Token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return response.UserTokenRefreshResponse{}, common.ErrInvalidToken
	}

	// 从 Redis 验证 Token 是否存在
	storedToken, err := repository.Rdb.GetRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil || storedToken != req.RefreshToken {
		return response.UserTokenRefreshResponse{}, common.ErrInvalidToken
	}

	// 检查是否在黑名单
	inBlacklist, err := repository.Rdb.IsInBlacklist(ctx, claims.TokenID)
	if err != nil {
		return response.UserTokenRefreshResponse{}, common.ErrInternalServer
	}
	if inBlacklist {
		return response.UserTokenRefreshResponse{}, common.ErrInvalidToken
	}

	// 查询用户信息（获取最新的 role 等信息）
	user, err := s.userRepo.GetUserByUserID(ctx, claims.UserID)
	if err != nil {
		return response.UserTokenRefreshResponse{}, err
	}
	if user == nil {
		return response.UserTokenRefreshResponse{}, common.ErrInvalidToken
	}

	// 生成新的 Token Pair（推荐：同时刷新 Refresh Token）
	newAccessToken, newRefreshToken, newTokenID, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return response.UserTokenRefreshResponse{}, common.ErrInternalServer
	}

	// 删除旧的 Refresh Token
	if err := repository.Rdb.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID); err != nil {
		return response.UserTokenRefreshResponse{}, common.ErrInternalServer
	}

	// 存储新的 Refresh Token
	if err := repository.Rdb.StoreRefreshToken(ctx, user.ID, newTokenID, newRefreshToken); err != nil {
		return response.UserTokenRefreshResponse{}, common.ErrInternalServer
	}

	userResponse := response.UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	}

	data := response.UserTokenRefreshResponse{
		AccessToken: newAccessToken,	
		RefreshToken: newRefreshToken,
		TokenType: "Bearer",
		ExpiresIn: 86400,
		User: userResponse,
	}

	return data, nil
}