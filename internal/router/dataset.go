package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sapphire-server/internal/middleware"
	"sapphire-server/internal/service"
	"sapphire-server/pkg/util"
	"strconv"
	"strings"
)

type DatasetRouter struct {
}

func NewDatasetRouter(engine *gin.Engine) *DatasetRouter {
	router := &DatasetRouter{}
	datasetGroup := engine.Group("/dataset")
	datasetGroup.GET("/list", router.HandleList)

	authRouter := datasetGroup.Group("/").Use(middleware.AuthMiddleware())
	{
		authRouter.GET("/created/list", router.HandleCreatedList)
		authRouter.GET("/joined/list", router.HandleJoinedList)
		authRouter.GET("/user/list", router.HandleUserList)

		authRouter.GET("/joined/users/:id", router.ListDatasetJoinedUsers)

		authRouter.POST("/query", router.HandleQuery)
		authRouter.POST("/create", router.HandleCreate)
		authRouter.PUT("/update/:id", router.HandleUpdate)
		authRouter.POST("/upload/:id", router.HandleUploadImg)
		authRouter.POST("/download/:id", router.HandleDownloadDataset)
		authRouter.DELETE("/:id", router.HandleDelete)
		//authRouter.POST("/register", router.HandleRegister)
		authRouter.GET("/:id", router.HandleGetByID)
		authRouter.POST("/join/:id", router.HandleJoin)
		authRouter.POST("/quit/:id", router.HandleQuit)
	}
	return router
}

var datasetService = service.NewDatasetService()

