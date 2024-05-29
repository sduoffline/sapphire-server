package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
)

type Dataset struct {
	gorm.Model
	Name        string `gorm:"column:name"`
	CreatorID   int    `gorm:"column:creator_id"`
	Description string `gorm:"column:description"`
	TypeID      int    `gorm:"column:type_id"`
	Format      string `gorm:"column:format"`
	Size        int    `gorm:"column:size"`
	IsPublic    bool   `gorm:"column:is_public"`
	IsDeleted   bool   `gorm:"column:is_deleted"`
}

type DatasetType struct {
	Id       int    `gorm:"column:id"`
	TypeName string `gorm:"column:name"`
	Desc     string `gorm:"column:description"`
}

type ImgDataset struct {
	Id        int    `gorm:"column:id"`
	ImgUrl    string `gorm:"column:img_url"`
	DatasetId int    `gorm:"column:dataset_id"`
}

func NewDataset() *Dataset {
	return &Dataset{}
}

// CreateDataset 创建数据集
func (d *Dataset) CreateDataset(dataset dto.NewDataset) {
	d.Name = dataset.Name
	d.CreatorID = dataset.CreatorID
	d.Description = dataset.Description
	d.TypeID = dataset.TypeID
	d.Format = dataset.Format
	d.Size = dataset.Size
	d.IsPublic = dataset.IsPublic
	d.IsDeleted = dataset.IsDeleted
	err := dao.Save(d)
	if err != nil {
		return
	}
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
	res, err := dao.First[Dataset]("id = ?", d.ID)
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
