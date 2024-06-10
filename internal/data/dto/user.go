package dto

// Register 注册请求参数
type Register struct {
	Name   string `json:"name" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}

// Login 登录请求参数
type Login struct {
	Name   string `json:"name" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
}

// ChangePasswd 修改密码请求参数
type ChangePasswd struct {
	Old string `json:"old" binding:"required"`
	New string `json:"new" binding:"required"`
}