// HandleList 获取公开且未删除的数据集
func (t *DatasetRouter) HandleList(ctx *gin.Context) {
	userID := ctx.Keys["id"].(uint)

	datasets := datasetService.GetAllDatasetList(userID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleCreatedList 获取用户创建的数据集
func (t *DatasetRouter) HandleCreatedList(ctx *gin.Context) {
	userID := ctx.Keys["id"].(uint)
	datasets := datasetService.GetUserCreatedDatasetList(userID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleJoinedList 获取用户加入的数据集
func (t *DatasetRouter) HandleJoinedList(ctx *gin.Context) {
	userID := ctx.Keys["id"].(uint)
	datasets := datasetService.GetUserJoinedDatasetList(userID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleUserList 获取用户创建的数据集
func (t *DatasetRouter) HandleUserList(ctx *gin.Context) {
	userID := ctx.Keys["id"].(uint)
	datasets := datasetService.GetUserDatasetList(userID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}

// HandleJoin 加入数据集
func (t *DatasetRouter) HandleJoin(ctx *gin.Context) {
	var err error

	setID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}
	creatorID := ctx.Keys["id"].(uint)

	err = domain.NewDatasetDomain().AddUserToDataset(creatorID, uint(setID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleQuit 退出数据集
func (t *DatasetRouter) HandleQuit(ctx *gin.Context) {
	var err error

	setID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}
	creatorID := ctx.Keys["id"].(uint)

	err = domain.NewDatasetDomain().RemoveUserFromDataset(creatorID, uint(setID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleCreate 创建数据集
func (t *DatasetRouter) HandleCreate(ctx *gin.Context) {
	var err error
	creatorID := ctx.Keys["id"].(uint)
	// 做个creatorID的校验
	//res, err := dao.First[domain.User]("id = ?", creatorID)
	//if res == nil || err != nil {
	//	ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("creator_id is invalid"))
	//	return
	//}

	// 提取请求体到 NewDataset 结构体
	body := dto.NewDataset{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	// 创建数据集
	dataset, err := domain.NewDatasetDomain().CreateDataset(creatorID, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 返回创建的数据集
	res := datasetService.GetDatasetDetail(creatorID, dataset.ID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(res))
}

func (t *DatasetRouter) HandleUpdate(ctx *gin.Context) {
	var err error
	creatorID := ctx.Keys["id"].(uint)
	datasetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}

	// 提取请求体到 NewDataset 结构体
	body := dto.NewDataset{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	// 创建数据集
	dataset, err := datasetDomain.UpdateDataset(creatorID, uint(datasetID), body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 返回创建的数据集
	res := datasetService.GetDatasetDetail(creatorID, dataset.ID)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(res))
}

// HandleDelete 删除数据集
func (t *DatasetRouter) HandleDelete(ctx *gin.Context) {
	var err error
	datasetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}
	dataset, err := datasetDomain.GetDatasetByID(uint(datasetID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	// 删除数据集
	if dataset == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("dataset not found"))
		return
	} else {
		err = dataset.DeleteDataset()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
			return
		}
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// ListDatasetJoinedUsers 列出加入数据集的用户
func (t *DatasetRouter) ListDatasetJoinedUsers(ctx *gin.Context) {
	var err error
	datasetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}

	// 获取用户
	datasetUsers, err := datasetDomain.ListJoinedUserByDatasetID(uint(datasetID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	//// 遍历 datasetUsers ，去掉自己
	//for i, relation := range datasetUsers {
	//	if relation.UserID == ctx.Keys["id"] {
	//		datasetUsers = append(datasetUsers[:i], datasetUsers[i+1:]...)
	//	}
	//}

	ids := make([]uint, 0)
	for _, relation := range datasetUsers {
		ids = append(ids, relation.UserID)
	}

	var res []domain.User
	if len(ids) > 0 {
		res, err = userDomain.ListUsersByIds(ids)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(res))
}

// HandleUploadImg 上传图片
func (t *DatasetRouter) HandleUploadImg(ctx *gin.Context) {
	var err error
	// 读取dataset id
	_, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}

	// 读取表单的文件
	file, err := ctx.FormFile("file")
	// 检查文件是否存在
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("文件不存在"))
		return
	}

	// 将文件保存到本地
	savePath := "./files/" + file.Filename
	saveDir := strings.Replace(savePath, ".zip", "", 1)
	// 检查是否存在该文件
	if _, err := os.Stat(savePath); err == nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("文件已存在"))
		return
	}
	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 结束后删除文件
	defer func() {
		err := os.Remove(savePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
			return
		}
	}()

	// 解压缩文件
	// 先将文件读取为[]byte
	bytes, err := os.ReadFile(savePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}
	// 创建对应的目录
	err = os.MkdirAll(saveDir, os.ModePerm)
	// 检查目录是否创建成果
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 解压文件
	err = util.Unzip(bytes, "./files/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 取出该目录下的所有文件
	files, err := os.ReadDir(saveDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 遍历文件，将文件上传到图床
	for _, f := range files {
		// 读取文件
		filePath := saveDir + "/" + f.Name()
		_, err := os.ReadFile(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
			return
		}
	}

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
	return
}

// Deprecated: HandleRegister 注册图片
func (t *DatasetRouter) HandleRegister(ctx *gin.Context) {
	imgUrl := ctx.PostForm("img_url")
	datasetId, _ := strconv.Atoi(ctx.PostForm("dataset_id"))
	// Check if the dataset exists
	dataset, err := dao.First[domain.Dataset]("id = ?", datasetId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse("dataset not found"))
		return
	}
	dataset.RegisterImage(imgUrl, uint(datasetId))
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nil))
}

// HandleGetByID 根据 ID 获取数据集
func (t *DatasetRouter) HandleGetByID(ctx *gin.Context) {
	var err error
	datasetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}

	userID, _ := ctx.Keys["id"].(uint)
	dataset := datasetService.GetDatasetDetail(userID, uint(datasetID))

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(dataset))
}

// HandleDownloadDataset 下载数据集
func (t *DatasetRouter) HandleDownloadDataset(ctx *gin.Context) {
	var err error
	datasetID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid dataset id"))
		return
	}

	// 先读出数据集
	//userID, _ := ctx.Keys["id"].(uint)
	//dataset := datasetService.GetDatasetDetail(userID, uint(datasetID))

	// 默认下载所有数据
	url, err := datasetDomain.GetResultArchive(uint(datasetID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse(err.Error()))
		return
	}

	// 使用map包装结果
	res := make(map[string]string)
	res["url"] = url

	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(res))
}

// HandleQuery 查询数据集
func (t *DatasetRouter) HandleQuery(ctx *gin.Context) {
	var err error
	userID := ctx.Keys["id"].(uint)
	query := &dto.DatasetQuery{}
	if err = ctx.BindJSON(query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse(err.Error()))
		return
	}

	// 检查排序字段是否合法
	if query.Order != "" && query.Order != dto.OrderTime && query.Order != dto.OrderHot && query.Order != dto.OrderSize {
		ctx.JSON(http.StatusBadRequest, dto.NewFailResponse("invalid order"))
		return
	}

	datasets := datasetService.QueryDatasetList(userID, query)
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(datasets))
}
