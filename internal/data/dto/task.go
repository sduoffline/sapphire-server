package dto

type NewJob struct {
	ImgURL string `json:"img_url" binding:"required"`
	OnnxId int    `json:"onnx_id" binding:"required"`
}
