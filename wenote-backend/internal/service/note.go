package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"wenote-backend/pkg/ai"
)

var (
	ErrNoteNotFound = errors.New("笔记不存在")
)

// 全局依赖(由 main.go 初始化)
var globalAIClient ai.Client

// InitGlobalDeps 初始化全局依赖
func InitGlobalDeps(client ai.Client) {
	globalAIClient = client
}

// NoteService 笔记服务
type NoteService struct {
	noteRepo     *repo.NoteRepo
	notebookRepo *repo.NotebookRepo
	tagRepo      *repo.TagRepo
}

// NewNoteService 创建笔记服务实例
func NewNoteService() *NoteService {
	return &NoteService{
		noteRepo:     repo.NewNoteRepo(),
		notebookRepo: repo.NewNotebookRepo(),
		tagRepo:      repo.NewTagRepo(),
	}
}

// Create 创建笔记
func (s *NoteService) Create(userID uint64, req *model.NoteCreateReq) (*model.Note, error) {
	// 目的：确保笔记创建时指定的笔记本存在，且该笔记本归属当前用户，防止用户在不存在或非本人拥有的笔记本下创建笔记。
	notebook, err := s.notebookRepo.GetByIDAndUserID(req.NotebookID, userID)
	if err != nil {
		return nil, err
	}
	if notebook == nil {
		return nil, ErrNotebookNotFound
	}

	// 设置默认摘要长度
	summaryLen := req.SummaryLen
	if summaryLen <= 0 {
		summaryLen = 200
	}
	if summaryLen > 500 {
		summaryLen = 500
	}

	// 创建笔记
	note := &model.Note{
		UserID:     userID,
		NotebookID: req.NotebookID,
		Title:      req.Title,
		Content:    req.Content,
		SummaryLen: summaryLen,
		AIStatus:   model.AIStatusPending,
	}

	if err := s.noteRepo.Create(note); err != nil {
		return nil, err
	}

	// 如果有标签，更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.noteRepo.ReplaceNoteTags(note.ID, req.TagIDs); err != nil {
			return nil, err
		}
	}

	return note, nil
}

// GetByID 获取笔记详情
func (s *NoteService) GetByID(userID, noteID uint64) (*model.Note, error) {
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}
	return note, nil
}

