package router

import (
	"github.com/gin-gonic/gin"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
)

type ScoreRouter struct{}

func NewScoreRouter(engine *gin.Engine) *ScoreRouter {
	router := &ScoreRouter{}
	scoreGroup := engine.Group("/score")
	authRouter := scoreGroup.Group("/").Use(middleware.AuthMiddleware())
	{
		authRouter.POST("/create", router.HandleCreate)
	}
	return router
}

var scoreDomain = domain.NewScoreDomain()

func (t *ScoreRouter) HandleCreate(ctx *gin.Context) {
	// TODO
	var err error
	sc := &domain.Score{}
	err = scoreDomain.CreateScore(sc)
	if err != nil {
		ctx.JSON(500, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(200, dto.NewSuccessResponse(nil))
}
