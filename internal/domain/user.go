package domain

import (
	"errors"
	"gorm.io/gorm"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/pkg/util"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	Email    string `gorm:"column:email"`
	Uid      string `gorm:"column:uid"`
	Avatar   string `gorm:"column:avatar"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(register dto.Register) (token string, err error) {
	// 检查是否有重名用户
	existedUser := u.loadUser(map[string]interface{}{"name": register.Name})
	if existedUser != nil {
		return "", errors.New("existed user")
	}

	// TODO: 对口令进行加密
	encryptedPasswd := register.Passwd

	// 插入用户
	u.Name = register.Name
	u.Password = encryptedPasswd
	u.Email = register.Email
	// TODO: 生成 UID
	//strconv.Itoa()
	u.Uid = strconv.FormatInt(time.Now().Unix(), 10)
	u.Avatar = register.Avatar
	err = dao.Save(u)
	if err != nil {
		return "", err
	}

	// 生成 token
	token = util.GenerateJWT(u.ID)

	return token, nil
}

func (u *User) Login(login dto.Login) (token string, err error) {
	// 读取用户
	user := u.loadUser(map[string]interface{}{"name": login.Name})
	if user == nil {
		return "", errors.New("user not found")
	}
	// Redis DEMO
	// infra.Redis.Set(infra.Ctx, "name", user.Name, time.Duration(10)*time.Second)
	// 验证口令
	// TODO: 当然，在上了口令加密以后，这里也要进行加密对比
	if user.Password != login.Passwd {
		return "", errors.New("wrong password")
	}

	// 生成 token
	token = util.GenerateJWT(user.ID)

	return token, nil
}

func (u *User) loadUser(param map[string]interface{}) *User {
	user, err := dao.First[User]("name = ?", param["name"])
	if err != nil {
		return nil
	}
	return user
}