// Update 更新笔记
func (s *NoteService) Update(userID, noteID uint64, req *model.NoteUpdateReq) (*model.Note, error) {
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	// 如果要更换笔记本，先验证要更换到的目标笔记本是否存在且归属于当前用户
	if req.NotebookID != nil {
		// 通过notebookRepo查找新笔记本，传入NotebookID和userID防止越权
		notebook, err := s.notebookRepo.GetByIDAndUserID(*req.NotebookID, userID)
		if err != nil {
			// 如果查找出现数据库等错误则直接返回
			return nil, err
		}
		if notebook == nil {
			// 新笔记本不存在（或不是自己的），返回“笔记本不存在”业务错误
			return nil, ErrNotebookNotFound
		}
		// 校验通过，允许更换笔记本
		note.NotebookID = *req.NotebookID
	}

	// 讲讲每一步在做什么：
	// 1. 标记内容是否发生变化（影响 updated_at 的更新）
	contentChanged := false
	statusOnlyChange := false // 标记是否只改状态字段

	// 2. 如果标题有变，更新并标记内容变化（这样AI摘要等会重算）
	if req.Title != nil {
		note.Title = *req.Title
		contentChanged = true
	}

	// 3. 如果正文有变，更新并标记内容变化
	if req.Content != nil {
		note.Content = *req.Content
		contentChanged = true
	}

	// 3.5. 如果笔记本变化，也算内容变化
	if req.NotebookID != nil {
		contentChanged = true
	}

	// 4. 可选地更新AI摘要长度，只允许1-500
	if req.SummaryLen != nil {
		summaryLen := *req.SummaryLen
		if summaryLen > 0 && summaryLen <= 500 {
			note.SummaryLen = summaryLen
		}
	}

	// 5. 检查是否只更新状态字段（置顶、收藏）
	if !contentChanged && (req.IsPinned != nil || req.IsStarred != nil) {
		statusOnlyChange = true
	}

	// 6. 置顶与星标字段可以直接改
	if req.IsPinned != nil {
		note.IsPinned = *req.IsPinned
	}
	if req.IsStarred != nil {
		note.IsStarred = *req.IsStarred
	}

	// 7. 保持 AI 摘要和标签不变（每个笔记只能生成一次 AI，生成后永久保存）
	// 确保 ai_status 字段始终有效（防止 NULL 或空字符串导致数据库错误）
	if note.AIStatus == "" {
		note.AIStatus = model.AIStatusPending
	}

	// 8. 写回数据库（只更新需要的字段，不覆盖 AI 相关字段）
	if statusOnlyChange {
		fields := make(map[string]interface{})
		if req.IsPinned != nil {
			fields["is_pinned"] = note.IsPinned
		}
		if req.IsStarred != nil {
			fields["is_starred"] = note.IsStarred
		}
		if err := s.noteRepo.UpdateFieldsWithoutTime(note.ID, fields); err != nil {
			return nil, err
		}
	} else {
		// 只更新用户可编辑的字段，不覆盖 AI 字段
		fields := map[string]interface{}{
			"title":       note.Title,
			"content":     note.Content,
			"notebook_id": note.NotebookID,
			"is_pinned":   note.IsPinned,
			"is_starred":  note.IsStarred,
			"summary_len": note.SummaryLen,
		}
		if err := s.noteRepo.UpdateFields(note.ID, fields); err != nil {
			return nil, err
		}
	}

	// 9. 如有标签变动，更新标签（note_tags 关联表）
	if req.TagIDs != nil {
		if err := s.noteRepo.ReplaceNoteTags(note.ID, req.TagIDs); err != nil {
			return nil, err
		}
	}

	// 10. 返回这条笔记的完整信息（含最新标签/AI字段等）
	// 是的，这里已经更新到数据库，
	// 然后通过ID再次查询最新的笔记返回前端
	return s.noteRepo.GetByID(note.ID)
}

// Delete 软删除笔记
func (s *NoteService) Delete(userID, noteID uint64) error {
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return err
	}
	if note == nil {
		return ErrNoteNotFound
	}

	return s.noteRepo.SoftDelete(noteID)
}

