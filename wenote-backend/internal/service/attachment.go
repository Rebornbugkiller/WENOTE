package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
)

// AttachmentService 附件服务
type AttachmentService struct {
	repo     *repo.AttachmentRepo
	noteRepo *repo.NoteRepo
}

// NewAttachmentService 创建附件服务实例
func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		repo:     repo.NewAttachmentRepo(),
		noteRepo: repo.NewNoteRepo(),
	}
}

const (
	// 上传目录配置
	UploadBasePath = "./uploads/images"
	MaxFileSize    = 5 * 1024 * 1024 // 5MB
)

// AllowedImageTypes 允许的图片类型
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// UploadImage 上传图片
func (s *AttachmentService) UploadImage(userID uint64, noteID uint64, file *multipart.FileHeader) (*model.AttachmentUploadResp, error) {
	// 1. 验证笔记所有权
	note, err := s.noteRepo.GetByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("笔记不存在")
	}
	if note.UserID != userID {
		return nil, fmt.Errorf("无权限访问该笔记")
	}

	// 2. 验证文件大小
	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("文件大小超过限制（最大5MB）")
	}

	// 3. 验证文件类型
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return nil, fmt.Errorf("不支持的文件类型")
	}

	// 4. 生成存储路径
	timestamp := time.Now().UnixNano()
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", timestamp, noteID, ext)
	userDir := filepath.Join(UploadBasePath, fmt.Sprintf("user_%d", userID))
	storagePath := filepath.Join(userDir, filename)

	// 5. 确保目录存在
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 6. 保存文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 7. 生成访问URL（相对路径）
	url := fmt.Sprintf("/uploads/images/user_%d/%s", userID, filename)

	// 8. 保存数据库记录
	attachment := &model.NoteAttachment{
		NoteID:      noteID,
		UserID:      userID,
		Filename:    file.Filename,
		FileSize:    int(file.Size),
		MimeType:    file.Header.Get("Content-Type"),
		StoragePath: storagePath,
		URL:         url,
	}

	if err := s.repo.Create(attachment); err != nil {
		// 删除已保存的文件
		os.Remove(storagePath)
		return nil, fmt.Errorf("保存附件记录失败: %v", err)
	}

	return &model.AttachmentUploadResp{
		ID:  attachment.ID,
		URL: url,
	}, nil
}

// GetAttachments 获取笔记的附件列表
func (s *AttachmentService) GetAttachments(userID uint64, noteID uint64) ([]*model.NoteAttachment, error) {
	// 验证笔记所有权
	note, err := s.noteRepo.GetByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("笔记不存在")
	}
	if note.UserID != userID {
		return nil, fmt.Errorf("无权限访问该笔记")
	}

	return s.repo.GetByNoteID(noteID)
}

// DeleteAttachment 删除附件
func (s *AttachmentService) DeleteAttachment(userID uint64, attachmentID uint64) error {
	// 1. 获取附件信息
	attachment, err := s.repo.GetByID(attachmentID)
	if err != nil {
		return fmt.Errorf("附件不存在")
	}

	// 2. 验证所有权
	if attachment.UserID != userID {
		return fmt.Errorf("无权限删除该附件")
	}

	// 3. 删除文件
	if err := os.Remove(attachment.StoragePath); err != nil {
		// 文件可能已被删除，只记录日志，不影响数据库删除
		fmt.Printf("删除文件失败: %v\n", err)
	}

	// 4. 删除数据库记录
	return s.repo.Delete(attachmentID)
}

