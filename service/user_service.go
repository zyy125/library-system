package service

import (
	"context"
	"library-system/model"
	"library-system/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService{
	return &UserService{userRepo: repo}
}

func (s *UserService) Register(ctx context.Context, user *model.User) error {
	return s.userRepo.CreateUser(ctx, user)
}