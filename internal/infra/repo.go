package infra

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sapphire-server/internal/conf"
)

var (
	DB *gorm.DB
)

func InitDB() error {
	dsn := conf.GetDBConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	log.Println("数据库连接成功")
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
	// 这里不使用 `Take()` 方法，因为 `Take()` 方法在没有找到数据时会返回 ErrRecordNotFound 错误
	res := DB.Limit(1).Find(&obj, conditions...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}
