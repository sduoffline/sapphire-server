package domain

import (
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"log/slog"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/datatypes"
	"sapphire-server/internal/data/dto"
)

var annotationDomain = NewAnnotationDomain()

type Annotation struct {
	gorm.Model
	Status         int            `gorm:"column:status"`
	Content        datatypes.JSON `gorm:"column:content"`
	DatasetID      uint           `gorm:"column:dataset_id"`
	ImageID        uint           `gorm:"column:image_id"`
	UserID         uint           `gorm:"column:user_id"`
	IsQualified    bool           `gorm:"column:is_qualified"`
	ReplicaCount   int            `gorm:"column:replica_count"`
	QualifiedCount int            `gorm:"column:qualified_count"`
	DeliveredCount int            `gorm:"column:delivered_count"`
}

type AnnotationUser struct {
	ID           uint   `gorm:"primaryKey"`
	AnnotationID uint   `gorm:"column:annotation_id"`
	UserId       uint   `gorm:"column:user_id"`
	Status       int    `gorm:"column:status"`
	Result       string `gorm:"column:result"`
	// NOTE: 关于 Result 这边原来设置为 JSON 格式的，嫌麻烦先改成 string 了
}

func NewAnnotationDomain() *Annotation {
	return &Annotation{}
}

func newAnnotationFromDTO(userID uint, anno dto.NewAnnotation) *Annotation {
	// 将 Marks 从 JSON 转为 string
	marks, _ := json.Marshal(anno.Marks)
	if len(marks) == 0 {
		marks = []byte("[]")
	}
	marksStr := string(marks)
	slog.Info("marksStr", marksStr)

	return &Annotation{
		Content:        datatypes.JSON(marksStr),
		DatasetID:      anno.DatasetID,
		UserID:         userID,
		IsQualified:    true,
		ReplicaCount:   0,
		QualifiedCount: 0,
		DeliveredCount: 0,
	}
}

func (a *Annotation) CreateAnnotation(userID uint, anno dto.NewAnnotation) (*Annotation, error) {
	var err error

	// 创建并保存标注
	annotation := newAnnotationFromDTO(userID, anno)
	err = dao.Save(annotation)
	if err != nil {
		return nil, err
	}
	slog.Info("Create Annotation Success", annotation)

	// 读出已经保存的标注，检查是否符合要求
	annotations, err := a.ListAnnotationsByImageID(annotation.ImageID)
	if err != nil {
		return nil, err
	}
	slog.Debug("annotations", annotations)

	return annotation, nil
}

// ListAnnotationsByUserID 根据用户 ID 获取该用户的所有标注
func (a *Annotation) ListAnnotationsByUserID(userID uint) ([]Annotation, error) {
	var err error
	annotations, err := dao.FindAll[Annotation]("user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

func (a *Annotation) ListAnnotationsByImageID(imageID uint) ([]Annotation, error) {
	var err error
	annotations, err := dao.FindAll[Annotation]("image_id = ?", imageID)
	if err != nil {
		return nil, err
	}
	return annotations, nil
}
