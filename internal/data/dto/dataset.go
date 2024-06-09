package dto

type NewDataset struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	EndTime     string   `json:"endTime"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
}
