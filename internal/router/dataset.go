package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/service"
	"strconv"
)

type DatasetRouter struct {
}

func NewDatasetRouter(engine *gin.Engine) *DatasetRouter {
	router := &DatasetRouter{}
	datasetGroup := engine.Group("/dataset")
	datasetGroup.GET("/list", router.HandleList)
	datasetGroup.GET("/my-list", router.HandleMyList)
	datasetGroup.POST("/create", router.HandleCreate)
	datasetGroup.POST("/delete", router.HandleDelete)
	datasetGroup.POST("/register", router.HandleRegister)
	datasetGroup.GET("/:id", router.HandleGetByID)
	return router
}

var datasetService = service.NewDatasetService()

// HandleList 获取公开且未删除的数据集
func (t *DatasetRouter) HandleList(ctx *gin.Context) {
	datasets := datasetService.GetDatasetList()
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleMyList 获取用户创建的数据集
func (t *DatasetRouter) HandleMyList(ctx *gin.Context) {
	// TODO: 从token中获取用户ID
	creatorID, _ := strconv.Atoi(ctx.Query("creator_id"))
	datasets := datasetService.GetMyDatasetList(creatorID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleCreate 创建数据集
func (t *DatasetRouter) HandleCreate(ctx *gin.Context) {
	name := ctx.PostForm("name")
	creatorID, _ := strconv.Atoi(ctx.PostForm("creator_id"))
	typeID, _ := strconv.Atoi(ctx.PostForm("type_id"))

	// 做个creatorID的校验
	res, err := dao.First[domain.User]("id = ?", creatorID)
	if res == nil || err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("creator_id is invalid"))
		return
	}

	var description string
	if desc := ctx.PostForm("description"); desc != "" {
		description = desc
	}

	var format string
	if fmt := ctx.PostForm("format"); fmt != "" {
		format = fmt
	}

	var size int
	if szStr := ctx.PostForm("size"); szStr != "" {
		size, _ = strconv.Atoi(szStr)
	}

	var isPublic bool
	if isPubStr := ctx.PostForm("is_public"); isPubStr != "" {
		isPublic, _ = strconv.ParseBool(isPubStr)
	}

	var isDeleted bool
	if isDelStr := ctx.PostForm("is_deleted"); isDelStr != "" {
		isDeleted, _ = strconv.ParseBool(isDelStr)
	}

	datasetInfo := dto.NewDataset{
		Name:        name,
		CreatorID:   creatorID,
		TypeID:      typeID,
		Description: description,
		Format:      format,
		Size:        size,
		IsPublic:    isPublic,
		IsDeleted:   isDeleted,
	}

	NewDataset := domain.NewDataset()
	NewDataset.CreateDataset(datasetInfo)
}

// HandleDelete 删除数据集
func (t *DatasetRouter) HandleDelete(ctx *gin.Context) {
	datasetID, _ := strconv.Atoi(ctx.Query("dataset_id"))
	dataset := domain.NewDataset()
	dataset.ID = uint(datasetID)
	dataset.DeleteDataset()
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleRegister 注册图片
func (t *DatasetRouter) HandleRegister(ctx *gin.Context) {
	imgUrl := ctx.PostForm("img_url")
	datasetId, _ := strconv.Atoi(ctx.PostForm("dataset_id"))
	// Check if the dataset exists
	dataset, err := dao.First[domain.Dataset]("id = ?", datasetId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse("dataset not found"))
		return
	}
	dataset.RegisterImage(imgUrl, datasetId)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleGetByID 根据 ID 获取数据集
func (t *DatasetRouter) HandleGetByID(ctx *gin.Context) {
	datasetID, _ := strconv.Atoi(ctx.Param("id"))
	dataset := domain.NewDataset()
	dataset.ID = uint(datasetID)
	res, err := dao.First[domain.Dataset]("id = ?", datasetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(res))
}
