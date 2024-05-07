package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
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

func NewDataset() *Dataset {
	return &Dataset{}
}

// CreateDataset 创建数据集
func (d *Dataset) CreateDataset() {
	// TODO: 创建数据集
}

// DeleteDataset 删除数据集
func (d *Dataset) DeleteDataset() {
	dao.Modify(d, "is_deleted", "1")
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
func (d *Dataset) GetDatasetList() {
	// TODO: 获取数据集列表
}

// GetDatasetListByUserID 根据用户 ID 获取数据集列表
func (d *Dataset) GetDatasetListByUserID() {
	// TODO: 根据用户 ID 获取数据集列表
}
