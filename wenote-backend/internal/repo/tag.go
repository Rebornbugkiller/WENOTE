package repo

import (
	"wenote-backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

// TagRepo 标签数据访问
type TagRepo struct{}

// NewTagRepo 创建 TagRepo 实例
func NewTagRepo() *TagRepo {
	return &TagRepo{}
}

// Create 创建标签
func (r *TagRepo) Create(tag *model.Tag) error {
	return DB.Create(tag).Error
}

// GetByIDAndUserID 根据ID和用户ID获取标签
func (r *TagRepo) GetByIDAndUserID(id, userID uint64) (*model.Tag, error) {
	var tag model.Tag
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// Delete 删除标签
func (r *TagRepo) Delete(id uint64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 先删除关联关系
		if err := tx.Where("tag_id = ?", id).Delete(&model.NoteTag{}).Error; err != nil {
			return err
		}
		// 再删除标签
		return tx.Delete(&model.Tag{}, id).Error
	})
}

// Update 更新标签
func (r *TagRepo) Update(tag *model.Tag) error {
	return DB.Save(tag).Error
}

// ListByUserID 获取用户的标签列表
func (r *TagRepo) ListByUserID(userID uint64) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tags).Error
	return tags, err
}

// ListByIDs 根据ID列表获取标签
func (r *TagRepo) ListByIDs(ids []uint64) ([]*model.Tag, error) {
	var tags []*model.Tag
	err := DB.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

// ExistsByNameAndUserID 检查用户是否已有同名标签
func (r *TagRepo) ExistsByNameAndUserID(name string, userID uint64) (bool, error) {
	var count int64
	err := DB.Model(&model.Tag{}).
		Where("name = ? AND user_id = ?", name, userID).
		Count(&count).Error
	return count > 0, err
}

// CountNotesByTagID 统计标签下的笔记数量
func (r *TagRepo) CountNotesByTagID(tagID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.NoteTag{}).Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// GetOrCreate 获取或创建标签
func (r *TagRepo) GetOrCreate(userID uint64, name string) (*model.Tag, error) {
	var tag model.Tag
	err := DB.Where("user_id = ? AND name = ?", userID, name).First(&tag).Error
	if err == nil {
		return &tag, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		tag = model.Tag{UserID: userID, Name: name}
		err = DB.Create(&tag).Error
		return &tag, err
	}

	return nil, err
}
