package handler

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NotebookHandler 笔记本处理器
type NotebookHandler struct {
	notebookService *service.NotebookService
}

// NewNotebookHandler 创建笔记本处理器实例
func NewNotebookHandler() *NotebookHandler {
	return &NotebookHandler{
		notebookService: service.NewNotebookService(),
	}
}

// Create 创建笔记本
func (h *NotebookHandler) Create(c *gin.Context) {
	userID := c.GetUint64("userID")
	var req model.NotebookCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	notebook, err := h.notebookService.Create(userID, &req)
	if err != nil {
		response.InternalError(c, "创建笔记本失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", notebook)
}

// List 获取笔记本列表
func (h *NotebookHandler) List(c *gin.Context) {
	userID := c.GetUint64("userID")
	notebooks, err := h.notebookService.List(userID)
	if err != nil {
		response.InternalError(c, "获取笔记本列表失败")
		return
	}
	response.Success(c, &model.NotebookListResp{List: notebooks})
}

// GetByID 获取笔记本详情
func (h *NotebookHandler) GetByID(c *gin.Context) {
	userID := c.GetUint64("userID")
	notebookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记本ID")
		return
	}

	notebook, err := h.notebookService.GetByID(userID, notebookID)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.NotFound(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "获取笔记本失败")
		return
	}

	response.Success(c, notebook)
}

// Update 更新笔记本
func (h *NotebookHandler) Update(c *gin.Context) {
	userID := c.GetUint64("userID")
	notebookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记本ID")
		return
	}

	var req model.NotebookUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	notebook, err := h.notebookService.Update(userID, notebookID, &req)
	if err != nil {
		if err == service.ErrNotebookNotFound {
			response.NotFound(c, "笔记本不存在")
			return
		}
		response.InternalError(c, "更新笔记本失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", notebook)
}

// Delete 删除笔记本
func (h *NotebookHandler) Delete(c *gin.Context) {
	userID := c.GetUint64("userID")
	notebookID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记本ID")
		return
	}

	if err := h.notebookService.Delete(userID, notebookID); err != nil {
		if err == service.ErrNotebookNotFound {
			response.NotFound(c, "笔记本不存在")
			return
		}
		if err == service.ErrCannotDeleteDefault {
			response.BadRequest(c, "默认笔记本不能删除")
			return
		}
		response.InternalError(c, "删除笔记本失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetDefault 获取或创建默认笔记本
func (h *NotebookHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint64("userID")
	notebook, err := h.notebookService.GetOrCreateDefault(userID)
	if err != nil {
		response.InternalError(c, "获取默认笔记本失败")
		return
	}
	response.Success(c, notebook)
}
