package response

import (
	"net/http"

	"wenote-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess         = 0
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeTooManyRequests = 429
	CodeInternalError   = 500
)

var codeMessages = map[int]string{
	CodeSuccess:         "success",
	CodeBadRequest:      "请求参数错误",
	CodeUnauthorized:    "未授权，请先登录",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeTooManyRequests: "请求过于频繁",
	CodeInternalError:   "服务器内部错误",
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: codeMessages[CodeSuccess],
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	httpStatus := http.StatusOK
	if code == CodeUnauthorized {
		httpStatus = http.StatusUnauthorized
	} else if code == CodeForbidden {
		httpStatus = http.StatusForbidden
	} else if code == CodeNotFound {
		httpStatus = http.StatusNotFound
	} else if code == CodeTooManyRequests {
		httpStatus = http.StatusTooManyRequests
	} else if code == CodeInternalError {
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeBadRequest]
	}
	Fail(c, CodeBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeUnauthorized]
	}
	Fail(c, CodeUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeNotFound]
	}
	Fail(c, CodeNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeInternalError]
	}
	Fail(c, CodeInternalError, message)
}

func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = codeMessages[CodeTooManyRequests]
	}
	Fail(c, CodeTooManyRequests, message)
}

// ValidationError 处理验证错误，返回用户友好的提示
func ValidationError(c *gin.Context, err error) {
	msg := validator.TranslateValidationError(err)
	BadRequest(c, msg)
}
