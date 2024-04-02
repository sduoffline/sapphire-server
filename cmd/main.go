package main

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/conf"
	"sapphire-server/internal/infra"
	"sapphire-server/internal/router"
)

func main() {
	var err error
	// 初始化并读取配置
	conf.InitConfig()

	// 连接数据库
	err = infra.InitDB()
	if err != nil {
		panic(err)
	}

	// init gin
	engine := gin.Default()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	// init http routes
	// 调用 internal/router/user.go 中的 NewUserRouter 方法
	router.NewUserRouter(engine)

	// start http server
	err = engine.Run(conf.GetServerAddr())
	if err != nil {
		panic(err)
	}
}
