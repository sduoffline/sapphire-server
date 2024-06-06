package middleware

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/pkg/util"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 从 Authorization 头部中获取 JWT
		tokenString := strings.Split(authHeader, " ")[1]

		// 解析 JWT
		userId, err := util.ParseJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "请重新登录"})
			c.Abort()
			return
		}

		// 将用户ID保存到Header中
		c.Set("id", userId)

		// 继续处理请求
		c.Next()
	}
}
