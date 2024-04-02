package router

import (
	"github.com/gin-gonic/gin"
	"log"

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
	engine.POST("/register", router.HandleRegister)
	return router
}

// HandleRegister 注册
func (u *UserRouter) HandleRegister(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := &dto.Register{}
	if err := ctx.BindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(body)

	// 由于功能简单，直接调用 domain 层的方法
	user := domain.NewUser()
	if err := user.Register(*body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 复杂的话在 service 层处理

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
