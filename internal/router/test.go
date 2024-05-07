package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
)

type TestRouter struct {
}

func NewTestRouter(engine *gin.Engine) *TestRouter {
	router := &TestRouter{}
	testGroup := engine.Group("/test")
	testGroup.GET("/hw", router.HandleHw)
	testGroup.PUT("/upload", router.HandleUpload)
	testGroup.POST("/changeRole", router.HandleChangeRole)
	return router
}

func (t *TestRouter) HandleHw(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

// HandleUpload Upload a picture
func (t *TestRouter) HandleUpload(ctx *gin.Context) {
	//url := ""
	//basicAuth := ""
	//imgFile := ctx.FormFile("img")
}

// HandleChangeRole 根据用户 id 更改权限
func (t *TestRouter) HandleChangeRole(ctx *gin.Context) {
	userId := ctx.Query("userId")
	role := ctx.PostForm("role")
	//fmt.Println("userId:", userId, "role:", role)
	if role == "" {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("role is required"))
		return
	}
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("userId is required"))
		return
	}
	user, err := dao.First[domain.User]("id = ?", userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	} else if user == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("user not found"))
		return
	}
	dao.Modify(user, "role", role)
}
