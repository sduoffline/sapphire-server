package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/data/cnt"
	"sapphire-server/pkg/util"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// 从 Authorization 头部中获取 JWT
		tokenString := strings.Split(authHeader, " ")[1]

		// 解析 JWT
		userId, err := util.ParseJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户ID保存到Header中
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), cnt.USER_ID, userId))

		// 继续处理请求
		c.Next()
	}
}
