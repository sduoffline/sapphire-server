package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"time"
)

type Dataset struct {
	gorm.Model
	Name        string    `gorm:"column:name"`
	CreatorID   int       `gorm:"column:creator_id"`
	Description string    `gorm:"column:description"`
	TypeID      int       `gorm:"column:type_id"`
	Format      string    `gorm:"column:format"`
	Size        int       `gorm:"column:size"`
	Schedule    time.Time `gorm:"column:schedule"`
	IsPublic    bool      `gorm:"column:is_public"`
}

type DatasetType struct {
	gorm.Model
	TypeName string `gorm:"column:name"`
	Desc     string `gorm:"column:description"`
}

type ImgDataset struct {
	gorm.Model
	ImgUrl       string `gorm:"column:img_url" json:"imgUrl"`
	DatasetId    int    `gorm:"column:dataset_id" json:"datasetId"`
	EmbeddingUrl string `gorm:"column:embedding_url" json:"embeddingUrl"`
}

func NewDatasetDomain() *Dataset {
	return &Dataset{}
}

// CreateDataset 创建数据集
func (d *Dataset) CreateDataset(creatorId int, dto dto.NewDataset) (*Dataset, error) {
	datasetInfo := &Dataset{
		Name:        dto.DatasetName,
		CreatorID:   creatorId,
		Description: dto.TaskInfo,
	}

	scheduleTime := time.Now()
	if dto.Schedule != "" {
		scheduleTime, _ = time.Parse("2006-01-02 15:04:05", dto.Schedule)

	}
	datasetInfo.Schedule = scheduleTime

	err := dao.Save(datasetInfo)
	if err != nil {
		return nil, err
	}
	return datasetInfo, nil
}

// DeleteDataset 删除数据集
func (d *Dataset) DeleteDataset() {
	err := dao.Modify(d, "is_deleted", "1")
	if err != nil {
		return
	}
}

// GetDatasetByID 根据 ID 获取数据集
func (d *Dataset) GetDatasetByID(id int) (*Dataset, error) {
	res, err := dao.First[Dataset]("id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetDatasetList 获取数据集列表
func (d *Dataset) GetDatasetList() ([]Dataset, error) {
	res, err := dao.FindAll[Dataset]()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetDatasetListByUserID 根据用户 ID 获取数据集列表
func (d *Dataset) GetDatasetListByUserID(createdID int) ([]Dataset, error) {
	res, err := dao.FindAll[Dataset]("creator_id = ?", createdID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetDatasetTypeByID 根据 ID 获取数据集类型
func (d *Dataset) GetDatasetTypeByID(id int) (*DatasetType, error) {
	res, err := dao.First[DatasetType]("id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RegisterImage 注册图片
func (d *Dataset) RegisterImage(imgUrl string, datasetID int) {
	img := &ImgDataset{
		ImgUrl:    imgUrl,
		DatasetId: datasetID,
	}
	err := dao.Save(img)
	if err != nil {
		return
	}
}

// GetDatasetDataList 获取数据集数据列表
func (d *Dataset) GetDatasetDataList(id int) ([]ImgDataset, error) {
	res, err := dao.FindAll[ImgDataset]("dataset_id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Dataset) GetImgByDatasetID(id int, size int) ([]ImgDataset, error) {
	res, err := dao.Query[ImgDataset]("select * from img_datasets where dataset_id = ? limit ?", id, size)
	if err != nil {
		return nil, err
	}
	return res, nil
}
