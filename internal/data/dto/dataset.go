package dto

type NewDataset struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	EndTime     string   `json:"endTime"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
}

type DatasetQuery struct {
	Myself  bool   `json:"myself"`
	Owner   bool   `json:"owner"`
	Keyword string `json:"keyword"`
	Order   string `json:"order" binding:"omitempty,oneof=time hot size"`
}

// Order Enums
const (
	OrderTime = "time"
	OrderHot  = "hot"
	OrderSize = "size"
)

type AddImage struct {
	DatasetID uint     `json:"datasetId"`
	Images    []string `json:"images"`
}
