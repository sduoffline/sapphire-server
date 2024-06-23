package domain

import (
	"gorm.io/gorm"
	"log/slog"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
)

var messageDomain = NewMessageDomain()

type Message struct {
	gorm.Model
	CreatorID  uint
	ReceiverID uint
	Content    string
	Title      string
	Type       int
}

const (
	// DEFAULT 默认
	DEFAULT = 0
	// MessageTypeTREND 动态
	MessageTypeTREND = 1
	// NOTIFICATION 通知
	NOTIFICATION = 2
	// SYSTEM 系统消息
	SYSTEM = 3
)

func NewMessageDomain() *Message {
	return &Message{}
}

func (m *Message) SendMessage(content, title string, messageType int, receiverID uint) *Message {
	var err error
	message := Message{
		CreatorID:  0,
		ReceiverID: receiverID,
		Content:    content,
		Title:      title,
		Type:       messageType,
	}
	err = dao.Save(&message)
	if err != nil {
		return nil
	}
	return &message
}

// CreateMessage 创建消息
func (m *Message) CreateMessage(creatorID uint, dto dto.NewMessage) *Message {
	var err error

	message := Message{
		CreatorID: creatorID,
		Content:   dto.Content,
		Title:     dto.Title,
		Type:      dto.Type,
	}
	for _, receiverID := range dto.ReceiverID {
		message.ReceiverID = receiverID
	}

	err = dao.Save(&message)
	if err != nil {
		return nil
	}
	return &message
}

// ListMessageByReceiverID 获取接收者的消息
func (m *Message) ListMessageByReceiverID(receiverID uint) []Message {
	messages, err := dao.FindAll[Message]("receiver_id = ?", receiverID)
	if err != nil {
		return nil
	}
	return messages
}

// ReadMessage 标记消息为已读
func (m *Message) ReadMessage(messageID uint) {
	var err error
	message, err := dao.FindOne[Message]("id = ?", messageID)
	if err != nil {
		return
	}

	slog.Debug("message", message)
	err = dao.Delete(message)
	if err != nil {
		return
	}
}
