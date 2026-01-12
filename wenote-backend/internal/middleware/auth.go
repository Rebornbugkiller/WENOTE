package middleware

import (
	"wenote-backend/config"
	"wenote-backend/pkg/jwt"
	"wenote-backend/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	cfg := config.GlobalConfig.JWT
	jwtManager := jwt.NewJWTManager(cfg.Secret, cfg.Expire)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "Token格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				response.Unauthorized(c, "Token已过期，请重新登录")
			case jwt.ErrTokenMalformed:
				response.Unauthorized(c, "Token格式错误")
			default:
				response.Unauthorized(c, "无效的Token")
			}
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

func GetUserID(c *gin.Context) uint64 {
	if userID, exists := c.Get("userID"); exists {
		return userID.(uint64)
	}
	return 0
}
