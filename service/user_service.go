package service

import (
	"context"
	"library-system/model"
	"library-system/repository"
	"library-system/common"
	"library-system/utils"

)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService{
	return &UserService{userRepo: repo}
}

func (s *UserService) Register(ctx context.Context, user *model.User) error {
	existingUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return common.ErrUsernameExist
	}

	existingUser, err = s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return common.ErrEmailExist
	}

	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPwd

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}