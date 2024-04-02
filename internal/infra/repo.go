package infra

import "gorm.io/gorm"

var (
	DB *gorm.DB
)

func InitDB() error {
	// TODO: implement this
	return nil
}

type BaseDAO[T any] struct {
	data T
}

// NewBaseDAO 创建一个BaseDAO
func NewBaseDAO[T any]() *BaseDAO[T] {
	return &BaseDAO[T]{}
}

// Insert 插入数据
func (dao *BaseDAO[T]) Insert() error {
	res := DB.Create(&dao.data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// FindOne 查询一条数据
func (dao *BaseDAO[T]) FindOne(conditions ...interface{}) (*T, error) {
	var obj T
	res := DB.Take(&obj, conditions...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}
