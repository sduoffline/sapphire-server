package dto

// Register 注册请求参数
type Register struct {
	Name   string `json:"name" binding:"required"`
	Passwd string `json:"passwd" binding:"required"`
}
