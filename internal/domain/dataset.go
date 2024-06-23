package domain

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"sapphire-server/internal/conf"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"time"
)

var datasetDomain = NewDatasetDomain()

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
	Status       int    `gorm:"column:status" json:"status"`
	EmbeddingUrl string `gorm:"column:embedding_url" json:"embeddingUrl"`
}

const (
	ImgDatasetStatusDefault           = 0
	ImgDatasetStatusEmbedding         = 1
	ImgDatasetStatusAnnotated         = 2
	ImgDatasetStatusReAnnotation      = 3
	ImgDatasetStatusAnnotationFailed  = 4
	ImgDatasetStatusAnnotationSuccess = 5
)

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

// ListJoinedUserByDatasetID 列出数据集的用户
func (d *Dataset) ListJoinedUserByDatasetID(datasetID uint) ([]DatasetUser, error) {
	res, err := dao.FindAll[DatasetUser]("dataset_id = ?", datasetID)
	if err != nil {
		return nil, err
	}
	return res, nil
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
func (d *Dataset) GetResultArchive(id uint) (string, error) {
	annotations, err := dao.FindAll[Annotation]("dataset_id = ?", id)
	if err != nil {
		return "", err
	}

	var imgIDs []uint
	for _, a := range annotations {
		imgIDs = append(imgIDs, a.ImageID)
	}
	images, err := d.ListImagesByIDs(imgIDs)

	// 将img构建为以id为key的map，方便后续查找
	imgMap := make(map[uint]ImgDataset)
	for _, img := range images {
		imgMap[img.ID] = img
	}

	// 将数据写入一个.txt文件
	// 保存到本地
	// 返回文件路径
	fileName := fmt.Sprintf("result_%d.txt", id)
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("close file failed", err)
		}
	}(f)

	for _, anno := range annotations {
		img := imgMap[anno.ImageID]
		strPattern := "id: %d, content: %s, deliveredCount: %d, qualifiedCount: %d, replicaCount: %d, status: %d, imageUrl: %s\n"
		str := fmt.Sprintf(strPattern, anno.ID, anno.Content, anno.DeliveredCount, anno.QualifiedCount, anno.ReplicaCount, anno.Status, img.ImgUrl)
		_, err := f.WriteString(str)
		if err != nil {
			return "", err
		}
	}

	// 结束写入
	err = f.Sync()
	if err != nil {
		return "", err
	}

	// 将文件上传到OSS
	// 返回文件路径
	info := conf.GetImgConfig()
	re := regexp.MustCompile(`svrUrl: (.*?); directUrl: (.*?); auth string: (.*?);`)
	matches := re.FindStringSubmatch(info)
	var svrUrl, directUrl, auth string
	if len(matches) == 4 {
		svrUrl = matches[1]
		directUrl = matches[2]
		auth = matches[3]
	} else {
		fmt.Println("String format is not valid.")
	}

	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close file:", err)
		}
	}(file)

	// 删除文件
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println("Failed to remove file:", err)
		}
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 生成一个随机文件名
	// 对当前时间戳进行哈希
	timeStr := time.Now().String()
	h := sha1.New()
	h.Write([]byte(timeStr))
	hashedTimeStr := fmt.Sprintf("%x", h.Sum(nil)) // 使用SHA256哈希
	// 使用SHA256哈希
	fileName = fileName + "_" + hashedTimeStr + ".txt"

	req, err := http.NewRequest("PUT", svrUrl+fileName, bytes.NewReader(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", auth)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	// 检查响应状态
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return directUrl + fileName, nil
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

func (d *Dataset) ListImagesByIDs(ids []uint) ([]ImgDataset, error) {
	res, err := dao.FindAll[ImgDataset]("id in ?", ids)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListImagesByStatusAndDatasetID 根据数据集 ID 和状态列出图片
func (d *Dataset) ListImagesByStatusAndDatasetID(datasetID uint, status int) ([]ImgDataset, error) {
	var err error
	var res []ImgDataset
	if status == -1 {
		res, err = dao.FindAll[ImgDataset]("dataset_id = ?", datasetID)
		if err != nil {
			return nil, err
		}

	} else {
		res, err = dao.FindAll[ImgDataset]("dataset_id = ? and status = ?", datasetID, status)
		if err != nil {
			return nil, err
		}
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

// AddImageList 添加图片列表
func (d *Dataset) AddImageList(dataset *Dataset, images []string) error {
	var err error

	var imageList []ImgDataset
	for _, img := range images {
		image := ImgDataset{
			ImgUrl:    img,
			DatasetId: dataset.ID,
		}
		imageList = append(imageList, image)
	}

	println("imageList", imageList)
	err = dao.SaveAll[ImgDataset](imageList)
	if err != nil {
		return err
	}

	return nil
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
