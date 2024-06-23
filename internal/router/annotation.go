package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
	"strconv"
)

type AnnotationRouter struct {
}

func NewAnnotationRouter(engine *gin.Engine) {
	router := &AnnotationRouter{}
	annotationGroup := engine.Group("/annotate").Use(middleware.AuthMiddleware()).Use(middleware.UserIDMiddleware())
	annotationGroup.GET("/:set_id", router.HandleGetAnnotation)
	annotationGroup.POST("/make", router.HandleMake)
}

var datasetDomain = domain.NewDatasetDomain()
var annotationDomain = domain.NewAnnotationDomain()

// HandleGetAnnotation godoc
//
//	@Summary		获取标注图片信息
//	@Description	根据数据集ID获取标注图片信息
//	@Tags			annotation
//	@Accept			json
//	@Produce		json
//	@Param			set_id	path		int	true	"Dataset ID"
//	@Param			size	query		int	false	"Number of images"
//	@Success		200		{object}	dto.Response{data=[]interface{}}
//	@Router			/annotate/{set_id} [get]
func (a *AnnotationRouter) HandleGetAnnotation(ctx *gin.Context) {
	datasetID, _ := strconv.Atoi(ctx.Param("set_id"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	images, err := datasetDomain.GetImgByDatasetID(uint(datasetID), size)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "get images failed"})
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(images))
}

// HandleMake godoc
//
//	@Summary		创建标注
//	@Description	创建标注
//	@Tags			annotation
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.NewAnnotation	true	"New Annotation"
//	@Success		200		{object}	dto.Response{data=domain.Annotation}
//	@Router			/annotate/make [post]
func (a *AnnotationRouter) HandleMake(ctx *gin.Context) {
	var err error
	body := dto.NewAnnotation{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	userID := ctx.Keys["id"].(uint)

	annotation, err := annotationDomain.CreateAnnotation(userID, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(annotation))
}
