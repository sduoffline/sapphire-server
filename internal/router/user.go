package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
	"strconv"
)

// UserRouter 用户路由
type UserRouter struct {
}

var userDomain = domain.NewUser()

// NewUserRouter 创建用户路由
func NewUserRouter(engine *gin.Engine) *UserRouter {
	router := &UserRouter{}
	// 创建一个路由组
	userGroup := engine.Group("/user")
	userGroup.POST("/register", router.HandleRegister)
	userGroup.POST("/login", router.HandleLogin)
	userGroup.POST("/change-role", router.HandleChangeRole)
	userGroup.GET("/profile/:id", router.HandleProfile)

	authGroup := userGroup.Group("").Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/passwd/change", router.HandleChangePasswd)
		authGroup.POST("/info/change", router.HandleChangeInfo)
	}

	statisticGroup := userGroup.Group("/statistic")
	statisticGroup.GET("/credit", router.HandleCredit)
	return router
}

// HandleProfile godoc
//
//	@Summary		获取用户信息
//	@Description	根据用户 id 获取用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		string	true	"User ID"
//	@Success		200		{object}	dto.Response{data=domain.User}
//	@Router			/user/profile [get]
func (u *UserRouter) HandleProfile(ctx *gin.Context) {
	var err error
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	user, err := userDomain.GetUserInfo(uint(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	} else if user == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("user not found"))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(user))
}

// HandleRegister godoc
//
//	@Summary		注册
//	@Description	用户注册
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.Register	true	"Register"
//	@Success		200		{object}	dto.Response{data=map[string]interface{}}
//	@Router			/user/register [post]
func (u *UserRouter) HandleRegister(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := dto.Register{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	token, user, err := userDomain.Register(body)
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
		"user":  user,
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(payload))
}

// HandleLogin godoc
//
//	@Summary		登录
//	@Description	用户登录
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.Login	true	"Login"
//	@Success		200		{object}	dto.Response{data=map[string]interface{}}
//	@Router			/user/login [post]
func (u *UserRouter) HandleLogin(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	body := &dto.Login{}
	if err := ctx.BindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	token, user, err := userDomain.Login(*body)
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

// HandleChangeRole godoc
//
//	@Summary		修改用户角色
//	@Description	根据用户 id 更改权限
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		string	true	"User ID"
//	@Param			role	formData	string	true	"Role"
//	@Success		200		{object}	dto.Response{data=domain.User}
//	@Router			/user/change-role [post]
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

// HandleCredit godoc
//
//	@Summary		获取用户积分
//	@Description	根据用户 id 获取用户积分
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		string	true	"User ID"
//	@Success		200		{object}	dto.Response{data=uint}
//	@Router			/user/statistic/credit [get]
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

// HandleChangePasswd godoc
//
//	@Summary		修改密码
//	@Description	修改密码
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ChangePasswd	true	"Change Password"
//	@Success		200		{object}	dto.Response{data=interface{}}
//	@Router			/user/passwd/change [post]
func (u *UserRouter) HandleChangePasswd(ctx *gin.Context) {
	// 提取请求体到 Register 结构体
	var err error
	body := &dto.ChangePasswd{}
	if err := ctx.BindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}
	userId := ctx.Keys["id"].(uint)

	err = userDomain.ChangePasswd(*body, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleChangeInfo godoc
//
//	@Summary		修改用户信息
//	@Description	修改用户信息
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.ChangeUserInfo	true	"Change User Info"
//	@Success		200		{object}	dto.Response{data=domain.User}
//	@Router			/user/info/change [post]
func (u *UserRouter) HandleChangeInfo(ctx *gin.Context) {
	var err error
	userId := ctx.Keys["id"].(uint)
	body := dto.ChangeUserInfo{}
	if err = ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	user, err := userDomain.ChangeInfo(body, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(user))
}
