package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// PageResponse 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// 状态码定义
const (
	CodeSuccess      = 0     // 成功
	CodeParamError   = 10001 // 参数错误
	CodeNotLogin     = 10002 // 未登录
	CodeNoPermission = 10003 // 无权限
	CodeNotFound     = 10004 // 资源不存在
	CodeDuplicate    = 10005 // 重复操作

	CodeUserNotFound     = 20001 // 用户不存在
	CodePasswordError    = 20002 // 密码错误
	CodeUserAlreadyExist = 20003 // 用户已存在

	CodeGoodsNotFound  = 30001 // 商品不存在
	CodeGoodsOffShelf  = 30002 // 商品已下架
	CodeGoodsSoldOut   = 30003 // 商品已售出

	CodeExchangeNotFound = 40001 // 交换请求不存在
	CodeCannotHandleSelf = 40002 // 不能处理自己的请求

	CodeServerError = 50001 // 服务器错误
)

// 常用错误别名（向后兼容）
const (
	ErrInvalidParam = CodeParamError
	ErrServerError  = CodeServerError
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithData 错误响应（带数据）
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageResponse{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}
