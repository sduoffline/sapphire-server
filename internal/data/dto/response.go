package dto

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func ok() {
	NewResponse(200, "success", nil)
}

func NewSuccessResponse(data interface{}) *Response {
	return NewResponse(200, "success", data)
}

func NewFailResponse(message string) *Response {
	return NewResponse(500, message, nil)
}
