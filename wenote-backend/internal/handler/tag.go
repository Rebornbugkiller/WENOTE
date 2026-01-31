package handler

import (
	"strconv"
	"wenote-backend/internal/model"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: service.NewTagService(),
	}
}

// Create 创建标签
func (h *TagHandler) Create(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.TagCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	tag, err := h.tagService.Create(userID, &req)
	if err != nil {
		if err == service.ErrTagNameExists {
			response.BadRequest(c, "标签名称已存在")
			return
		}
		response.InternalError(c, "创建标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "创建成功", tag)
}

// List 获取标签列表
func (h *TagHandler) List(c *gin.Context) {
	userID := c.GetUint64("userID")

	tags, err := h.tagService.List(userID)
	if err != nil {
		response.InternalError(c, "获取标签列表失败")
		return
	}

	response.Success(c, &model.TagListResp{List: tags})
}

// Update 更新标签
func (h *TagHandler) Update(c *gin.Context) {
	userID := c.GetUint64("userID")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	var req model.TagUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	tag, err := h.tagService.Update(userID, tagID, &req)
	if err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		if err == service.ErrTagNameExists {
			response.BadRequest(c, "标签名称已存在")
			return
		}
		response.InternalError(c, "更新标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", tag)
}

// Delete 删除标签
func (h *TagHandler) Delete(c *gin.Context) {
	userID := c.GetUint64("userID")
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}

	if err := h.tagService.Delete(userID, tagID); err != nil {
		if err == service.ErrTagNotFound {
			response.NotFound(c, "标签不存在")
			return
		}
		response.InternalError(c, "删除标签失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
