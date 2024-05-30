package dto

type NewDataset struct {
	DatasetName string `json:"dataSetName"`
	TaskInfo    string `json:"taskInfo"`
	Schedule    string `json:"schedule"`
}
