package domain

import (
	"gorm.io/gorm"
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
	dao.Save(t)
}

// GetLatestTask 获取最新可做的task
func (t *Task) GetLatestTask() *Task {
	task, err := dao.First[Task]("status = ?", READY)
	if err != nil {
		return nil
	}
	return task
}
