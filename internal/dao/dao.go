package dao

import (
	"sapphire-server/internal/infra"
)

func Save[T any](data T) error {
	var err error
	// 根据有无id字段判断是插入还是更新
	if infra.HasID(data) {
		err = infra.Update(data)
	} else {
		err = infra.Insert(data)
	}
	if err != nil {
		return err
	}
	return nil
}

func Delete[T any](data T) error {
	err := infra.Delete(data)
	if err != nil {
		return err
	}
	return nil
}

func FindOne[T any](conditions ...interface{}) (*T, error) {
	result, err := infra.FindOne[T](conditions...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func First[T any](conditions ...interface{}) (*T, error) {
	result, err := infra.First[T](conditions...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func FindAll[T any](conditions ...interface{}) ([]T, error) {
	all, err := infra.FindAll[T](conditions...)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func FindPage[T any](page int, pageSize int, conditions ...interface{}) ([]T, error) {
	return nil, nil
}

func Query[T any](sql string, args ...interface{}) ([]T, error) {
	result, err := infra.Query[T](sql, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Modify[T any](data T, column string, value string) error {
	err := infra.UpdateSingleColumn(data, column, value)
	if err != nil {
		return err
	}
	return nil
}
