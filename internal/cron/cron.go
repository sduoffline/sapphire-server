package cron

type Service struct {
	// 保存各类 cron
	EmbeddingCron *EmbeddingCron
}

func NewCronService() *Service {
	embeddingCron := NewEmbeddingCron()

	return &Service{
		EmbeddingCron: embeddingCron,
	}
}

func (c *Service) Init() {
	c.EmbeddingCron.Init()
}

func (c *Service) Stop() {
	// 销毁各类 cron
}

func (c *Service) Start() {
	c.EmbeddingCron.Start()
}
