package dto

type NewDiscussion struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ReplyID uint   `json:"replyID"`
}
