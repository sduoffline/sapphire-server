package domain

import (
	"github.com/goccy/go-json"
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
)

type Annotation struct {
	gorm.Model
	Status         int    `gorm:"column:status"`
	Content        string `gorm:"column:content"`
	DatasetID      uint   `gorm:"column:dataset_id"`
	ReplicaCount   int    `gorm:"column:replica_count"`
	QualifiedCount int    `gorm:"column:qualified_count"`
	DeliveredCount int    `gorm:"column:delivered_count"`
}

type AnnotationUser struct {
	ID           uint   `gorm:"primaryKey"`
	AnnotationId uint   `gorm:"column:annotation_id"`
	UserId       uint   `gorm:"column:user_id"`
	Status       int    `gorm:"column:status"`
	Result       string `gorm:"column:result"`
	// NOTE: 关于 Result 这边原来设置为 JSON 格式的，嫌麻烦先改成 string 了
}

func NewAnnotationDomain() *Annotation {
	return &Annotation{}
}

func newAnnotationFromDTO(anno dto.NewAnnotation) *Annotation {
	// 将 Marks 从 JSON 转为 string
	marks, _ := json.Marshal(anno.Marks)
	if len(marks) == 0 {
		marks = []byte("[]")
	}
	marksStr := string(marks)

	return &Annotation{
		Content:        marksStr,
		DatasetID:      anno.DatasetID,
		ReplicaCount:   0,
		QualifiedCount: 0,
		DeliveredCount: 0,
	}
}

func (a *Annotation) CreateAnnotation(anno dto.NewAnnotation) (*Annotation, error) {
	var err error
	annotation := newAnnotationFromDTO(anno)
	err = dao.Save(annotation)
	if err != nil {
		return nil, err
	}
	return annotation, nil
}
