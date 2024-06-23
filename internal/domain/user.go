package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"math"
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
	Password    string `gorm:"column:password" json:"-"`
	Email       string `gorm:"column:email" json:"email"`
	Uid         string `gorm:"column:uid" json:"uid"`
	Description string `gorm:"column:description" json:"description"`
	Avatar      string `gorm:"column:avatar" json:"avatar"`
	Role        int    `gorm:"column:role" json:"role"`
	Score       int    `gorm:"column:score" json:"score"`
}

type UserRole struct {
	ID       int    `gorm:"column:id"`
	RoleName string `gorm:"column:role_name"`
}

type UserResult struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Uid            string `json:"uid"`
	Description    string `json:"description"`
	Avatar         string `json:"avatar"`
	Role           int    `json:"role"`
	Score          int    `json:"score"`
	JoinedCount    int    `json:"joinedCount"`
	CreatedCount   int    `json:"createdCount"`
	AnnotatedCount int    `json:"annotatedCount"`
}

func NewUserDomain() *User {
	return &User{}
}

func buildUserResult(user *User, joinedCount, createdCount, annotatedCount int) *UserResult {
	// 使用一个增长越来越缓慢的函数来计算积分
	createScore := math.Pow(float64(createdCount), 1.1)
	joinScore := math.Pow(float64(joinedCount), 1.2)
	annotateScore := math.Pow(float64(annotatedCount), 1.3)
	slog.Info("createScore", createScore)
	slog.Info("joinScore", joinScore)
	slog.Info("annotateScore", annotateScore)
	totalScore := int(10 * (createScore + joinScore + annotateScore) / 3)

	return &UserResult{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		Uid:            user.Uid,
		Description:    user.Description,
		Avatar:         user.Avatar,
		Role:           user.Role,
		Score:          totalScore,
		JoinedCount:    joinedCount,
		CreatedCount:   createdCount,
		AnnotatedCount: annotatedCount,
	}
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

func (u *User) ChangePasswd(passwd dto.ChangePasswd, userId uint) error {
	var err error
	user, err := dao.First[User](userId)
	if err != nil {
		return err
	}

	// 验证口令
	err = user.verifyPassword(user.Password, passwd.Old)
	if err != nil {
		return errors.New("password error")
	}

	// 更新密码
	encryptedPasswd, err := u.hashPassword(passwd.New)
	if err != nil {
		return err
	}
	user.Password = encryptedPasswd

	err = dao.Modify(user, "password", encryptedPasswd)
	if err != nil {
		return err
	}

	return nil
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

func (u *User) GetUserInfo(userId uint) (*User, error) {
	user, err := dao.First[User](userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) GetUserDetail(userId uint) (*UserResult, error) {
	var err error
	user, err := dao.FindOne[User](userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	createdDatasets, err := datasetDomain.ListUserCreatedDatasets(userId)
	if err != nil {
		slog.Error("ListUserCreatedDatasets", err)
		return nil, err
	}
	slog.Info("createdDatasets", createdDatasets)

	joinedDatasets, err := datasetDomain.ListUserJoinedDatasetList(userId)
	if err != nil {
		slog.Error("ListUserJoinedDatasetList", err)
		return nil, err
	}
	slog.Info("joinedDatasets", joinedDatasets)

	annotations, err := annotationDomain.ListAnnotationsByUserID(userId)
	if err != nil {
		slog.Error("ListAnnotationsByUserID", err)
		return nil, err
	}
	slog.Info("annotations", annotations)

	res := buildUserResult(user, len(joinedDatasets), len(createdDatasets), len(annotations))

	return res, nil
}

func (u *User) ChangeInfo(info dto.ChangeUserInfo, userId uint) (user *User, err error) {
	user, err = dao.First[User](userId)
	if err != nil {
		return nil, err
	}

	if info.Name != "" {
		user.Name = info.Name
	}
	if info.Avatar != "" {
		user.Avatar = info.Avatar
	}
	if info.Email != "" {
		user.Email = info.Email
	}
	if info.Description != "" {
		user.Description = info.Description
	}

	err = dao.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
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

// ListUsersByIds 根据 ID 列出用户
func (u *User) ListUsersByIds(ids []uint) ([]User, error) {
	var err error
	var users []User
	users, err = dao.FindAll[User]("id in (?)", ids)
	if err != nil {
		return nil, err
	}
	return users, nil
}
