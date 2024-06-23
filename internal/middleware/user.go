package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
)

func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		customUserHeader := c.GetHeader("Sapphire-User-ID")
		// 如果没有自定义的Header，直接跳过
		if customUserHeader == "" {
			c.Next()
			return
		}
		slog.Info(customUserHeader)

		userIDStr := customUserHeader
		// 读取为uint类型
		userID64, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		userID := uint(userID64)
		slog.Info("Read ID from header: ", userID)

		_, exists := c.Get("id")
		if exists {
			slog.Info("User ID already exists, Update it: ", userID)
			c.Set("id", userID)
		} else {
			slog.Info("User ID not exists, Set it: ", userID)
			c.Set("id", userID)
		}

		curID, _ := c.Get("id")
		slog.Info("Current User ID: ", curID)

		// 继续处理请求
		c.Next()
	}
}
