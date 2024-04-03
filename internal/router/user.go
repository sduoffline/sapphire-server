package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
)

// UserRouter 用户路由
type UserRouter struct {
}

// NewUserRouter 创建用户路由
func NewUserRouter(engine *gin.Engine) *UserRouter {
	router := &UserRouter{}
	// 创建一个路由组
	userGroup := engine.Group("/user")
	userGroup.POST("/register", router.HandleRegister)
	userGroup.POST("/login", router.HandleLogin)
	return router
}

// HandleRegister 注册
func (u *UserRouter) HandleRegister(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := &dto.Register{}
	if err := ctx.BindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	// 由于功能简单，直接调用 domain 的方法
	user := domain.NewUser()
	token, err := user.Register(*body)
	if err != nil {
		if err.Error() == "existed user" {
			ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("existed user"))
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		}
		return
	}
	// 复杂的话在 service 层处理

	payload := map[string]interface{}{
		"token": token,
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(payload))
}

func (u *UserRouter) HandleLogin(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := &dto.Login{}
	if err := ctx.BindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	// 由于功能简单，直接调用 domain 层的方法
	user := domain.NewUser()
	token, err := user.Login(*body)
	if err != nil {
		if err.Error() == "wrong password" {
			ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("wrong password"))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
			return
		}
	}
	// 复杂的话在 service 层处理

	payload := map[string]interface{}{
		"token": token,
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(payload))
}
