package middleware

import (
	"wenote-backend/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
// 记录每个 HTTP 请求的详细信息，用于监控和调试
//
// 记录的信息包括：
//   - method: HTTP 方法（GET、POST 等）
//   - path: 请求路径（包含查询参数）
//   - status: HTTP 状态码
//   - duration_ms: 请求处理耗时（毫秒）
//   - ip: 客户端 IP
//   - user_id: 用户 ID（已登录时）
//   - error: 错误信息（如果有）
//
// 特殊处理：
//   - 慢请求（>500ms）使用 Warn 级别记录
//   - 正常请求使用 Info 级别记录
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 执行后续处理（Handler 和其他中间件）
		c.Next()

		// 计算请求处理耗时
		latency := time.Since(start)
		status := c.Writer.Status()

		// 拼接完整路径（包含查询参数）
		if query != "" {
			path = path + "?" + query
		}

		// 获取用户 ID（如果已登录）
		userID := GetUserID(c)

		// 收集错误信息
		// Gin 的 c.Errors 包含处理过程中添加的所有错误
		var errMsg string
		if len(c.Errors) > 0 {
			errMsgs := make([]string, len(c.Errors))
			for i, e := range c.Errors {
				errMsgs[i] = e.Error()
			}
			errMsg = strings.Join(errMsgs, "; ")
		}

		// 构建日志字段
		// 使用 key-value 格式，便于日志系统解析
		fields := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration_ms", latency.Milliseconds(),
			"ip", c.ClientIP(),
		}

		// 只有已登录用户才记录 user_id
		if userID > 0 {
			fields = append(fields, "user_id", userID)
		}

		// 只有发生错误时才记录 error
		if errMsg != "" {
			fields = append(fields, "error", errMsg)
		}

		// 根据请求耗时选择日志级别
		// 慢请求（>500ms）使用 Warn 级别，便于后续排查性能问题
		if latency > 500*time.Millisecond {
			logger.Warn("slow request", fields...)
		} else {
			logger.Info("request", fields...)
		}
	}
}
