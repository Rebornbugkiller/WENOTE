package service

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
)

// UserService 用户服务
type UserService struct {
	userRepo *repo.UserRepo
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repo.NewUserRepo(),
	}
}

// GetProfile 获取用户信息
func (s *UserService) GetProfile(userID uint64) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
