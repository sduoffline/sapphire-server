package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/domain"
	"strconv"
)

type AnnotationRouter struct {
}

func NewAnnotationRouter(engine *gin.Engine) {
	router := &AnnotationRouter{}
	annotationGroup := engine.Group("/annotation")
	annotationGroup.POST("/make", router.HandleMake)
}

func (a *AnnotationRouter) HandleMake(ctx *gin.Context) {
	datasetId, _ := strconv.Atoi(ctx.PostForm("dataset_id"))
	content := ctx.PostForm("content")
	markerId, _ := strconv.Atoi(ctx.PostForm("marker_id"))
	var ReplicaCount int
	if rc := ctx.PostForm("replica_count"); rc != "" {
		ReplicaCount, _ = strconv.Atoi(rc)
	}
	var QualifiedCount int
	if qc := ctx.PostForm("qualified_count"); qc != "" {
		QualifiedCount, _ = strconv.Atoi(qc)
	}
	var DeliveredCount int
	if dc := ctx.PostForm("delivered_count"); dc != "" {
		DeliveredCount, _ = strconv.Atoi(dc)
	}
	// Create a new annotation
	annotation := domain.NewAnnotation()
	annotation.DatasetID = uint(datasetId)
	annotation.Content = content
	annotation.ReplicaCount = ReplicaCount
	annotation.QualifiedCount = QualifiedCount
	annotation.DeliveredCount = DeliveredCount
	// Save the annotation
	err := dao.Save(annotation)
	if err != nil {
		return
	}

	annotationUser := domain.AnnotationUser{
		ID:           0,
		AnnotationId: annotation.ID,
		UserId:       uint(markerId),
		Status:       0,
		Result:       "",
	}
	err = dao.Save(annotationUser)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, annotation)
}
