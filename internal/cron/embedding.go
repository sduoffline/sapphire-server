package cron

import (
	"github.com/robfig/cron/v3"
	"log/slog"
	"sapphire-server/internal/domain"
	"time"
)

type EmbeddingCron struct {
	// Interval of the embedding cron
	Interval int
	// 启动等待时间
	WaitTime      int
	Cron          *cron.Cron
	DatasetDomain *domain.Dataset
	IsEmbedding   bool
}

func NewEmbeddingCron() *EmbeddingCron {
	cron := &EmbeddingCron{
		Interval:      60,
		WaitTime:      10,
		DatasetDomain: domain.NewDatasetDomain(),
	}

	return cron
}

func (e *EmbeddingCron) Init() {
	slog.Info("Embedding cron is initializing")
	e.Cron = cron.New(cron.WithSeconds())
	e.Cron.AddFunc("@every 10s", func() {
		images, err := e.DatasetDomain.ListAllNotEmbeddedImg(1)
		if err != nil {
			slog.Error("Failed to list all not embedded images")
			return
		}
		if len(images) == 0 {
			return
		}
		if e.IsEmbedding {
			slog.Debug("Embedding is running")
			return
		}
		// Start embedding
		e.IsEmbedding = true
		go func() {
			// 等待几分钟, 模拟embedding过程
			slog.Info("Start embedding image", images[0].ID)
			// 线程休眠
			time.Sleep(1 * time.Minute)
			for _, img := range images {
				slog.Info("Start embedding image", img.ID)
				err := e.DatasetDomain.EmbeddingImg(img.ID)
				if err != nil {
					slog.Error("Failed to embedding image", img.ID)
					continue
				}
				slog.Info("Embedding image", img.ID, "successfully")
			}
			e.IsEmbedding = false
		}()
	})
}

func (e *EmbeddingCron) Start() {
	slog.Info("Embedding cron is starting")
	e.Cron.Start()
}
