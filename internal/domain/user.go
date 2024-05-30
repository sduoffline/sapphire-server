package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"sapphire-server/internal/dao"
	"sapphire-server/internal/data/dto"
	"sapphire-server/pkg/util"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"column:name" json:"name"`
	// 序列化为 JSON 时忽略字段
	Password string `gorm:"column:password" json:"-"`
	Email    string `gorm:"column:email" json:"email"`
	Uid      string `gorm:"column:uid" json:"uid"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
	Role     int    `gorm:"column:role" json:"role"`
	Score    int    `gorm:"column:score" json:"score"`
}

type UserRole struct {
	ID       int    `gorm:"column:id"`
	RoleName string `gorm:"column:role_name"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Register(register dto.Register) (token string, user *User, err error) {
	// 检查是否有重名用户
	existedUser := u.loadUser(map[string]interface{}{"name": register.Name})
	if existedUser != nil {
		return "", nil, errors.New("existed user")
	}

	encryptedPasswd, err := u.hashPassword(register.Passwd)
	if err != nil {
		return "", nil, err
	}

	// 插入用户
	u.Name = register.Name
	u.Password = encryptedPasswd
	if u.isValidEmail(register.Email) {
		u.Email = register.Email
	} else {
		return "", nil, errors.New("invalid email")
	}
	// TODO: 生成 UID
	u.Uid = strconv.FormatInt(time.Now().Unix(), 10)
	u.Avatar = register.Avatar

	// 默认给普通用户权限
	role, err := u.FindRoleIdByRoleName("USER")
	if err != nil {
		return "", nil, err
	} else {
		u.Role = role
	}

	err = dao.Save(u)
	if err != nil {
		return "", nil, err
	}

	// 生成 token
	token = util.GenerateJWT(u.ID)

	return token, nil, nil
}

func (u *User) Login(login dto.Login) (token string, user *User, err error) {
	// 读取用户
	user = u.loadUser(map[string]interface{}{"name": login.Name})
	if user == nil {
		return "", nil, errors.New("user not found")
	}
	// Redis DEMO
	// infra.Redis.Set(infra.Ctx, "name", user.Name, time.Duration(10)*time.Second)
	// 验证口令

	err = user.verifyPassword(user.Password, login.Passwd)
	if err != nil {
		return "", nil, errors.New("密码错误")
	}

	// 生成 token
	token = util.GenerateJWT(user.ID)

	return token, user, nil
}

func (u *User) loadUser(param map[string]interface{}) *User {
	user, err := dao.First[User]("name = ?", param["name"])
	if err != nil {
		return nil
	}
	return user
}

// hashPassword 生成密码哈希
func (u *User) hashPassword(password string) (string, error) {
	// 使用 bcrypt 生成哈希，第二个参数是哈希强度，越高越安全
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// verifyPassword 验证密码是否匹配哈希
func (u *User) verifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

// isValidEmail 验证邮箱格式
func (u *User) isValidEmail(email string) bool {
	// 姑且先用正则表达式来验证邮箱
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	match, err := regexp.MatchString(regex, email)
	if err != nil {
		return false
	}

	return match
}

// FindRoleIdByRoleName 通过权限名查找权限 ID
func (u *User) FindRoleIdByRoleName(roleName string) (int, error) {
	role, err := dao.First[UserRole]("role_name = ?", roleName)
	if err != nil {
		return 0, err
	} else {
		return role.ID, nil
	}
}
