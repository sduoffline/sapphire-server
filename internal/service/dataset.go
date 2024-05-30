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

func NewDatasetResult(dataset *domain.Dataset, isOwner bool, isClaim bool) *DatasetResult {
	return &DatasetResult{
		DatasetId:   int(dataset.ID),
		DatasetName: dataset.Name,
		TaskInfo:    dataset.Description,
		ObjectCnt:   dataset.Size,
		Owner:       isOwner,
		Claim:       isClaim,
		Schedule:    dataset.Schedule.Format("2006-01-02 15:04:05"),
		Total:       0,
		Finished:    0,
	}
}

func newDatasetItem(data *domain.ImgDataset) DatasetItem {
	return DatasetItem{
		ImgUrl:       data.ImgUrl,
		EmbeddingUrl: data.EmbeddingUrl,
		Id:           int(data.ID),
	}
}

var datasetDomain = domain.NewDatasetDomain()

// GetAllDatasetList 获取数据集列表
func (service *DatasetService) GetAllDatasetList() []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetList()
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := service.buildResultList(datasets, false, false)

	return results
}

// GetUserCreatedDatasetList 获取用户创建的数据集列表
func (service *DatasetService) GetUserCreatedDatasetList(creatorID int) []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetListByUserID(creatorID)
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := service.buildResultList(datasets, true, false)

	return results
}

// GetUserJoinedDatasetList 获取用户加入的数据集列表
func (service *DatasetService) GetUserJoinedDatasetList(userID int) []*DatasetResult {
	datasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}

	results := service.buildResultList(datasets, false, true)

	return results
}

// GetUserDatasetList 获取用户的数据集列表
func (service *DatasetService) GetUserDatasetList(userID int) []*DatasetResult {
	createdDatasets, err := datasetDomain.GetDatasetListByUserID(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	createdResults := service.buildResultList(createdDatasets, true, false)

	joinedDatasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	joinedResults := service.buildResultList(joinedDatasets, false, true)

	// 创建一个新的列表，用于存储用户创建的数据集和加入的数据集
	// 同时创建一个map检查是否有重复的数据集
	results := make([]*DatasetResult, 0)
	datasetIdMap := make(map[int]bool)
	allDatasets := append(createdResults, joinedResults...)
	for _, result := range allDatasets {
		// 当数据集ID已经存在时，跳过
		if _, ok := datasetIdMap[result.DatasetId]; ok {
			continue
		}
		results = append(results, result)
		datasetIdMap[result.DatasetId] = true
	}

	return results
}

// GetDatasetDetail 获取数据集详情
func (service *DatasetService) GetDatasetDetail(id int) *DatasetResult {
	dataset, err := datasetDomain.GetDatasetByID(id)
	if err != nil {
		return nil
	}
	datas, err := datasetDomain.GetDatasetDataList(id)
	if err != nil {
		return nil
	}

	// TODO: 这里的 isOwner 和 isClaim 需要根据用户ID来判断
	result := NewDatasetResult(dataset, false, false)
	result.Datas = make([]DatasetItem, 0)
	for _, data := range datas {
		result.Datas = append(result.Datas, newDatasetItem(&data))
	}

	return result
}

// 将 domain.Dataset 转换为 DatasetResult 的列表
func (service *DatasetService) buildResultList(datasets []domain.Dataset, isOwner bool, isClaim bool) []*DatasetResult {
	results := make([]*DatasetResult, 0)
	for _, dataset := range datasets {
		result := NewDatasetResult(&dataset, isOwner, isClaim)
		results = append(results, result)
	}
	return results
}

// 将 domain.Dataset 转换为 DatasetResult，包含数据集的数据
//func (service *DatasetService) buildResult(dataset domain.Dataset, datas []domain.DatasetData) *DatasetResult {
//	result := NewDatasetResult(&dataset)
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
