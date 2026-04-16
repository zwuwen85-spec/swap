package handler

import (
	"campus-swap-shop/internal/logic"
	"campus-swap-shop/internal/model"
	"campus-swap-shop/internal/svc"
	"campus-swap-shop/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExchangeHandler 交换处理器
type ExchangeHandler struct {
	exchangeLogic *logic.ExchangeLogic
}

// NewExchangeHandler 创建交换处理器
func NewExchangeHandler(serviceCtx *svc.ServiceContext) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeLogic: logic.NewExchangeLogic(serviceCtx.DB),
	}
}

// Create 发起交换请求
func (h *ExchangeHandler) Create(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.ExchangeCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 创建交换请求
	exchange, err := h.exchangeLogic.Create(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		if err.Error() == "商品不存在或已下架" {
			response.Error(c, response.CodeGoodsNotFound, err.Error())
		} else if err.Error() == "不能给自己的商品发起交换" {
			response.Error(c, response.CodeDuplicate, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	response.SuccessWithMessage(c, "交换请求已发送", gin.H{
		"exchange_id": exchange.ID,
	})
}

// GetList 获取交换列表
func (h *ExchangeHandler) GetList(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.ExchangeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 获取列表
	list, total, err := h.exchangeLogic.GetList(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	// 分页参数
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	response.PageSuccess(c, list, total, page, pageSize)
}

// GetDetail 获取交换详情
func (h *ExchangeHandler) GetDetail(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	exchangeIDStr := c.Query("id")
	if exchangeIDStr == "" {
		response.Error(c, response.CodeParamError, "缺少交换ID")
		return
	}

	exchangeID, err := strconv.ParseInt(exchangeIDStr, 10, 64)
	if err != nil {
		response.Error(c, response.CodeParamError, "交换ID格式错误")
		return
	}

	// 获取详情
	exchange, err := h.exchangeLogic.GetDetail(c.Request.Context(), exchangeID)
	if err != nil {
		if err.Error() == "交换请求不存在" {
			response.Error(c, response.CodeExchangeNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	// 检查权限
	if exchange.InitiatorID != userID.(int64) && exchange.TargetID != userID.(int64) {
		response.Error(c, response.CodeNoPermission, "无权查看此交换请求")
		return
	}

	response.Success(c, exchange)
}

// Handle 处理交换请求
func (h *ExchangeHandler) Handle(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	var req model.HandleExchangeDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeParamError, "参数错误："+err.Error())
		return
	}

	// 处理交换请求
	if err := h.exchangeLogic.Handle(c.Request.Context(), userID.(int64), &req); err != nil {
		if err.Error() == "无权操作此交换请求" {
			response.Error(c, response.CodeNoPermission, err.Error())
		} else if err.Error() == "交换请求不存在" {
			response.Error(c, response.CodeExchangeNotFound, err.Error())
		} else {
			response.Error(c, response.CodeServerError, err.Error())
		}
		return
	}

	// 根据操作类型返回不同消息
	var message string
	switch req.Action {
	case "accept":
		message = "已接受交换请求"
	case "reject":
		message = "已拒绝交换请求"
	case "cancel":
		message = "已取消交换请求"
	case "complete":
		message = "交换已完成"
	}

	response.SuccessWithMessage(c, message, nil)
}

// GetPendingCount 获取待处理交换数量
func (h *ExchangeHandler) GetPendingCount(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeNotLogin, "未登录")
		return
	}

	count, err := h.exchangeLogic.GetPendingCount(c.Request.Context(), userID.(int64))
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"count": count,
	})
}
