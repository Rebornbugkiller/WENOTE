package repo

import (
	"wenote-backend/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

// NoteRepo 笔记数据访问
type NoteRepo struct{}

// NewNoteRepo 创建 NoteRepo 实例
func NewNoteRepo() *NoteRepo {
	return &NoteRepo{}
}

// Create 创建笔记
func (r *NoteRepo) Create(note *model.Note) error {
	return DB.Create(note).Error
}

// GetByID 根据ID获取笔记（包含标签）
func (r *NoteRepo) GetByID(id uint64) (*model.Note, error) {
	var note model.Note
	err := DB.Preload("Tags").Where("id = ? AND deleted_at IS NULL", id).First(&note).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

// GetByIDAndUserID 根据ID和用户ID获取笔记
func (r *NoteRepo) GetByIDAndUserID(id, userID uint64) (*model.Note, error) {
	var note model.Note
	err := DB.Preload("Tags").
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		First(&note).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

// GetDeletedByIDAndUserID 获取已删除的笔记
func (r *NoteRepo) GetDeletedByIDAndUserID(id, userID uint64) (*model.Note, error) {
	var note model.Note
	err := DB.Where("id = ? AND user_id = ? AND deleted_at IS NOT NULL", id, userID).
		First(&note).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &note, err
}

// Update 更新笔记
func (r *NoteRepo) Update(note *model.Note) error {
	return DB.Save(note).Error
}

// UpdateFields 更新指定字段
func (r *NoteRepo) UpdateFields(id uint64, fields map[string]interface{}) error {
	return DB.Model(&model.Note{}).Where("id = ?", id).Updates(fields).Error
}

// UpdateFieldsWithoutTime 更新指定字段但不更新 updated_at
func (r *NoteRepo) UpdateFieldsWithoutTime(id uint64, fields map[string]interface{}) error {
	return DB.Model(&model.Note{}).Where("id = ?", id).UpdateColumns(fields).Error
}

// SoftDelete 软删除笔记
func (r *NoteRepo) SoftDelete(id uint64) error {
	now := time.Now()
	return DB.Model(&model.Note{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// Restore 恢复已删除的笔记
func (r *NoteRepo) Restore(id uint64) error {
	return DB.Model(&model.Note{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// List 获取笔记列表
func (r *NoteRepo) List(userID uint64, req *model.NoteListReq) ([]*model.Note, int64, error) {
	var notes []*model.Note
	var total int64

	query := DB.Model(&model.Note{}).Where("user_id = ? AND deleted_at IS NULL", userID)

	// 笔记本筛选
	if req.NotebookID != nil {
		query = query.Where("notebook_id = ?", *req.NotebookID)
	}

	// 星标筛选
	if req.IsStarred != nil {
		query = query.Where("is_starred = ?", *req.IsStarred)
	}

	// 置顶筛选
	if req.IsPinned != nil {
		query = query.Where("is_pinned = ?", *req.IsPinned)
	}

	// 关键词搜索 - 使用MySQL全文搜索
	if req.Keyword != "" {
		// 使用全文搜索（IN NATURAL LANGUAGE MODE）
		// 利用notes表的FULLTEXT INDEX ft_title_content (title, content)
		query = query.Where(
			"MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)",
			req.Keyword,
		)
	}

	// 标签筛选
	if req.TagID != nil {
		query = query.Joins("JOIN note_tags ON notes.id = note_tags.note_id").
			Where("note_tags.tag_id = ?", *req.TagID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	
	// 如果有关键词搜索，按相关性排序；否则按置顶和更新时间排序
	var err error
	if req.Keyword != "" {
		// 全文搜索时按相关性排序，置顶次之
		// 使用Expr构建MATCH表达式用于排序
		err = query.Preload("Tags").
			Order("is_pinned DESC").
			Order(DB.Raw("MATCH(title, content) AGAINST(?) DESC", req.Keyword)).
			Offset(offset).
			Limit(req.PageSize).
			Find(&notes).Error
	} else {
		// 普通查询按置顶和更新时间排序
		err = query.Preload("Tags").
			Order("is_pinned DESC, updated_at DESC").
			Offset(offset).
			Limit(req.PageSize).
			Find(&notes).Error
	}

	return notes, total, err
}

// UpdateAIStatus 更新 AI 任务状态
func (r *NoteRepo) UpdateAIStatus(id uint64, status model.AIStatus, aiError string) error {
	fields := map[string]interface{}{
		"ai_status": status,
		"ai_error":  aiError,
	}
	return DB.Model(&model.Note{}).Where("id = ?", id).Updates(fields).Error
}

// UpdateAIResult 更新 AI 处理结果
func (r *NoteRepo) UpdateAIResult(id uint64, summary string, suggestedTags []string) error {
	fields := map[string]interface{}{
		"ai_status":      model.AIStatusDone,
		"summary":        summary,
		"suggested_tags": model.StringSlice(suggestedTags),
		"ai_error":       "",
	}
	return DB.Model(&model.Note{}).Where("id = ?", id).Updates(fields).Error
}

// CountByNotebookID 统计笔记本下的笔记数量
func (r *NoteRepo) CountByNotebookID(notebookID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.Note{}).
		Where("notebook_id = ? AND deleted_at IS NULL", notebookID).
		Count(&count).Error
	return count, err
}

// CountByUserID 统计用户的笔记数量
func (r *NoteRepo) CountByUserID(userID uint64) (int64, error) {
	var count int64
	err := DB.Model(&model.Note{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// ReplaceNoteTags 替换笔记的所有标签
func (r *NoteRepo) ReplaceNoteTags(noteID uint64, tagIDs []uint64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除现有关联
		if err := tx.Where("note_id = ?", noteID).Delete(&model.NoteTag{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新关联
		if len(tagIDs) > 0 {
			noteTags := make([]model.NoteTag, len(tagIDs))
			for i, tagID := range tagIDs {
				noteTags[i] = model.NoteTag{NoteID: noteID, TagID: tagID}
			}
			if err := tx.Create(&noteTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetNoteTagIDs 获取笔记的所有标签ID
func (r *NoteRepo) GetNoteTagIDs(noteID uint64) ([]uint64, error) {
	var tagIDs []uint64
	err := DB.Model(&model.NoteTag{}).
		Where("note_id = ?", noteID).
		Pluck("tag_id", &tagIDs).Error
	return tagIDs, err
}

// ListDeleted 获取已删除的笔记列表（回收站）
func (r *NoteRepo) ListDeleted(userID uint64, page, pageSize int) ([]*model.Note, int64, error) {
	var notes []*model.Note
	var total int64

	query := DB.Unscoped().Model(&model.Note{}).Where("user_id = ? AND deleted_at IS NOT NULL", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按删除时间倒序
	offset := (page - 1) * pageSize
	err := query.Preload("Tags").
		Order("deleted_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notes).Error

	return notes, total, err
}

// HardDeleteOld 硬删除超过指定天数的软删除笔记
func (r *NoteRepo) HardDeleteOld(days int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days)

	result := DB.Unscoped().
		Where("deleted_at IS NOT NULL AND deleted_at < ?", cutoffTime).
		Delete(&model.Note{})

	return result.RowsAffected, result.Error
}

// ClearSuggestedTags 清空建议标签
func (r *NoteRepo) ClearSuggestedTags(noteID uint64) error {
	return DB.Model(&model.Note{}).Where("id = ?", noteID).
		Update("suggested_tags", nil).Error
}

// FilterByUserID 过滤属于用户的笔记ID（包括已删除的）
func (r *NoteRepo) FilterByUserID(noteIDs []uint64, userID uint64) ([]uint64, error) {
	var validIDs []uint64
	err := DB.Unscoped().Model(&model.Note{}).
		Where("id IN ? AND user_id = ?", noteIDs, userID).
		Pluck("id", &validIDs).Error
	return validIDs, err
}

// BatchSoftDelete 批量软删除
func (r *NoteRepo) BatchSoftDelete(noteIDs []uint64) (int64, error) {
	now := time.Now()
	result := DB.Model(&model.Note{}).
		Where("id IN ? AND deleted_at IS NULL", noteIDs).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

// BatchRestore 批量恢复
func (r *NoteRepo) BatchRestore(noteIDs []uint64) (int64, error) {
	result := DB.Model(&model.Note{}).
		Where("id IN ? AND deleted_at IS NOT NULL", noteIDs).
		Update("deleted_at", nil)
	return result.RowsAffected, result.Error
}

// BatchHardDelete 批量硬删除（永久删除）
func (r *NoteRepo) BatchHardDelete(noteIDs []uint64) (int64, error) {
	result := DB.Unscoped().Where("id IN ?", noteIDs).Delete(&model.Note{})
	return result.RowsAffected, result.Error
}

// EmptyTrash 清空用户回收站（永久删除所有已删除笔记）
func (r *NoteRepo) EmptyTrash(userID uint64) (int64, error) {
	result := DB.Unscoped().Where("user_id = ? AND deleted_at IS NOT NULL", userID).Delete(&model.Note{})
	return result.RowsAffected, result.Error
}

// BatchUpdateNotebook 批量移动到笔记本
func (r *NoteRepo) BatchUpdateNotebook(noteIDs []uint64, notebookID uint64) (int64, error) {
	result := DB.Model(&model.Note{}).
		Where("id IN ? AND deleted_at IS NULL", noteIDs).
		Update("notebook_id", notebookID)
	return result.RowsAffected, result.Error
}

// SoftDeleteByNotebookID 按笔记本ID软删除所有笔记
func (r *NoteRepo) SoftDeleteByNotebookID(notebookID uint64) (int64, error) {
	now := time.Now()
	result := DB.Model(&model.Note{}).
		Where("notebook_id = ? AND deleted_at IS NULL", notebookID).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

