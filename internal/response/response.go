package response

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 错误码常量
const (
	CodeSuccess      = 0
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeConflict     = 409
	CodeInternal     = 500
)

// HandleError 统一错误处理
func HandleError(c *gin.Context, err error) {
	errMsg := err.Error()

	switch {
	case errMsg == "tenant not found" || errMsg == "tenant code already exists":
		if errMsg == "tenant code already exists" {
			Error(c, CodeConflict, "Tenant code already exists")
		} else {
			Error(c, CodeNotFound, "Tenant not found")
		}
	case len(errMsg) > 0 && errMsg[0] == '[': // validator errors
		Error(c, CodeBadRequest, errMsg)
	default:
		Error(c, CodeInternal, "Internal server error")
	}
}
