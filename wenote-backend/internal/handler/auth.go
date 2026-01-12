package handler

import (
	"wenote-backend/internal/model"
	"wenote-backend/internal/service"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if err == service.ErrUsernameExists {
			response.BadRequest(c, "用户名已存在")
			return
		}
		response.InternalError(c, "注册失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "注册成功", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			response.BadRequest(c, "用户不存在")
		case service.ErrPasswordIncorrect:
			response.BadRequest(c, "密码错误")
		default:
			response.InternalError(c, "登录失败: "+err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "登录成功", resp)
}
