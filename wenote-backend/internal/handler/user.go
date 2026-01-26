package handler

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// GetMe 获取当前用户信息（含统计）
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint64("userID")

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "获取用户信息失败")
		return
	}

	response.Success(c, profile)
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	profile, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.NotFound(c, "用户不存在")
		case service.ErrEmailExists:
			response.BadRequest(c, "邮箱已被使用")
		case service.ErrInvalidAvatar:
			response.BadRequest(c, "无效的头像样式或颜色")
		case service.ErrInvalidEmail:
			response.BadRequest(c, "邮箱格式不正确")
		default:
			response.InternalError(c, "更新失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "个人资料已更新", profile)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err := h.userService.ChangePassword(userID, &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.NotFound(c, "用户不存在")
		case service.ErrPasswordIncorrect:
			response.BadRequest(c, "当前密码错误")
		case service.ErrPasswordSame:
			response.BadRequest(c, "新密码不能与旧密码相同")
		default:
			response.InternalError(c, "修改密码失败")
		}
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// DeleteAccount 注销账号
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint64("userID")

	var req model.DeleteAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err := h.userService.DeleteAccount(userID, &req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.NotFound(c, "用户不存在")
		case service.ErrPasswordIncorrect:
			response.BadRequest(c, "密码错误")
		case service.ErrConfirmMismatch:
			response.BadRequest(c, "请输入 DELETE 确认删除")
		default:
			response.InternalError(c, "删除账号失败")
		}
		return
	}

	response.SuccessWithMessage(c, "账号已删除", nil)
}
