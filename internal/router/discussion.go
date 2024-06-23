package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
	"strconv"
)

type DiscussionRouter struct {
}

var discussionDomain = domain.NewDiscussionDomain()

func NewDiscussionRouter(engine *gin.Engine) *DiscussionRouter {
	router := &DiscussionRouter{}

	routerGroup := engine.Group("/discussion").
		Use(middleware.AuthMiddleware()).
		Use(middleware.UserIDMiddleware())
	{
		routerGroup.POST("/create", router.HandleCreate)
		routerGroup.GET("/list/:id", router.HandleList)
	}

	return router
}

func (r *DiscussionRouter) HandleCreate(ctx *gin.Context) {
	var err error
	userID := ctx.Keys["id"].(uint)

	body := dto.NewDiscussion{}
	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	slog.Info("create discussion, userID: %d, body: %+v", userID, body)

	// Create discussion
	res := discussionDomain.CreateDiscussion(userID, body)
	ctx.JSON(200, dto.NewSuccessResponse(res))
}

func (r *DiscussionRouter) HandleList(ctx *gin.Context) {
	var err error
	userID := ctx.Keys["id"].(uint)
	datasetID64, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	datasetID := uint(datasetID64)

	// List discussion
	res := discussionDomain.ListDiscussionsByDatasetID(userID, datasetID)
	ctx.JSON(200, dto.NewSuccessResponse(res))
}
