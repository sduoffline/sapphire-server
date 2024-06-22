package domain

import (
	"fmt"
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"time"
)

type Dataset struct {
	gorm.Model
	Name        string    `gorm:"column:name"`
	CreatorID   uint      `gorm:"column:creator_id"`
	Description string    `gorm:"column:description"`
	Cover       string    `gorm:"column:cover"`
	TypeID      int       `gorm:"column:type_id"`
	Format      string    `gorm:"column:format"`
	Tags        string    `gorm:"column:tags"`
	Size        int       `gorm:"column:size"`
	EndTime     time.Time `gorm:"column:end_time"`
	IsPublic    bool      `gorm:"column:is_public"`
}

type DatasetTag struct {
	gorm.Model
	Tag string `gorm:"column:tag"`
}

type DatasetType struct {
	gorm.Model
	TypeName string `gorm:"column:name"`
	Desc     string `gorm:"column:description"`
}

type ImgDataset struct {
	gorm.Model
	ImgUrl       string `gorm:"column:img_url" json:"imgUrl"`
	DatasetId    uint   `gorm:"column:dataset_id" json:"datasetId"`
	EmbeddingUrl string `gorm:"column:embedding_url" json:"embeddingUrl"`
}

type DatasetUser struct {
	gorm.Model
	UserID    uint `gorm:"column:user_id"`
	DatasetID uint `gorm:"column:dataset_id"`
}

func NewDatasetDomain() *Dataset {
	return &Dataset{}
}

// AddUserToDataset 添加用户到数据集
func (d *Dataset) AddUserToDataset(userID uint, datasetID uint) error {
	var err error
	exist, err := dao.FindOne[DatasetUser]("user_id = ? and dataset_id = ?", userID, datasetID)
	if err != nil {
		return err
	}
	if exist != nil {
		return nil
	}

	datasetUser := &DatasetUser{
		UserID:    userID,
		DatasetID: datasetID,
	}
	err = dao.Save(datasetUser)

	if err != nil {
		return err
	}
	return nil
}

// RemoveUserFromDataset 移除用户从数据集
func (d *Dataset) RemoveUserFromDataset(userID uint, datasetID uint) error {
	var err error
	record, err := dao.First[DatasetUser]("user_id = ? and dataset_id = ?", userID, datasetID)
	if err != nil {
		return err
	}
	if record == nil {
		return nil
	}

	err = dao.Delete(record)
	if err != nil {
		return err
	}
	return nil
}

// IsUserClaimDataset 判断用户是否拥有数据集
func (d *Dataset) IsUserClaimDataset(userID uint, datasetID uint) bool {
	record, err := dao.FindOne[DatasetUser]("user_id = ? and dataset_id = ?", userID, datasetID)
	if err != nil {
		return false
	}
	if record == nil {
		return false
	}
	return true
}

// CreateDataset 创建数据集
func (d *Dataset) CreateDataset(creatorId uint, dto dto.NewDataset) (*Dataset, error) {
	// 创建数据集记录
	datasetInfo := &Dataset{
		Name:        dto.Name,
		CreatorID:   creatorId,
		Description: dto.Description,
		Cover:       dto.Cover,
	}

	scheduleTime := time.Now()
	if dto.EndTime != "" {
		scheduleTime, _ = time.Parse("2006-01-02 15:04:05", dto.EndTime)
	}
	datasetInfo.EndTime = scheduleTime

	// 添加标注tag的记录
	tags := dto.Tags
	tagStr := ""
	for _, tag := range tags {
		tagStr += tag + ","
	}
	datasetInfo.Tags = tagStr

	err := dao.Save(datasetInfo)
	if err != nil {
		return nil, err
	}
	return datasetInfo, nil
}

func (d *Dataset) UpdateDataset(creatorID uint, id uint, dto dto.NewDataset) (*Dataset, error) {
	var err error
	dataset, err := d.GetDatasetByID(id)
	if err != nil {
		return nil, err
	}
	if dataset == nil {
		return nil, fmt.Errorf("dataset not found")
	}

	if creatorID != dataset.CreatorID {
		return nil, fmt.Errorf("no permission")
	}
	dataset.Name = dto.Name
	dataset.Description = dto.Description
	dataset.Cover = dto.Cover
	dataset.EndTime, _ = time.Parse("2006-01-02 15:04:05", dto.EndTime)
	tagStr := ""
	for _, tag := range dto.Tags {
		tagStr += tag + ","
	}
	dataset.Tags = tagStr

	err = dao.Save(dataset)
	if err != nil {
		return nil, err
	}

	return dataset, err

}

// DeleteDataset 删除数据集
func (d *Dataset) DeleteDataset() error {
	err := dao.Delete(d)
	if err != nil {
		return err
	}
	return nil
}

// GetDatasetByID 根据 ID 获取数据集
func (d *Dataset) GetDatasetByID(id uint) (*Dataset, error) {
	res, err := dao.First[Dataset]("id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetResultArchive 获取结果归档
// TODO: 未实现
func (d *Dataset) GetResultArchive(id uint) (string, error) {
	_, err := dao.FindAll[ImgDataset]("dataset_id = ?", id)
	if err != nil {
		return "", err
	}
	return "", nil
}

// GetDatasetList 获取数据集列表
func (d *Dataset) GetDatasetList() ([]Dataset, error) {
	res, err := dao.FindAll[Dataset]()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListByKeywords 根据关键字列出数据集
func (d *Dataset) ListByKeywords(keywords []string) ([]Dataset, error) {
	sql := "select * from datasets where name like ?"
	for i := 1; i < len(keywords); i++ {
		sql += " and name like ?"
	}
	res, err := dao.Query[Dataset](sql, keywords[0])
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListAllDataset 列出所有记录
func (d *Dataset) ListAllDataset() ([]Dataset, error) {
	res, err := dao.FindAll[Dataset]()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListUserJoinedDatasetList 列出用户加入的数据集
func (d *Dataset) ListUserJoinedDatasetList(userID uint) ([]Dataset, error) {
	sql := "select * from datasets where id in (select dataset_id from dataset_users where user_id = ?) and creator_id != ?"
	res, err := dao.Query[Dataset](sql, userID, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListUserCreatedDatasets 列出用户创建的数据集
func (d *Dataset) ListUserCreatedDatasets(createdID uint) ([]Dataset, error) {
	res, err := dao.FindAll[Dataset]("creator_id = ?", createdID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetDatasetTypeByID 根据 ID 获取数据集类型
func (d *Dataset) GetDatasetTypeByID(id uint) (*DatasetType, error) {
	res, err := dao.First[DatasetType]("id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RegisterImage 注册图片
func (d *Dataset) RegisterImage(imgUrl string, datasetID uint) {
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
func (d *Dataset) GetDatasetDataList(id uint) ([]ImgDataset, error) {
	res, err := dao.FindAll[ImgDataset]("dataset_id = ?", id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *Dataset) GetImgByDatasetID(id uint, size int) ([]ImgDataset, error) {
	res, err := dao.Query[ImgDataset]("select * from img_datasets where dataset_id = ? limit ?", id, size)
	if err != nil {
		return nil, err
	}
	return res, nil
}
