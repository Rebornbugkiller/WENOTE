package repo

import (
	"wenote-backend/internal/model"
)

// AttachmentRepo 附件数据访问层
type AttachmentRepo struct{}

// NewAttachmentRepo 创建附件仓库实例
func NewAttachmentRepo() *AttachmentRepo {
	return &AttachmentRepo{}
}

// Create 创建附件记录
func (r *AttachmentRepo) Create(attachment *model.NoteAttachment) error {
	return DB.Create(attachment).Error
}

// GetByID 根据ID获取附件
func (r *AttachmentRepo) GetByID(id uint64) (*model.NoteAttachment, error) {
	var attachment model.NoteAttachment
	err := DB.First(&attachment, id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// GetByNoteID 获取笔记的所有附件
func (r *AttachmentRepo) GetByNoteID(noteID uint64) ([]*model.NoteAttachment, error) {
	var attachments []*model.NoteAttachment
	err := DB.Where("note_id = ?", noteID).Find(&attachments).Error
	return attachments, err
}

// Delete 删除附件记录
func (r *AttachmentRepo) Delete(id uint64) error {
	return DB.Delete(&model.NoteAttachment{}, id).Error
}

// DeleteByNoteID 删除笔记的所有附件
func (r *AttachmentRepo) DeleteByNoteID(noteID uint64) error {
	return DB.Where("note_id = ?", noteID).Delete(&model.NoteAttachment{}).Error
}

