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

// Update 更新用户信息
func (r *UserRepo) Update(user *model.User) error {
	return DB.Save(user).Error
}

// UpdatePassword 更新用户密码
func (r *UserRepo) UpdatePassword(userID uint64, passwordHash string) error {
	return DB.Model(&model.User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error
}

// ExistsByEmail 检查邮箱是否已被其他用户使用
func (r *UserRepo) ExistsByEmail(email string, excludeUserID uint64) (bool, error) {
	var count int64
	err := DB.Model(&model.User{}).Where("email = ? AND id != ?", email, excludeUserID).Count(&count).Error
	return count > 0, err
}

// Delete 删除用户及其所有关联数据
func (r *UserRepo) Delete(userID uint64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户的附件
		if err := tx.Where("user_id = ?", userID).Delete(&model.NoteAttachment{}).Error; err != nil {
			return err
		}
		// 删除用户的笔记
		if err := tx.Where("user_id = ?", userID).Delete(&model.Note{}).Error; err != nil {
			return err
		}
		// 删除用户的笔记本
		if err := tx.Where("user_id = ?", userID).Delete(&model.Notebook{}).Error; err != nil {
			return err
		}
		// 删除用户的标签
		if err := tx.Where("user_id = ?", userID).Delete(&model.Tag{}).Error; err != nil {
			return err
		}
		// 删除用户的游戏化数据
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserGamification{}).Error; err != nil {
			return err
		}
		// 删除用户的成就
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserAchievement{}).Error; err != nil {
			return err
		}
		// 最后删除用户
		if err := tx.Delete(&model.User{}, userID).Error; err != nil {
			return err
		}
		return nil
	})
}
