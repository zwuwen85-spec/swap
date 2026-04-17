package middleware

import (
	"campus-swap-shop/pkg/jwt"
	"campus-swap-shop/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth JWT认证中间件
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Token
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			// 验证Bearer前缀
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		// 如果Header中没有，尝试从Query参数获取（主要用于WebSocket）
		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			response.Error(c, response.CodeNotLogin, "未登录")
			c.Abort()
			return
		}

		// 创建JWT管理器并解析Token
		jwtManager := jwt.NewJWTManager(jwtSecret)
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			response.Error(c, response.CodeNotLogin, "Token无效")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
