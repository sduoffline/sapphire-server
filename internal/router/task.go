package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"strconv"
)

type TaskRouter struct {
}

func NewTaskRouter(engine *gin.Engine) *TaskRouter {
	router := &TaskRouter{}
	taskGroup := engine.Group("/task")
	taskGroup.GET("/list", router.HandleList)
	taskGroup.GET("/next", router.HandleNext)
	taskGroup.POST("/create", router.HandleCreate)
	taskGroup.POST("/update", router.HandleUpdate)
	return router
}

func (t *TaskRouter) HandleList(ctx *gin.Context) {
	tasks, _ := dao.FindAll[domain.Task]()
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(tasks))
}

func (t *TaskRouter) HandleCreate(ctx *gin.Context) {
	imgUrl := ctx.PostForm("img_url")
	onnxId := ctx.PostForm("onnx_id")
	// Convert onnxId to int
	onnxIdInt, _ := strconv.Atoi(onnxId)
	// Check if the onnxId is valid
	_, err := dao.First[domain.Sam]("id = ?", onnxIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse("onnxId is invalid"))
		return
	}
	taskInfo := dto.NewJob{ImgURL: imgUrl, OnnxId: onnxIdInt}
	NewTask := domain.NewTask()
	NewTask.CreateTask(taskInfo)
	ctx.JSON(http.StatusOK, NewTask)
}

func (t *TaskRouter) HandleNext(ctx *gin.Context) {
	task := domain.NewTask()
	nextTask := task.GetLatestTask()
	if nextTask == nil {
		ctx.JSON(http.StatusOK, dto.NewFailResponse("No task available"))
		return
	}
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(nextTask))
}

func (t *TaskRouter) HandleUpdate(ctx *gin.Context) {
	taskId := ctx.PostForm("task_id")
	status := ctx.PostForm("status")
	// Convert taskId to int
	taskIdInt, _ := strconv.Atoi(taskId)
	// Convert status to int
	statusInt, _ := strconv.Atoi(status)
	// Check if the task exists
	task, _ := dao.First[domain.Task]("id = ?", taskIdInt)
	if task == nil {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse("task not found"))
		return
	}
	// Check if the status is valid
	if statusInt != domain.READY && statusInt != domain.RUNNING && statusInt != domain.SUCCESS && statusInt != domain.FAILED {
		ctx.JSON(http.StatusInternalServerError, dto.NewFailResponse("status is invalid"))
		return
	}
	// Update the task status
	dao.Modify(task, "status", strconv.Itoa(statusInt))
	ctx.JSON(http.StatusOK, dto.NewSuccessResponse(task))
}
