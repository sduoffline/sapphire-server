package domain

import (
	"gorm.io/gorm"
	"log/slog"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
)

type Task struct {
	gorm.Model
	ID           int    `gorm:"column:id"`
	ImgURL       string `gorm:"column:img_url"`
	EmbeddingURL string `gorm:"column:embedding_url"`
	Status       int    `gorm:"column:status"`
	OnnxId       int    `gorm:"column:onnx_id"`
}

const (
	READY   = 0
	RUNNING = 1
	SUCCESS = 2
	FAILED  = -1
)

func NewTask() *Task {
	return &Task{}
}

func (t *Task) CreateTask(job dto.NewJob) {
	t.ImgURL = job.ImgURL
	t.OnnxId = job.OnnxId
	t.Status = READY
	t.EmbeddingURL = ""
	err := dao.Save(t)
	if err != nil {
		return
	}
}

// GetLatestTask 获取最新可做的task
func (t *Task) GetLatestTask() *Task {
	task, err := dao.First[Task]("status = ?", READY)
	if err != nil {
		return nil
	}
	slog.Debug("task", task)

	return task
}

// UpdateTaskStatus 更新task状态
func (t *Task) UpdateTaskStatus(status int) {
	t.Status = status
	err := dao.Save(t)
	if err != nil {
		return
	}
}

// UpdateTaskEmbeddingURL 更新task的embedding url
func (t *Task) UpdateTaskEmbeddingURL(embeddingURL string) {
	t.EmbeddingURL = embeddingURL
	err := dao.Save(t)
	if err != nil {
		return
	}
}

// GetTaskByID 根据id获取task
func (t *Task) GetTaskByID(id int) *Task {
	task, err := dao.First[Task]("id = ?", id)
	if err != nil {
		return nil
	}
	return task
}

// GetAllTasks 获取所有task
func (t *Task) GetAllTasks() []Task {
	tasks, err := dao.FindAll[Task]()
	if err != nil {
		return nil
	}
	slog.Debug("tasks", tasks)

	return tasks
}
