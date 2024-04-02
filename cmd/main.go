package main

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/router"
)

func main() {
	engine := gin.Default()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	// 调用 internal/router/user.go 中的 NewUserRouter 方法
	router.NewUserRouter(engine)
	engine.Run(":8080")
}
