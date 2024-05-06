package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestRouter struct {
}

func NewTestRouter(engine *gin.Engine) *TestRouter {
	router := &TestRouter{}
	testGroup := engine.Group("/test")
	testGroup.GET("/hw", router.HandleHw)
	testGroup.PUT("/upload", router.HandleUpload)
	return router
}

func (t *TestRouter) HandleHw(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

// HandleUpload Upload a picture
func (t *TestRouter) HandleUpload(ctx *gin.Context) {
	url := ""
	basicAuth := ""
	imgFile := ctx.FormFile("img")
}
