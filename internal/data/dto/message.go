package dto

type NewMessage struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	ReceiverID []uint `json:"receiverId"`
	Type       int    `json:"type" default:"0"`
}
