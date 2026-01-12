package repo

import (
	"wenote-backend/internal/model"
	"errors"

	"gorm.io/gorm"
)

// NotebookRepo 笔记本数据访问
type NotebookRepo struct{}

// NewNotebookRepo 创建 NotebookRepo 实例
func NewNotebookRepo() *NotebookRepo {
	return &NotebookRepo{}
}

// Create 创建笔记本
func (r *NotebookRepo) Create(notebook *model.Notebook) error {
	return DB.Create(notebook).Error
}

// GetByIDAndUserID 根据ID和用户ID获取笔记本
func (r *NotebookRepo) GetByIDAndUserID(id, userID uint64) (*model.Notebook, error) {
	var notebook model.Notebook
	err := DB.Where("id = ? AND user_id = ?", id, userID).First(&notebook).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &notebook, err
}

// Update 更新笔记本
func (r *NotebookRepo) Update(notebook *model.Notebook) error {
	return DB.Save(notebook).Error
}

// Delete 删除笔记本
func (r *NotebookRepo) Delete(id uint64) error {
	return DB.Delete(&model.Notebook{}, id).Error
}

// ListByUserID 获取用户的笔记本列表
func (r *NotebookRepo) ListByUserID(userID uint64) ([]*model.Notebook, error) {
	var notebooks []*model.Notebook
	err := DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notebooks).Error
	return notebooks, err
}

// CountNotesByNotebookID 统计笔记本下的笔记数量（不含已删除）
func (r *NotebookRepo) CountNotesByNotebookID(notebookID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.Note{}).
		Where("notebook_id = ? AND deleted_at IS NULL", notebookID).
		Count(&count).Error
	return count, err
}

// ExistsByUserIDAndName 检查用户是否已有同名笔记本
func (r *NotebookRepo) ExistsByUserIDAndName(userID uint64, name string, excludeID uint64) (bool, error) {
	var count int64
	query := DB.Model(&model.Notebook{}).Where("user_id = ? AND name = ?", userID, name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetOrCreateDefault 获取或创建默认笔记本
func (r *NotebookRepo) GetOrCreateDefault(userID uint64) (*model.Notebook, error) {
	var notebook model.Notebook
	err := DB.Where("user_id = ? AND is_default = ?", userID, true).First(&notebook).Error
	if err == nil {
		return &notebook, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// 创建默认笔记本
	notebook = model.Notebook{
		UserID:    userID,
		Name:      "未分类",
		IsDefault: true,
	}
	if err := DB.Create(&notebook).Error; err != nil {
		return nil, err
	}
	return &notebook, nil
}
