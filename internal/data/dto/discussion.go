package dto

type NewDiscussion struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	DatasetID uint   `json:"datasetId"`
	ReplyID   uint   `json:"replyId"`
}
