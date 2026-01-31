package handler

import (
	"context"
	"strconv"
	"time"

	"wenote-backend/internal/model"
	"wenote-backend/internal/repo"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// NoteHandler 笔记处理器
type NoteHandler struct {
	noteService *service.NoteService
	auditRepo   *repo.AuditRepo
}

// NewNoteHandler 创建笔记处理器实例
func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		noteService: service.NewNoteService(),
		auditRepo:   repo.NewAuditRepo(),
	}
}

// Create 创建笔记
func (h *NoteHandler) Create(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.NoteCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	note, err := h.noteService.Create(userID, &req)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "创建笔记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", note)
}

// GetByID 获取笔记详情
func (h *NoteHandler) GetByID(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	note, err := h.noteService.GetByID(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "获取笔记失败")
		return
	}

	response.Success(c, note)
}

// Update 更新笔记
func (h *NoteHandler) Update(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	var req model.NoteUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	note, err := h.noteService.Update(userID, noteID, &req)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "更新笔记失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", note)
}

// Delete 删除笔记（软删除）
func (h *NoteHandler) Delete(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	if err := h.noteService.Delete(userID, noteID); err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "删除笔记失败")
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "delete",
		ResourceType: "note",
		ResourceID:   noteID,
		IPAddress:    c.ClientIP(),
	})

	response.SuccessWithMessage(c, "删除成功", nil)
}

// Restore 恢复已删除的笔记
func (h *NoteHandler) Restore(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	note, err := h.noteService.Restore(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "恢复笔记失败")
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "restore",
		ResourceType: "note",
		ResourceID:   noteID,
		IPAddress:    c.ClientIP(),
	})

	response.SuccessWithMessage(c, "恢复成功", note)
}

// List 获取笔记列表
func (h *NoteHandler) List(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.NoteListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	resp, err := h.noteService.List(userID, &req)
	if err != nil {
		response.InternalError(c, "获取笔记列表失败")
		return
	}

	response.Success(c, resp)
}

// UpdateTags 更新笔记标签
func (h *NoteHandler) UpdateTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	var req model.NoteTagsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	note, err := h.noteService.UpdateTags(userID, noteID, req.TagIDs)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.InternalError(c, "更新标签失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", note)
}

// ApplySuggestedTags 应用AI建议的标签
func (h *NoteHandler) ApplySuggestedTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	if err := h.noteService.ApplySuggestedTags(userID, noteID); err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "标签应用成功", nil)
}

// ListDeleted 获取回收站笔记列表
func (h *NoteHandler) ListDeleted(c *gin.Context) {
	userID := c.GetUint64("userID")

	// 获取分页参数
	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 {
			pageSize = val
		}
	}

	resp, err := h.noteService.ListDeleted(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, "获取回收站列表失败")
		return
	}

	response.Success(c, resp)
}

// BatchDelete 批量删除笔记
func (h *NoteHandler) BatchDelete(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	count, err := h.noteService.BatchHardDelete(req.NoteIDs, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_delete",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":      req.NoteIDs,
			"deleted_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量删除成功", map[string]interface{}{
		"deleted_count": count,
	})
}

// BatchRestore 批量恢复笔记
func (h *NoteHandler) BatchRestore(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchRestoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	count, err := h.noteService.BatchRestore(req.NoteIDs, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_restore",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":       req.NoteIDs,
			"restored_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量恢复成功", map[string]interface{}{
		"restored_count": count,
	})
}

// EmptyTrash 清空回收站
func (h *NoteHandler) EmptyTrash(c *gin.Context) {
	userID := c.GetUint64("userID")

	count, err := h.noteService.EmptyTrash(userID)
	if err != nil {
		response.InternalError(c, "清空回收站失败")
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "empty_trash",
		ResourceType: "note",
		Details: map[string]interface{}{
			"deleted_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "回收站已清空", map[string]interface{}{
		"deleted_count": count,
	})
}

// BatchMove 批量移动笔记
func (h *NoteHandler) BatchMove(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.BatchMoveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	count, err := h.noteService.BatchMove(req.NoteIDs, req.NotebookID, userID)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.BadRequest(c, "笔记本不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	_ = h.auditRepo.Create(&model.AuditLog{
		UserID:       userID,
		Action:       "batch_move",
		ResourceType: "note",
		Details: map[string]interface{}{
			"note_ids":    req.NoteIDs,
			"notebook_id": req.NotebookID,
			"moved_count": count,
		},
		IPAddress: c.ClientIP(),
	})

	response.SuccessWithMessage(c, "批量移动成功", map[string]interface{}{
		"moved_count": count,
	})
}

// GenerateSummaryAndTags 同步生成摘要和标签
func (h *NoteHandler) GenerateSummaryAndTags(c *gin.Context) {
	userID := c.GetUint64("userID")
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	result, err := h.noteService.GenerateSummaryAndTagsSync(userID, noteID)
	if err != nil {
		if err == service.ErrNoteNotFound {
			response.NotFound(c, "笔记不存在")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"summary": result.Summary,
		"tags":    result.Tags,
	})
}

// AIAssist AI写作助手
// POST /api/v1/notes/ai/assist
func (h *NoteHandler) AIAssist(c *gin.Context) {
	var req model.AIAssistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	// 创建超时上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	result, err := h.noteService.AIAssist(ctx, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}
