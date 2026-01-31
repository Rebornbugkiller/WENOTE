package handler

import (
	"strconv"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AttachmentHandler 附件处理器
type AttachmentHandler struct {
	service *service.AttachmentService
}

// NewAttachmentHandler 创建附件处理器实例
func NewAttachmentHandler() *AttachmentHandler {
	return &AttachmentHandler{
		service: service.NewAttachmentService(),
	}
}

// UploadImage 上传图片
// POST /api/v1/notes/:id/attachments
func (h *AttachmentHandler) UploadImage(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")

	// 获取笔记ID
	noteIDStr := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "未找到上传文件")
		return
	}

	// 上传文件
	result, err := h.service.UploadImage(userID.(uint64), noteID, file)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetAttachments 获取笔记的附件列表
// GET /api/v1/notes/:id/attachments
func (h *AttachmentHandler) GetAttachments(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")

	// 获取笔记ID
	noteIDStr := c.Param("id")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的笔记ID")
		return
	}

	// 获取附件列表
	attachments, err := h.service.GetAttachments(userID.(uint64), noteID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"list": attachments})
}

// DeleteAttachment 删除附件
// DELETE /api/v1/attachments/:id
func (h *AttachmentHandler) DeleteAttachment(c *gin.Context) {
	// 获取用户ID
	userID, _ := c.Get("userID")

	// 获取附件ID
	attachmentIDStr := c.Param("id")
	attachmentID, err := strconv.ParseUint(attachmentIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的附件ID")
		return
	}

	// 删除附件
	if err := h.service.DeleteAttachment(userID.(uint64), attachmentID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

