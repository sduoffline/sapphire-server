package router

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
	"strconv"
)

type MessageRouter struct{}

func NewMessageRouter(engine *gin.Engine) *MessageRouter {
	router := &MessageRouter{}
	messageGroup := engine.Group("/message")
	authRouter := messageGroup.Group("/").Use(middleware.AuthMiddleware())
	{
		authRouter.POST("/create", router.HandleCreate)
		authRouter.GET("/list", router.HandleList)
		authRouter.POST("/read/:id", router.HandleRead)
	}
	return router
}

var messageDomain = domain.NewMessageDomain()

func (t *MessageRouter) HandleCreate(ctx *gin.Context) {
	var err error
	creatorID := ctx.Keys["id"].(uint)
	body := dto.NewMessage{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(400, dto.NewFailResponse(err.Error()))
		return
	}

	message := messageDomain.CreateMessage(creatorID, body)
	if message == nil {
		ctx.JSON(500, dto.NewFailResponse("create message failed"))
		return
	}

	ctx.JSON(200, dto.NewSuccessResponse(nil))
}

func (t *MessageRouter) HandleList(ctx *gin.Context) {
	var _ error
	receiverID := ctx.Keys["id"].(uint)
	messages := messageDomain.ListMessageByReceiverID(receiverID)
	if messages == nil {
		ctx.JSON(500, dto.NewFailResponse("list message failed"))
		return
	}
	ctx.JSON(200, dto.NewSuccessResponse(messages))
}

func (t *MessageRouter) HandleRead(ctx *gin.Context) {
	var _ error
	messageID, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	messageDomain.ReadMessage(uint(messageID))
	ctx.JSON(200, dto.NewSuccessResponse(nil))
}