// Restore 恢复已删除的笔记
func (s *NoteService) Restore(userID, noteID uint64) (*model.Note, error) {
	note, err := s.noteRepo.GetDeletedByIDAndUserID(noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	if err := s.noteRepo.Restore(noteID); err != nil {
		return nil, err
	}

	return s.noteRepo.GetByID(noteID)
}

// List 获取笔记列表
func (s *NoteService) List(userID uint64, req *model.NoteListReq) (*model.NoteListResp, error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	notes, total, err := s.noteRepo.List(userID, req)
	if err != nil {
		return nil, err
	}

	return &model.NoteListResp{
		Total: total,
		List:  notes,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// UpdateTags 更新笔记标签
func (s *NoteService) UpdateTags(userID, noteID uint64, tagIDs []uint64) (*model.Note, error) {
	// 1. 验证笔记归属
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil || note.DeletedAt != nil {
		return nil, ErrNoteNotFound
	}

	// 2. 验证所有标签都属于该用户
	if len(tagIDs) > 0 {
		tags, err := s.tagRepo.ListByIDs(tagIDs)
		if err != nil {
			return nil, err
		}
		if len(tags) != len(tagIDs) {
			return nil, errors.New("部分标签不存在")
		}
		for _, tag := range tags {
			if tag.UserID != userID {
				return nil, errors.New("无权使用该标签")
			}
		}
	}

	// 3. 替换标签关联
	if err := s.noteRepo.ReplaceNoteTags(noteID, tagIDs); err != nil {
		return nil, err
	}

	// 4. 返回更新后的笔记
	return s.GetByID(userID, noteID)
}

// ApplySuggestedTags 应用AI建议的标签
func (s *NoteService) ApplySuggestedTags(userID, noteID uint64) error {
	// 1. 获取笔记及其 suggested_tags
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return err
	}
	if note == nil {
		return ErrNoteNotFound
	}

	if len(note.SuggestedTags) == 0 {
		return errors.New("暂无标签建议")
	}

	// 2. 为每个建议标签创建或获取 Tag
	var tagIDs []uint64
	for _, tagName := range note.SuggestedTags {
		tag, err := s.tagRepo.GetOrCreate(userID, tagName)
		if err != nil {
			continue
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	// 3. 关联标签到笔记
	if err := s.noteRepo.ReplaceNoteTags(noteID, tagIDs); err != nil {
		return err
	}

	// 4. 清空 suggested_tags 避免重复应用
	return s.noteRepo.ClearSuggestedTags(noteID)
}

// BatchHardDelete 批量永久删除笔记（用于回收站）
func (s *NoteService) BatchHardDelete(noteIDs []uint64, userID uint64) (int64, error) {
	validNoteIDs, err := s.noteRepo.FilterByUserID(noteIDs, userID)
	if err != nil {
		return 0, err
	}

	if len(validNoteIDs) == 0 {
		return 0, errors.New("无有效笔记可删除")
	}

	return s.noteRepo.BatchHardDelete(validNoteIDs)
}

// BatchRestore 批量恢复笔记
func (s *NoteService) BatchRestore(noteIDs []uint64, userID uint64) (int64, error) {
	// 权限校验
	validNoteIDs, err := s.noteRepo.FilterByUserID(noteIDs, userID)
	if err != nil {
		return 0, err
	}

	if len(validNoteIDs) == 0 {
		return 0, errors.New("无有效笔记可恢复")
	}

	return s.noteRepo.BatchRestore(validNoteIDs)
}

// BatchMove 批量移动笔记到指定笔记本
func (s *NoteService) BatchMove(noteIDs []uint64, notebookID, userID uint64) (int64, error) {
	// 验证笔记本是否存在且属于该用户
	notebook, err := s.notebookRepo.GetByIDAndUserID(notebookID, userID)
	if err != nil {
		return 0, err
	}
	if notebook == nil {
		return 0, ErrNotebookNotFound
	}

	// 权限校验
	validNoteIDs, err := s.noteRepo.FilterByUserID(noteIDs, userID)
	if err != nil {
		return 0, err
	}

	if len(validNoteIDs) == 0 {
		return 0, errors.New("无有效笔记可移动")
	}

	return s.noteRepo.BatchUpdateNotebook(validNoteIDs, notebookID)
}

// ListDeleted 获取回收站笔记列表
func (s *NoteService) ListDeleted(userID uint64, page, pageSize int) (*model.NoteListResp, error) {
	// 设置默认分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	notes, total, err := s.noteRepo.ListDeleted(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &model.NoteListResp{
		Total: total,
		List:  notes,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// CleanupDeletedNotes 清理已删除的笔记
func (s *NoteService) CleanupDeletedNotes(days int) (int64, error) {
	slog.Info("开始清理过期笔记", "days", days)

	count, err := s.noteRepo.HardDeleteOld(days)
	if err != nil {
		slog.Error("清理失败", "error", err)
		return 0, err
	}

	slog.Info("清理完成", "deleted_count", count)
	return count, nil
}

// GenerateSummaryAndTagsSync 同步生成摘要和标签建议
func (s *NoteService) GenerateSummaryAndTagsSync(userID, noteID uint64) (*ai.SummaryResult, error) {
	// 验证笔记归属
	note, err := s.noteRepo.GetByIDAndUserID(noteID, userID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, ErrNoteNotFound
	}

	// 检查是否已生成过（每个笔记只能生成一次）
	if note.AIStatus == model.AIStatusDone {
		return nil, errors.New("该笔记已生成过AI摘要，不允许重复生成")
	}

	// 检查内容是否为空
	if note.Content == "" {
		return nil, errors.New("笔记内容为空")
	}

	// 检查 AI 客户端是否初始化
	if globalAIClient == nil {
		return nil, errors.New("AI 服务未初始化")
	}

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 同步调用 AI 生成摘要和标签
	result, err := globalAIClient.GenerateSummaryAndTags(ctx, note.Content, note.SummaryLen)
	if err != nil {
		slog.Error("AI generate summary and tags failed", "note_id", noteID, "error", err)
		return nil, fmt.Errorf("AI 生成失败: %w", err)
	}

	// 更新到数据库
	err = s.noteRepo.UpdateAIResult(noteID, result.Summary, result.Tags)
	if err != nil {
		slog.Error("Failed to update AI result", "note_id", noteID, "error", err)
		return nil, fmt.Errorf("保存失败: %w", err)
	}

	return result, nil
}

// AIAssist AI写作助手
func (s *NoteService) AIAssist(ctx context.Context, req *model.AIAssistReq) (*model.AIAssistResp, error) {
	// 检查 AI 客户端是否初始化
	if globalAIClient == nil {
		return nil, errors.New("AI 服务未初始化")
	}

	var prompt string
	var userInput string

	switch req.Action {
	case "continue":
		// 续写：基于上下文生成后续内容
		prompt = "你是一个专业的写作助手。请根据用户提供的上下文内容，自然地续写接下来的内容。续写内容应该流畅、连贯，风格与上文保持一致。请直接输出续写内容，不要包含任何解释性文字。"
		userInput = fmt.Sprintf("上下文：\n%s\n\n请续写接下来的内容（约100-200字）：", req.Context)

	case "rewrite":
		// 改写/润色：优化选中的文本
		if req.SelectedText == "" {
			return nil, errors.New("改写操作需要提供选中的文本")
		}
		prompt = "你是一个专业的文字编辑。请对用户提供的文本进行改写和润色，使其更加流畅、准确、优雅。保持原意不变，但提升表达质量。请直接输出改写后的文本，不要包含任何解释性文字。"
		userInput = fmt.Sprintf("原文：\n%s\n\n请改写润色：", req.SelectedText)

	case "expand":
		// 扩写：将简短要点扩展为完整段落
		if req.SelectedText == "" {
			return nil, errors.New("扩写操作需要提供选中的文本")
		}
		prompt = "你是一个专业的写作助手。请将用户提供的简短内容或要点扩展为详细、完整的段落。扩写应该内容充实、逻辑清晰、表达流畅。请直接输出扩写后的内容，不要包含任何解释性文字。"
		userInput = fmt.Sprintf("简短内容：\n%s\n\n请扩写为详细内容：", req.SelectedText)

	case "translate":
		// 翻译：中英文互译
		if req.SelectedText == "" {
			return nil, errors.New("翻译操作需要提供选中的文本")
		}
		targetLang := "英文"
		if req.Language == "zh" {
			targetLang = "中文"
		}
		prompt = fmt.Sprintf("你是一个专业的翻译助手。请将用户提供的文本翻译为%s。翻译应该准确、流畅、地道。请直接输出翻译结果，不要包含任何解释性文字。", targetLang)
		userInput = fmt.Sprintf("原文：\n%s\n\n请翻译为%s：", req.SelectedText, targetLang)

	default:
		return nil, errors.New("不支持的操作类型")
	}

	// 调用 AI 生成
	result, err := globalAIClient.GenerateText(ctx, prompt, userInput)
	if err != nil {
		slog.Error("AI assist failed", "action", req.Action, "error", err)
		return nil, fmt.Errorf("AI 生成失败: %w", err)
	}

	return &model.AIAssistResp{
		Result: result,
	}, nil
}
