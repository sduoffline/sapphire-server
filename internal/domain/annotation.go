package domain

type Annotation struct {
	ID             uint   `gorm:"primaryKey"`
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

func NewAnnotation() *Annotation {
	return &Annotation{}
}
