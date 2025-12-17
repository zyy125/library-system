package service

import (
	"context"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"
	"library-system/utils"
	"time"
	"errors"
	"gorm.io/gorm"
	"math"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) Register(ctx context.Context, req *request.UserRegisterRequest) (*response.UserRegisterResponse, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 判断是否存在相同用户名
	if _, err := s.userRepo.GetUserByUsername(ctx, user.Username); err == nil {
		return nil, common.ErrUsernameExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 判断是否存在相同邮箱
	if _, err := s.userRepo.GetUserByUsername(ctx, user.Username); err == nil {
		return nil, common.ErrEmailExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 加密密码
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPwd

	// 调用数据库函数
	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	// 构建返回值
	data := &response.UserRegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return data, nil
}

func (s *UserService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {
	// 调用数据库函数
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
    		return nil, common.ErrInvalidAuth
		}
		return nil, err
	}

	if user.Status == "disabled" {
		return nil, common.ErrUserDisabled
	}

	// 验证密码
	err = utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		return nil, common.ErrInvalidAuth
	}

	accessToken, refreshToken, tokenID, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, common.ErrInternalServer
	}

	err = repository.Rdb.StoreRefreshToken(ctx, user.ID, tokenID, refreshToken)
	if err != nil {
		return nil, err
	}

	userResponse := response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	// 构建返回值
	data := &response.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User:         userResponse,
	}

	return data, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *request.UserRefreshTokenRequest) (*response.UserTokenRefreshResponse, error) {
	// 验证 Refresh Token
	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, common.ErrInvalidToken
	}

	// 从 Redis 验证 Token 是否存在
	storedToken, err := repository.Rdb.GetRefreshToken(ctx, claims.UserID, claims.TokenID)
	if err != nil || storedToken != req.RefreshToken {
		return nil, common.ErrInvalidToken
	}

	// 检查是否在黑名单
	inBlacklist, err := repository.Rdb.IsInBlacklist(ctx, claims.TokenID)
	if err != nil {
		return nil, common.ErrInternalServer
	}
	if inBlacklist {
		return nil, common.ErrInvalidToken
	}

	// 查询用户信息（获取最新的 role 等信息）
	user, err := s.userRepo.GetUserByUserID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
    		return nil, common.ErrInvalidAuth
		}
		return nil, err
	}

	// 生成新的 Token Pair
	newAccessToken, newRefreshToken, newTokenID, err := utils.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Role,
	)
	if err != nil {
		return nil, common.ErrInternalServer
	}

	// 删除旧的 Refresh Token
	if err := repository.Rdb.DeleteRefreshToken(ctx, claims.UserID, claims.TokenID); err != nil {
		return nil, common.ErrInternalServer
	}

	// 存储新的 Refresh Token
	if err := repository.Rdb.StoreRefreshToken(ctx, user.ID, newTokenID, newRefreshToken); err != nil {
		return nil, common.ErrInternalServer
	}

	userResponse := response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	data := &response.UserTokenRefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		User:         userResponse,
	}

	return data, nil
}

