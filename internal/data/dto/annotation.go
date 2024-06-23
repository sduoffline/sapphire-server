package dto

type NewAnnotation struct {
	Marks     []AnnotationResult `json:"marks"`
	ImgID     uint               `json:"imgId"`
	DatasetID uint               `json:"datasetId"`
}

type AnnotationResult struct {
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	Width   float64 `json:"w"`
	Height  float64 `json:"h"`
	ID      uint    `json:"id"`
}
