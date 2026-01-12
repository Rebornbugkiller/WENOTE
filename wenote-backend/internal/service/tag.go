package service

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"errors"
)

var (
	ErrTagNotFound   = errors.New("标签不存在")
	ErrTagNameExists = errors.New("标签名称已存在")
)

// TagService 标签服务
type TagService struct {
	tagRepo *repo.TagRepo
}

// NewTagService 创建标签服务实例
func NewTagService() *TagService {
	return &TagService{
		tagRepo: repo.NewTagRepo(),
	}
}

// Create 创建标签
func (s *TagService) Create(userID uint64, req *model.TagCreateReq) (*model.Tag, error) {
	// 检查名称是否已存在
	exists, err := s.tagRepo.ExistsByNameAndUserID(req.Name, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTagNameExists
	}

	color := req.Color
	if color == "" {
		color = "#6B7280"
	}

	tag := &model.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  color,
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetByID 获取标签详情
func (s *TagService) GetByID(userID, tagID uint64) (*model.Tag, error) {
	tag, err := s.tagRepo.GetByIDAndUserID(tagID, userID)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}

	// 获取笔记数量
	count, err := s.tagRepo.CountNotesByTagID(tagID)
	if err != nil {
		return nil, err
	}
	tag.NoteCount = count

	return tag, nil
}

// Delete 删除标签
func (s *TagService) Delete(userID, tagID uint64) error {
	tag, err := s.tagRepo.GetByIDAndUserID(tagID, userID)
	if err != nil {
		return err
	}
	if tag == nil {
		return ErrTagNotFound
	}

	return s.tagRepo.Delete(tagID)
}

// Update 更新标签
func (s *TagService) Update(userID, tagID uint64, req *model.TagUpdateReq) (*model.Tag, error) {
	tag, err := s.tagRepo.GetByIDAndUserID(tagID, userID)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}

	if req.Name != nil {
		// 检查新名称是否与其他标签重复
		exists, err := s.tagRepo.ExistsByNameAndUserID(*req.Name, userID)
		if err != nil {
			return nil, err
		}
		if exists && *req.Name != tag.Name {
			return nil, ErrTagNameExists
		}
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// List 获取标签列表
func (s *TagService) List(userID uint64) ([]*model.Tag, error) {
	tags, err := s.tagRepo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取每个标签的笔记数量
	for _, tag := range tags {
		count, err := s.tagRepo.CountNotesByTagID(tag.ID)
		if err != nil {
			return nil, err
		}
		tag.NoteCount = count
	}

	return tags, nil
}
