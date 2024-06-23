package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	docs "sapphire-server/cmd/docs"
	"sapphire-server/internal/conf"
	"sapphire-server/internal/infra"
	"sapphire-server/internal/middleware"
	"sapphire-server/internal/router"
)

import "github.com/swaggo/files" // swagger embed files

// @title			Sapphire Server API
// @version		1.0
// @description	Sapphire is a platform for image annotation and dataset management.
// @contact.name	API Support
// @contact.url	https://www.example.com/support
// @license.name	Apache 2.0
// @license.url	https://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	var err error
	// 初始化并读取配置
	conf.InitConfig()
	slog.Info("Server started")

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
	docs.SwaggerInfo.BasePath = "/api/v1"
	// 调用 internal/router/user.go 中的 NewUserRouter 方法
	router.NewUserRouter(engine)
	router.NewTaskRouter(engine)
	router.NewImgRouter(engine)
	router.NewDatasetRouter(engine)
	router.NewAnnotationRouter(engine)
	router.NewTestRouter(engine)
	router.NewScoreRouter(engine)
	router.NewMessageRouter(engine)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// start http server
	err = engine.Run(conf.GetServerAddr())
	if err != nil {
		panic(err)
	}
}
