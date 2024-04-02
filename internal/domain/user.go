package domain

import (
	"gorm.io/gorm"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/infra"
)

type User struct {
	gorm.Model
	Name   string
	Passwd string
	dao    *infra.BaseDAO[User]
}

func NewUser() *User {
	user := &User{
		dao: infra.NewBaseDAO[User](),
	}
	return user
}

func (u *User) Register(register dto.Register) error {
	var err error
	// 检查是否有重名用户
	_, err = u.dao.FindOne("name = ?", register.Name)
	if err != nil {
		return err
	}

	// TODO: 对口令进行加密
	encryptedPasswd := register.Passwd

	// 插入用户
	u.Name = register.Name
	u.Passwd = encryptedPasswd
	err = u.dao.Insert()
	if err != nil {
		return err
	}

	return nil
}
