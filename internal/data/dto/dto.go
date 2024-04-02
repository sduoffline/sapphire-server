package dto

// Request
//
// 定义一个interface，用于对请求参数进行规范
// 要求实现该接口的struct必须实现Validate方法
type Request interface {
	// Validate 验证请求参数
	Validate() error
	// BindJSON 绑定请求参数
	BindJSON() error
}
