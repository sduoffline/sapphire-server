package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/data/dto"
)

type TestRouter struct {
}

func NewTestRouter(engine *gin.Engine) *TestRouter {
	router := &TestRouter{}
	testGroup := engine.Group("/test")
	testGroup.GET("/hw", router.HandleHelloWorld)
	return router
}

func (t *TestRouter) HandleHelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse("Hello World"))
}
