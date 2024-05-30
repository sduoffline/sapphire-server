package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
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
	userGroup.POST("/change-role", router.HandleChangeRole)
	userGroup.GET("/profile", router.HandleProfile)

	statisticGroup := userGroup.Group("/statistic")
	statisticGroup.GET("/credit", router.HandleCredit)
	return router
}

// HandleProfile 获取用户信息
func (u *UserRouter) HandleProfile(ctx *gin.Context) {
	userId := ctx.Query("userId")
	user, err := dao.First[domain.User]("id = ?", userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	} else if user == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("user not found"))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(user))
}

// HandleRegister 注册
func (u *UserRouter) HandleRegister(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := dto.Register{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}
	// 由于功能简单，直接调用 domain 的方法
	user := domain.NewUser()
	token, user, err := user.Register(body)
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
	token, user, err := user.Login(*body)
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
		"user":  user,
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(payload))
}

// HandleChangeRole 根据用户 id 更改权限
func (u *UserRouter) HandleChangeRole(ctx *gin.Context) {
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
	err = dao.Modify(user, "role", role)
	if err != nil {
		return
	}
}

// HandleCredit 获取用户积分
func (u *UserRouter) HandleCredit(ctx *gin.Context) {
	userId := ctx.Query("userId")
	user, err := dao.First[domain.User]("id = ?", userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	} else if user == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("user not found"))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(user.Score))
}
