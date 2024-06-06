package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"strconv"
)

type AnnotationRouter struct {
}

func NewAnnotationRouter(engine *gin.Engine) {
	router := &AnnotationRouter{}
	annotationGroup := engine.Group("/annotate")
	annotationGroup.GET("/:set_id", router.HandleGetAnnotation)
	annotationGroup.POST("/make", router.HandleMake)
}

var datasetDomain = domain.NewDatasetDomain()
var annotationDomain = domain.NewAnnotationDomain()

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

func (a *AnnotationRouter) HandleMake(ctx *gin.Context) {
	body := dto.NewAnnotation{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	annotation, err := annotationDomain.CreateAnnotation(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(annotation))
}
