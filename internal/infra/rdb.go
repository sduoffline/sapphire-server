package infra

import (
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	slog.Info("DB connected")
	return nil
}

// Insert
// 通过范型实现通用 Insert 函数
func Insert[T any](data T) error {
	res := DB.Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// FindOne 查询一条数据
func FindOne[T any](conditions ...interface{}) (*T, error) {
	var obj T
	// 这里不使用 `Take()` 方法，因为 `Take()` 方法在没有找到数据时会返回 ErrRecordNotFound 错误
	res := DB.Limit(1).Find(&obj, conditions...)
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}
