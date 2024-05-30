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

// GetDatasetList 获取数据集列表
func (service *DatasetService) GetDatasetList() []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetList()
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := service.buildResultList(datasets)

	return results
}

// GetMyDatasetList 获取用户创建的数据集列表
func (service *DatasetService) GetMyDatasetList(creatorID int) []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetListByUserID(creatorID)
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := service.buildResultList(datasets)

	return results
}

// 将 domain.Dataset 转换为 DatasetResult 的列表
func (service *DatasetService) buildResultList(datasets []domain.Dataset) []*DatasetResult {
	results := make([]*DatasetResult, 0)
	for _, dataset := range datasets {
		result := newDatasetResult(&dataset)
		results = append(results, result)
	}
	return results
}

// 将 domain.Dataset 转换为 DatasetResult，包含数据集的数据
//func (service *DatasetService) buildResult(dataset domain.Dataset, datas []domain.DatasetData) *DatasetResult {
//	result := newDatasetResult(&dataset)
//	result.Datas = make([]DatasetItem, 0)
//	for _, data := range datas {
//		result.Datas = append(result.Datas, DatasetItem{
//			ImgUrl:       data.ImgUrl,
//			EmbeddingUrl: data.EmbeddingUrl,
//			Id:           data.ID,
//		})
//	}
//	return result
//}
