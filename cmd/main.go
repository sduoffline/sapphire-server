package main

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/conf"
	"sapphire-server/internal/infra"
	"sapphire-server/internal/middleware"
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
	err = infra.InitRedis()
	if err != nil {
		panic(err)
	}

	// init gin
	engine := gin.Default()
	//engine.Use(middleware.AuthMiddleware())
	// init http routes
	engine.Use(middleware.Cors())

	// 调用 internal/router/user.go 中的 NewUserRouter 方法
	router.NewUserRouter(engine)
	router.NewTaskRouter(engine)
	router.NewImgRouter(engine)
	router.NewDatasetRouter(engine)
	router.NewAnnotationRouter(engine)
	router.NewTestRouter(engine)

	// start http server
	err = engine.Run(conf.GetServerAddr())
	if err != nil {
		panic(err)
	}
}
