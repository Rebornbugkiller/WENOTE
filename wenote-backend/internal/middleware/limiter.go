package middleware

import (
	"sync"
	"wenote-backend/config"
	"wenote-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流中间件
func RateLimiter() gin.HandlerFunc {
	cfg := config.GlobalConfig.RateLimit

	globalLimiter := rate.NewLimiter(
		rate.Limit(cfg.GlobalRate),
		int(cfg.GlobalBurst),
	)

	userLimiters := &sync.Map{}
	ipLimiters := &sync.Map{}

	return func(c *gin.Context) {
		if !globalLimiter.Allow() {
			response.TooManyRequests(c, "服务繁忙，请稍后重试")
			c.Abort()
			return
		}

		clientIP := c.ClientIP()
		ipLimiter := getIPLimiter(ipLimiters, clientIP, cfg)
		if !ipLimiter.Allow() {
			response.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		if userID, exists := c.Get("userID"); exists {
			limiter := getUserLimiter(userLimiters, userID.(uint64), cfg)
			if !limiter.Allow() {
				response.TooManyRequests(c, "请求过于频繁，请稍后再试")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// getUserLimiter 获取或创建用户限流器
func getUserLimiter(userLimiters *sync.Map, userID uint64, cfg config.RateLimitConfig) *rate.Limiter {
	if limiter, ok := userLimiters.Load(userID); ok {
		return limiter.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(
		rate.Limit(cfg.UserRate),
		int(cfg.UserBurst),
	)

	actual, _ := userLimiters.LoadOrStore(userID, limiter)
	return actual.(*rate.Limiter)
}

// getIPLimiter 获取或创建 IP 限流器
func getIPLimiter(ipLimiters *sync.Map, ip string, cfg config.RateLimitConfig) *rate.Limiter {
	if limiter, ok := ipLimiters.Load(ip); ok {
		return limiter.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(
		rate.Limit(cfg.IPRate),
		int(cfg.IPBurst),
	)

	actual, _ := ipLimiters.LoadOrStore(ip, limiter)
	return actual.(*rate.Limiter)
}