func (s *UserService) Logout(ctx context.Context, userID uint64, tokenID string) error {
	if err := repository.Rdb.DeleteRefreshToken(ctx, userID, tokenID); err != nil {
		return err
	}

	if err := repository.Rdb.AddToBlacklist(ctx, tokenID, 24*time.Hour); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserMsg(ctx context.Context, userID uint64) (*response.GetUserMsgResponse, error) {
	user, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	data := &response.GetUserMsgResponse{
		ID:             userID,
		Username:       user.Username,
		Email:          user.Email,
		Phone:          user.Phone,
		Role:           user.Phone,
		Status:         user.Status,
		BorrowLimit:    user.BorrowLimit,
		BorrowingCount: user.BorrowingCount,
		OverdueCount:   user.OverdueCount,
		CreatedAt:		user.CreatedAt.UTC().Format(time.RFC3339),
	}

	return data, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID uint64, req *request.UpdateUserRequest) (*response.UpdateUserResponse, error) {
	if req.Email == nil &&
		req.Phone == nil&&
		req.Username == nil {
		return nil, common.ErrBadRequest
	}

	updates := make(map[string]interface{})

	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}

	if err := s.userRepo.UpdateUserFields(ctx, userID, updates); err != nil {
		return nil, err
	}

	data := &response.UpdateUserResponse{
		ID:        userID,
		Username:  *req.Username,
		Email:     *req.Email,
		Phone:     *req.Phone,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	return data, nil
}

func (s *UserService) ChangePwd(ctx context.Context, userID uint64, tokenID string, req *request.ChangePwdRequest) error {
	user, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := utils.CheckPassword(user.Password, req.OldPassword); err != nil {
		return common.ErrInvalidAuth
	}

	new, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	update := make(map[string]interface{})
	update["password"] = new

	err = s.userRepo.UpdateUserFields(ctx, userID, update)
	if err != nil {
		return err
	}

	if err := repository.Rdb.DeleteRefreshToken(ctx, userID, tokenID); err != nil {
		return err
	}

	if err := repository.Rdb.AddToBlacklist(ctx, tokenID, 24*time.Hour); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserList(ctx context.Context, req *request.GetUserListRequest) (*response.GetUserListResponse, error) {
	// 参数校验
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Role != "" && req.Role != "admin" && req.Role != "user" {
		return nil, common.ErrNotFound
	}
	if req.Status != "" && req.Status != "active" && req.Status != "disabled" {
		return nil, common.ErrNotFound
	}

	users, count, err := s.userRepo.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}

	list := make([]response.GetUserMsgResponse, 0, len(users))
	for _, u := range users {
		list = append(list, response.GetUserMsgResponse{
			ID:             u.ID,
			Username:       u.Username,
			Email:          u.Email,
			Phone:          u.Phone,
			Role:           u.Role,
			Status:         u.Status,
			BorrowLimit:    u.BorrowLimit,
			BorrowingCount: u.BorrowingCount,
			OverdueCount:   u.OverdueCount,
			CreatedAt:      u.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	totalPages := int(math.Ceil(float64(count) / float64(req.Limit)))
	res := &response.GetUserListResponse{
		Total: int(count),
		Page: req.Page,	
		Limit: req.Limit,
		TotalPages: totalPages,
		Users: list,
	}

	return res, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.CreateUserResponse, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Email: req.Email,
		Phone: req.Phone,
		Role: req.Role,
		BorrowLimit: req.BorrowLimit,
	}

	// 判断是否存在相同用户名
	if _, err := s.userRepo.GetUserByUsername(ctx, user.Username); err == nil {
		return nil, common.ErrUsernameExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 判断是否存在相同邮箱
	if _, err := s.userRepo.GetUserByUsername(ctx, user.Username); err == nil {
		return nil, common.ErrEmailExist
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 加密密码
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPwd	

	// 调用数据库函数
	err = s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	// 构建返回值
	data := &response.CreateUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Status:      user.Status,
		BorrowLimit: user.BorrowLimit,
		CreatedAt:   user.CreatedAt.UTC().Format(time.RFC3339),
	}

	return data, nil
}

func (s *UserService) UpdateUserByAdmin(ctx context.Context, id uint64, req *request.UpdateUserByAdminRequest) (*response.UpdateUserByAdminResponse, error) {
	updates := make(map[string]interface{})

	if req.Email == nil &&
		req.Role == nil &&
		req.Status == nil &&
		req.BorrowLimit == nil &&
		req.Phone == nil&&
		req.Username == nil {
		return nil, common.ErrBadRequest
	}

	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.BorrowLimit != nil {
		updates["borrow_limit"] = *req.BorrowLimit
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := s.userRepo.UpdateUserFields(ctx, id, updates); err != nil {
		return nil, err
	}

	res := &response.UpdateUserByAdminResponse{
		ID:          id,
		Username:    *req.Username,
		Email:       *req.Email,
		Phone:       *req.Phone,
		Role:        *req.Role,
		Status:      *req.Status,
		BorrowLimit: *req.BorrowLimit,
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint64) error {
	user, err := s.userRepo.GetUserByUserID(ctx, id)
	if err != nil {
		return err
	}

	if user.BorrowingCount > 0 {
		bizErr := common.NewBizError(400, "无法删除该用户", 400)
		details := make(map[string]interface{})
		details["reason"] = "用户有未归还图书"
		details["unreturned_books"] = user.BorrowingCount
		bizErr.WithDetails(details)
		return bizErr
	}

	err = s.userRepo.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}