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

// GetMe 获取当前用户信息
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetUint64("userID")

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "获取用户信息失败")
		return
	}

	resp := &model.UserResp{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	response.Success(c, resp)
}
