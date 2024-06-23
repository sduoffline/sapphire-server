package domain

import (
	"gorm.io/gorm"
	"log/slog"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
)

type Discussion struct {
	gorm.Model
	Title   string `gorm:"column:title"`
	Content string `gorm:"column:content"`
	UserID  uint   `gorm:"column:user_id"`
	ReplyID uint   `gorm:"column:reply_id"`
}

type DiscussionResult struct {
	ID       uint               `json:"id"`
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	UserName string             `json:"userName"`
	Avatar   string             `json:"avatar"`
	IsReply  bool               `json:"isReply"`
	Replies  []DiscussionResult `json:"replies"`
}

func NewDiscussionDomain() *Discussion {
	return &Discussion{}
}

func (d *Discussion) CreateDiscussion(userID uint, dto dto.NewDiscussion) *Discussion {
	var err error
	slog.Info("create discussion, userID: %d, dto: %+v", userID, dto)

	discussion := Discussion{
		Title:   dto.Title,
		Content: dto.Content,
		UserID:  userID,
		ReplyID: dto.ReplyID,
	}
	slog.Info("create discussion: %+v", discussion)

	err = dao.Save[Discussion](discussion)
	if err != nil {
		slog.Error("create discussion failed: %s", err.Error())
		return nil
	}

	return &discussion
}
