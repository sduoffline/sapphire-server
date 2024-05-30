package service

import "sapphire-server/internal/domain"

type DatasetService struct {
}

func NewDatasetService() *DatasetService {
	return &DatasetService{}
}

type DatasetItem struct {
	ImgUrl       string `json:"imgUrl"`
	EmbeddingUrl string `json:"embeddingUrl"`
	Id           int    `json:"id"`
}

type DatasetResult struct {
	DatasetId   int           `json:"dataSetId"`
	DatasetName string        `json:"dataSetName"`
	TaskInfo    string        `json:"taskInfo"`
	ObjectCnt   int           `json:"objectCnt"`
	Objects     []string      `json:"objects"`
	Owner       bool          `json:"owner"`
	Claim       bool          `json:"claim"`
	Datas       []DatasetItem `json:"datas"`
	Schedule    string        `json:"schedule"`
	Total       int           `json:"total"`
	Finished    int           `json:"finished"`
}

func newDatasetResult(dataset *domain.Dataset) *DatasetResult {
	return &DatasetResult{
		DatasetId:   int(dataset.ID),
		DatasetName: dataset.Name,
		TaskInfo:    dataset.Description,
		ObjectCnt:   dataset.Size,
		Owner:       false,
		Claim:       false,
		Schedule:    "",
		Total:       0,
		Finished:    0,
	}
}

var datasetDomain = domain.NewDataset()

func (service *DatasetService) GetDatasetList() []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetList()
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := make([]*DatasetResult, 0)
	for _, dataset := range datasets {
		result := newDatasetResult(&dataset)
		results = append(results, result)
	}

	return results
}
