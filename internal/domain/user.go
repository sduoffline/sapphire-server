package domain

import (
	"errors"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/internal/infra"
	"time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(register dto.Register) error {
	var err error
	// 检查是否有重名用户
	existed, err := dao.FindOne[User]("name = ?", register.Name)
	if err != nil {
		return err
	}
	slog.Info("existed", existed)
	if existed.ID != 0 {
		slog.Info("existed", existed)
		return errors.New("existed user")
	} else {
		slog.Info("Not existed", existed)
	}
	// TODO: 对口令进行加密
	encryptedPasswd := register.Passwd

	// 插入用户
	u.Name = register.Name
	u.Password = encryptedPasswd
	err = dao.Save(u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Login(login dto.Login) error {
	var err error
	// 读取用户
	user, err := dao.FindOne[User]("name = ?", login.Name)
	if err != nil {
		return err
	}
	infra.Redis.Set(infra.Ctx, "name", user.Name, time.Duration(10)*time.Second)
	if user.Password != login.Passwd {
		return errors.New("wrong password")
	}
	return nil
}
