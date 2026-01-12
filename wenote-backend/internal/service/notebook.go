package service

import (
	"errors"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
)

var (
	ErrNotebookNotFound      = errors.New("笔记本不存在")
	ErrCannotDeleteDefault   = errors.New("默认笔记本不能删除")
	ErrNotebookNameDuplicate = errors.New("笔记本名称已存在")
)

// NotebookService 笔记本服务
type NotebookService struct {
	notebookRepo *repo.NotebookRepo
	noteRepo     *repo.NoteRepo
}

// NewNotebookService 创建笔记本服务实例
func NewNotebookService() *NotebookService {
	return &NotebookService{
		notebookRepo: repo.NewNotebookRepo(),
		noteRepo:     repo.NewNoteRepo(),
	}
}

// Create 创建笔记本
func (s *NotebookService) Create(userID uint64, req *model.NotebookCreateReq) (*model.Notebook, error) {
	// 检查同名笔记本
	exists, err := s.notebookRepo.ExistsByUserIDAndName(userID, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrNotebookNameDuplicate
	}

	notebook := &model.Notebook{
		UserID: userID,
		Name:   req.Name,
	}

	if err := s.notebookRepo.Create(notebook); err != nil {
		return nil, err
	}

	return notebook, nil
}

// GetByID 获取笔记本详情
func (s *NotebookService) GetByID(userID, notebookID uint64) (*model.Notebook, error) {
	notebook, err := s.notebookRepo.GetByIDAndUserID(notebookID, userID)
	if err != nil {
		return nil, err
	}
	if notebook == nil {
		return nil, ErrNotebookNotFound
	}

	count, err := s.noteRepo.CountByNotebookID(notebookID)
	if err != nil {
		return nil, err
	}
	notebook.NoteCount = count

	return notebook, nil
}

// Update 更新笔记本
func (s *NotebookService) Update(userID, notebookID uint64, req *model.NotebookUpdateReq) (*model.Notebook, error) {
	notebook, err := s.notebookRepo.GetByIDAndUserID(notebookID, userID)
	if err != nil {
		return nil, err
	}
	if notebook == nil {
		return nil, ErrNotebookNotFound
	}

	// 检查同名笔记本（排除自己）
	exists, err := s.notebookRepo.ExistsByUserIDAndName(userID, req.Name, notebookID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrNotebookNameDuplicate
	}

	notebook.Name = req.Name

	if err := s.notebookRepo.Update(notebook); err != nil {
		return nil, err
	}

	return notebook, nil
}

// Delete 删除笔记本
func (s *NotebookService) Delete(userID, notebookID uint64) error {
	// 步骤1: 权限校验
	notebook, err := s.notebookRepo.GetByIDAndUserID(notebookID, userID)
	if err != nil {
		return err
	}
	if notebook == nil {
		return ErrNotebookNotFound
	}

	// 步骤2: 禁止删除默认笔记本
	if notebook.IsDefault {
		return ErrCannotDeleteDefault
	}

	// 步骤3: 软删除笔记本下的所有笔记（移入回收站）
	if _, err := s.noteRepo.SoftDeleteByNotebookID(notebookID); err != nil {
		return err
	}

	// 步骤4: 删除笔记本
	return s.notebookRepo.Delete(notebookID)
}

// List 获取笔记本列表
func (s *NotebookService) List(userID uint64) ([]*model.Notebook, error) {
	notebooks, err := s.notebookRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, notebook := range notebooks {
		count, err := s.noteRepo.CountByNotebookID(notebook.ID)
		if err != nil {
			return nil, err
		}
		notebook.NoteCount = count
	}

	return notebooks, nil
}

// GetOrCreateDefault 获取或创建默认笔记本
func (s *NotebookService) GetOrCreateDefault(userID uint64) (*model.Notebook, error) {
	notebook, err := s.notebookRepo.GetOrCreateDefault(userID)
	if err != nil {
		return nil, err
	}
	count, err := s.noteRepo.CountByNotebookID(notebook.ID)
	if err != nil {
		return nil, err
	}
	notebook.NoteCount = count
	return notebook, nil
}
