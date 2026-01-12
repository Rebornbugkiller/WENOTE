package repo

import (
	"errors"
	"wenote-backend/internal/model"

	"gorm.io/gorm"
)

// UserRepo 用户数据访问层
type UserRepo struct{}

// NewUserRepo 创建 UserRepo 实例
func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return DB.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepo) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := DB.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
