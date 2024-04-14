package dto

// Register 注册请求参数
type Register struct {
	Name   string `json:"name" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Avatar string `json:"avatar" binding:"required"`
}

type Login struct {
	Name   string `json:"name" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
}
