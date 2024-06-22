package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
)

type Score struct {
	gorm.Model
	DatasetID uint `gorm:"column:dataset_id"`
	ImgID     uint `gorm:"column:img_id"`
	UserID    uint `gorm:"column:user_id"`
	Score     int  `gorm:"column:score"`
}

type ScoreResult struct {
	Date  string `json:"date"`
	Score int    `json:"score"`
	Count int    `json:"count"`
}

// NewScoreDomain 创建一个新的 Score 实例
func NewScoreDomain() *Score {
	return &Score{}
}

// CreateScore 创建一个新的 Score 记录
func (s *Score) CreateScore(score *Score) error {
	err := dao.Save(score)
	if err != nil {
		return err
	}
	return nil
}

// ListScoreRecordsInDays 获取用户最近 days 天的评分记录
func (s *Score) ListScoreRecordsInDays(userID uint, days int) ([]ScoreResult, error) {
	var err error
	sql := "SELECT * FROM scores WHERE user_id = ? AND created_at >= now() - make_interval(days := ?)"
	scores, err := dao.Query[Score](sql, userID, days)
	if err != nil {
		return nil, err
	}
	var results []ScoreResult
	for _, score := range scores {
		results = append(results, ScoreResult{
			Date:  score.CreatedAt.Format("2006-01-02"),
			Score: score.Score,
			Count: 1,
		})
	}
	return results, nil
}
