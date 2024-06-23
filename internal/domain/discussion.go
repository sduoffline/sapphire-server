package domain

import (
	"gorm.io/gorm"
	"log/slog"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/pkg/util"
)

type Discussion struct {
	gorm.Model
	Title     string `gorm:"column:title"`
	Content   string `gorm:"column:content"`
	DatasetID uint   `gorm:"column:dataset_id"`
	UserID    uint   `gorm:"column:user_id"`
	ReplyID   uint   `gorm:"column:reply_id"`
}

type DiscussionResult struct {
	ID       uint               `json:"id"`
	Title    string             `json:"title"`
	Content  string             `json:"content"`
	UserName string             `json:"userName"`
	Avatar   string             `json:"avatar"`
	IsReply  bool               `json:"isReply"`
	Replies  []DiscussionResult `json:"replies"`
	Time     string             `json:"time"`
}

func NewDiscussionDomain() *Discussion {
	return &Discussion{}
}

var userDomain = NewUserDomain()

func buildResult(discussion Discussion, user User) *DiscussionResult {
	timeStr := util.FormatTimeStr(discussion.CreatedAt)
	result := DiscussionResult{
		ID:       discussion.ID,
		Title:    discussion.Title,
		Content:  discussion.Content,
		UserName: user.Name,
		Avatar:   user.Avatar,
		Time:     timeStr,
	}

	return &result
}

func (d *Discussion) CreateDiscussion(userID uint, dto dto.NewDiscussion) *DiscussionResult {
	var err error
	slog.Info("create discussion, userID: %d, dto: %+v", userID, dto)

	//timeStr ：= util.FormatTimeStr(time.Now())
	discussion := Discussion{
		Title:     dto.Title,
		Content:   dto.Content,
		DatasetID: dto.DatasetID,
		UserID:    userID,
		ReplyID:   dto.ReplyID,
	}
	slog.Info("create discussion: %+v", discussion)

	err = dao.Save(&discussion)
	if err != nil {
		slog.Error("create discussion failed: %s", err.Error())
		return nil
	}

	user, err := userDomain.GetUserInfo(discussion.UserID)
	if err != nil {
		slog.Error("get user info failed: %s", err.Error())
		return nil
	}

	res := buildResult(discussion, *user)

	return res
}

func (d *Discussion) GetDiscussion(id uint) *DiscussionResult {
	var err error
	slog.Info("get discussion, id: %d", id)

	discussion, err := dao.FindOne[Discussion]("id = ?", id)
	if err != nil {
		slog.Error("get discussion failed: %s", err.Error())
		return nil
	}

	user, err := userDomain.GetUserInfo(discussion.UserID)
	if err != nil {
		slog.Error("get user info failed: %s", err.Error())
		return nil
	}

	res := buildResult(*discussion, *user)
	return res
}

func (d *Discussion) ListDiscussionsByDatasetID(userID uint, datasetID uint) []DiscussionResult {
	var err error
	slog.Info("list discussions by datasetID: %d", datasetID)

	discussions, err := dao.FindAll[Discussion]("dataset_id = ?", datasetID)
	if err != nil {
		slog.Error("list discussions failed: %s", err.Error())
		return nil
	}

	var results []DiscussionResult
	for _, discussion := range discussions {
		user, err := userDomain.GetUserInfo(discussion.UserID)
		if err != nil {
			slog.Error("get user info failed: %s", err.Error())
			return nil
		}

		result := buildResult(discussion, *user)
		results = append(results, *result)
	}

	return results
}
