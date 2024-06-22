package service

import (
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/domain"
	"sort"
	"strings"
	"time"
)

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
	DatasetId   uint          `json:"dataSetId"`
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
	objects := make([]string, 0)
	if dataset.Tags != "" {
		tags := strings.Split(dataset.Tags, ",")
		for _, tag := range tags {
			if strings.TrimSpace(tag) == "" {
				continue
			} else {
				objects = append(objects, tag)
			}
		}
	}

	return &DatasetResult{
		DatasetId:   dataset.ID,
		DatasetName: dataset.Name,
		TaskInfo:    dataset.Description,
		ObjectCnt:   len(objects),
		Objects:     objects,
		Owner:       isOwner,
		Claim:       isClaim,
		Schedule:    dataset.EndTime.Format("2006-01-02 15:04:05"),
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
func (s *DatasetService) GetAllDatasetList(userID uint) []*DatasetResult {
	var err error
	datasets, err := datasetDomain.GetDatasetList()
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	// 读取用户创建的数据集和加入的数据集
	userCreatedDatasets, err := datasetDomain.ListUserCreatedDatasets(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	userCreatedMap := make(map[int]bool)
	for _, dataset := range userCreatedDatasets {
		userCreatedMap[int(dataset.ID)] = true
	}
	userJoinedDatasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	userJoinedMap := make(map[int]bool)
	for _, dataset := range userJoinedDatasets {
		userJoinedMap[int(dataset.ID)] = true
	}

	// 构建结果列表
	results := make([]*DatasetResult, 0)
	for _, dataset := range datasets {
		isOwner := userCreatedMap[int(dataset.ID)]
		isClaim := userJoinedMap[int(dataset.ID)]
		result := NewDatasetResult(&dataset, isOwner, isClaim)
		results = append(results, result)
	}
	return results
}

// GetUserCreatedDatasetList 获取用户创建的数据集列表
func (s *DatasetService) GetUserCreatedDatasetList(creatorID uint) []*DatasetResult {
	var err error
	datasets, err := datasetDomain.ListUserCreatedDatasets(creatorID)
	if err != nil {
		// 返回空列表
		return make([]*DatasetResult, 0)
	}

	results := s.buildResultList(datasets, true, false)

	return results
}

// GetUserJoinedDatasetList 获取用户加入的数据集列表
func (s *DatasetService) GetUserJoinedDatasetList(userID uint) []*DatasetResult {
	datasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}

	results := s.buildResultList(datasets, false, true)

	return results
}

// GetUserDatasetList 获取用户的数据集列表
func (s *DatasetService) GetUserDatasetList(userID uint) []*DatasetResult {
	createdDatasets, err := datasetDomain.ListUserCreatedDatasets(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	createdResults := s.buildResultList(createdDatasets, true, false)

	joinedDatasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
	if err != nil {
		return make([]*DatasetResult, 0)
	}
	joinedResults := s.buildResultList(joinedDatasets, false, true)

	// 创建一个新的列表，用于存储用户创建的数据集和加入的数据集
	// 同时创建一个map检查是否有重复的数据集
	results := make([]*DatasetResult, 0)
	datasetIdMap := make(map[uint]bool)
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

// QueryDatasetList 查询数据集列表
// TODO: 实现查询数据集列表的逻辑
func (s *DatasetService) QueryDatasetList(userID uint, query *dto.DatasetQuery) []*DatasetResult {
	var _ error
	print(query)
	print(userID)

	// 使用一个Map存储不同方式的查询结果
	queryMap := make(map[string][]domain.Dataset)
	// 用户创建的数据集
	if query.Myself {
		// 查询用户创建的数据集
		createdDatasets, err := datasetDomain.ListUserCreatedDatasets(userID)
		if err != nil {
			return make([]*DatasetResult, 0)
		}
		queryMap["created"] = createdDatasets
	}
	// 用户加入的数据集
	if query.Owner {
		// 查询用户加入的数据集
		joinedDatasets, err := datasetDomain.ListUserJoinedDatasetList(userID)
		if err != nil {
			return make([]*DatasetResult, 0)
		}
		queryMap["joined"] = joinedDatasets
	}
	// 关键字查询
	if query.Keyword != "" {
		// 对关键字进行分割
		keywords := strings.Split(query.Keyword, " ")
		// 查询关键字
		keyDatasets, err := datasetDomain.ListByKeywords(keywords)
		if err != nil {
			return make([]*DatasetResult, 0)
		}
		queryMap["keyword"] = keyDatasets
	}

	// 将所有查询结果合并
	mergedDatasets := make([]domain.Dataset, 0)
	mergedDatasets = append(mergedDatasets, queryMap["created"]...)
	mergedDatasets = append(mergedDatasets, queryMap["joined"]...)
	mergedDatasets = append(mergedDatasets, queryMap["keyword"]...)

	// 包装结果
	var results []*DatasetResult
	// 使用一个map存储数据集ID，避免重复
	datasetIdMap := make(map[uint]bool)
	// 构建结果列表
	for _, dataset := range mergedDatasets {
		// 当数据集ID已经存在时，跳过
		if _, ok := datasetIdMap[dataset.ID]; ok {
			continue
		}
		result := NewDatasetResult(&dataset, userID == dataset.CreatorID, datasetDomain.IsUserClaimDataset(userID, dataset.ID))
		results = append(results, result)
		datasetIdMap[dataset.ID] = true
	}

	// 排序
	if query.Order == dto.OrderTime {
		// 按时间排序
		sort.Slice(results, func(i, j int) bool {
			var iSchedule, jSchedule time.Time
			iSchedule, _ = time.Parse("2006-01-02 15:04:05", results[i].Schedule)
			jSchedule, _ = time.Parse("2006-01-02 15:04:05", results[j].Schedule)
			return iSchedule.After(jSchedule)
		})
	}
	if query.Order == dto.OrderHot {
		// 按热度排序
		sort.Slice(results, func(i, j int) bool {
			return results[i].Total > results[j].Total
		})
	}
	if query.Order == dto.OrderSize {
		// 按大小排序
		sort.Slice(results, func(i, j int) bool {
			return results[i].ObjectCnt > results[j].ObjectCnt
		})
	}

	return results
}

// GetDatasetDetail 获取数据集详情
func (s *DatasetService) GetDatasetDetail(userId uint, id uint) *DatasetResult {
	dataset, err := datasetDomain.GetDatasetByID(id)
	if err != nil {
		return nil
	}
	if dataset == nil {
		return nil
	}
	datas, err := datasetDomain.GetDatasetDataList(id)
	if err != nil {
		return nil
	}

	isOwner := dataset.CreatorID == userId
	isClaim := datasetDomain.IsUserClaimDataset(userId, id)
	result := NewDatasetResult(dataset, isOwner, isClaim)
	result.Datas = make([]DatasetItem, 0)
	for _, data := range datas {
		result.Datas = append(result.Datas, newDatasetItem(&data))
	}

	return result
}

// 将 domain.Dataset 转换为 DatasetResult 的列表
func (s *DatasetService) buildResultList(datasets []domain.Dataset, isOwner bool, isClaim bool) []*DatasetResult {
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
