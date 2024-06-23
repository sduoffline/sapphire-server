package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		customUserHeader := c.GetHeader("Sapphire-User-ID")
		// 如果没有自定义的Header，直接跳过
		if customUserHeader == "" {
			c.Next()
		}

		userIDStr := strings.Split(customUserHeader, " ")[1]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		c.Set("id", userID)

		// 继续处理请求
		c.Next()
	}
}
