package infra

import (
	"errors"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
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

func HasID[T any](data T) bool {
	// 通过反射获取字段值
	v := reflect.ValueOf(data).Elem()
	id := v.FieldByName("ID")
	print(id.IsValid())
	// 判断字段是否为空
	return id.IsValid() && !id.IsZero()
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

func InsertMany[T any](data []T) error {
	res := DB.Create(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func InsertManyWithTransaction[T any](data []T) error {
	tx := DB.Begin()
	res := tx.Create(&data)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

// Update 通过范型实现通用 Update 函数
func Update[T any](data T) error {
	res := DB.Save(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func Delete[T any](data T) error {
	res := DB.Delete(&data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// UpdateSingleColumn
// 通过范型实现通用 UpdateSingleColumn 函数
func UpdateSingleColumn[T any](data T, column string, value interface{}) error {
	res := DB.Model(&data).Update(column, value)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// FindOne 查询一条数据
func FindOne[T any](conditions ...interface{}) (*T, error) {
	var obj T
	// 这里不使用 `Take()` 方法，因为 `Take()` 方法在没有找到数据时会返回 ErrRecordNotFound 错误
	res := DB.Take(&obj, conditions...)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}

// First 查询第一条数据
func First[T any](conditions ...interface{}) (*T, error) {
	var obj T
	res := DB.Take(&obj, conditions...)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &obj, nil
}

// FindAll 查询所有数据
func FindAll[T any](conditions ...interface{}) ([]T, error) {
	var objs []T
	res := DB.Find(&objs, conditions...)
	if res.Error != nil {
		return nil, res.Error
	}
	return objs, nil
}

// Query 执行原生 SQL 查询
func Query[T any](sql string, args ...interface{}) ([]T, error) {
	var objs []T
	res := DB.Raw(sql, args...).Scan(&objs)
	if res.Error != nil {
		return nil, res.Error
	}
	return objs, nil
}
